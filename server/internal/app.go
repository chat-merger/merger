package internal

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/pelletier/go-toml"
	"github.com/urfave/cli"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/chat-merger/merger/server/internal/callback"
	"github.com/chat-merger/merger/server/internal/handlers"
)

type App struct {
	ConfigPath string
	cfg        *Config
	cbClient   callback.Client
	db         *gorm.DB
}

func (a *App) CBClient() callback.Client { return a.cbClient }
func (a *App) DB() *gorm.DB              { return a.db }

type Config struct {
	Port            int    `toml:"port"`
	DbConnection    string `toml:"db"`
	RedisConnection string `toml:"redis"`
}

func (a *App) Start(c *cli.Context) (err error) {
	if a.cfg, err = parseConfig(a.ConfigPath); err != nil {
		return err
	}

	a.cbClient = configureCallbackApi()

	var closeDB func() error
	if a.db, closeDB, err = connectToDB(a.cfg.DbConnection); err != nil {
		return err
	}
	defer closeDB()

	return startAndListenServer(a)
}

func startAndListenServer(app *App) error {
	router := http.NewServeMux()
	handlers.Setup(app, router)

	serv := http.Server{
		Addr:    fmt.Sprintf(":%d", app.cfg.Port),
		Handler: handlePanic(router),
	}
	defer serv.Close()

	return serv.ListenAndServe()
}

func handlePanic(handler http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			rErr := recover()
			if rErr == nil {
				return
			}

			if w != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Header().Set("Content-Type", "application/json; charset=UTF-8")
				if _, err := w.Write([]byte(`{"code":500,"message":"panic service"}`)); err != nil {
					slog.Error(err.Error())
				}
			}
		}()

		handler.ServeHTTP(w, r)
	}
}

func connectToDB(dsn string) (*gorm.DB, func() error, error) {
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, nil, err
	}
	db.Logger.LogMode(logger.Info)

	sqlDB, err := db.DB()
	if err != nil {
		return nil, nil, err
	}

	return db, sqlDB.Close, nil
}

func configureCallbackApi() callback.Client {
	return callback.NewClient()
}

func parseConfig(path string) (*Config, error) {
	var (
		err  error
		file []byte
		cfg  = &struct {
			Merger *Config `toml:"merger"`
		}{}
	)
	if file, err = os.ReadFile(path); err != nil {
		return nil, err
	}

	err = toml.Unmarshal(file, cfg)
	if err != nil {
		return nil, err
	}

	return cfg.Merger, nil
}

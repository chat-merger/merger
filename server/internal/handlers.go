package internal

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

func setup(router *http.ServeMux, app *App) {
	router.HandleFunc("GET /", wrap(app, handlerOk))
}

func wrap(app *App, fn func(app *App, r *http.Request) (any, error)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if response, err := fn(app, r); err != nil {
			w.WriteHeader(http.StatusBadRequest)
		} else {
			w.WriteHeader(http.StatusOK)
			var b []byte
			if b, err = json.Marshal(response); err != nil {
				slog.Error(err.Error())
			} else {
				if _, err = w.Write(b); err != nil {
					slog.Error(err.Error())
				}
			}
		}

	}
}

func handlerOk(app *App, r *http.Request) (any, error) {
	return struct{ R string }{R: "OK"}, nil
}

package internal

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
)

type Event interface{}
type File interface{}
type Context struct {
	app  *App
	repo Repo
}

type Repo interface {
	EventFrom(id int) []Event
	EventLast() Event
	EventSend(Event)

	File(id int) File
	FileUpload(File)

	Assoc(clientID string) int
	AssocWrite(clientID string, mergerID int)
}

type emptyRepo struct{}

func (_ emptyRepo) EventFrom(id int) []Event                 { return nil }
func (_ emptyRepo) EventLast() Event                         { return nil }
func (_ emptyRepo) EventSend(Event)                          {}
func (_ emptyRepo) File(id int) File                         { return nil }
func (_ emptyRepo) FileUpload(File)                          {}
func (_ emptyRepo) Assoc(clientID string) int                { return 1 }
func (_ emptyRepo) AssocWrite(clientID string, mergerID int) {}

func setup(router *http.ServeMux, app *App) {
	ctx := Context{app: app, repo: emptyRepo{}}
	router.HandleFunc("GET /", wrap(ctx, handlerOk))
	router.HandleFunc("GET /event", wrap(ctx, handlerEvent))
	router.HandleFunc("POST /event", wrap(ctx, handlerEventSend))
	router.HandleFunc("GET /file", wrap(ctx, handlerFile))
	router.HandleFunc("POST /file", wrap(ctx, handlerFileUpload))
	router.HandleFunc("GET /assoc", wrap(ctx, handlerAssoc))
	router.HandleFunc("POST /assoc", wrap(ctx, handlerAssocWrite))
}

func wrap(c Context, fn func(c Context, r *http.Request) (any, error)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if response, err := fn(c, r); err != nil {
			w.WriteHeader(http.StatusBadRequest)
		} else {
			w.WriteHeader(http.StatusOK)
			var b []byte
			if b, err = json.Marshal(response); err != nil {
				b = []byte(fmt.Sprintf("%v", response))
			}
			if _, err = w.Write(b); err != nil {
				slog.Error(err.Error())
			}

		}
	}
}

func text(str string) (any, error) {
	return struct {
		Text string `json:"text"`
	}{Text: str}, nil
}

func handlerOk(_ Context, _ *http.Request) (any, error) {
	return text("OK")
}

func handlerEvent(c Context, r *http.Request) (_ any, err error) {
	var last, from int
	if last, err = intFromQury(r, "last"); err != nil {
		return nil, err
	}
	if last != 0 {
		return c.repo.EventLast(), nil
	}

	if from, err = intFromQury(r, "from"); err != nil {
		return nil, err
	}
	if from != 0 {
		return c.repo.EventLast(), nil
	}

	return text("не переданы параметры")
}
func handlerEventSend(c Context, r *http.Request) (any, error)  { return nil, nil }
func handlerFile(c Context, r *http.Request) (any, error)       { return nil, nil }
func handlerFileUpload(c Context, r *http.Request) (any, error) { return nil, nil }
func handlerAssoc(c Context, r *http.Request) (any, error)      { return nil, nil }
func handlerAssocWrite(c Context, r *http.Request) (any, error) { return nil, nil }

func intFromQury(r *http.Request, name string) (int, error) {
	queryVal := r.URL.Query().Get(name)
	if queryVal == "" {
		return 0, nil
	}
	if val, err := strconv.Atoi(queryVal); err != nil {
		return 0, err
	} else {
		return val, nil
	}
}

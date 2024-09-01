package handlers

import (
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"strconv"

	"gorm.io/gorm"

	"github.com/chat-merger/merger/server/internal/callback"
	"github.com/chat-merger/merger/server/internal/event"
	"github.com/chat-merger/merger/server/internal/operation"
)

type Context interface {
	CallbackApi() callback.API
	DB() *gorm.DB
}

func Setup(c Context, router *http.ServeMux) {
	router.HandleFunc("GET /", wrap(c, handlerOk))
	router.HandleFunc("POST /events/newMessage", wrap(c, handlerEventMessageNew))
	router.HandleFunc("POST /files", wrap(c, handlerFileUpload))
}

func handlerOk(_ Context, _ *http.Request) (any, error) {
	return text("OK")
}

func handlerEventMessageNew(c Context, r *http.Request) (_ any, err error) {
	var appID int
	if appID, err = headerAppID(r); err != nil {
		return nil, err
	}

	var newMessage event.MessageNew
	if err = json.NewDecoder(r.Body).Decode(&newMessage); err != nil {
		return nil, err
	}
	newMessage.AppID = appID

	if err = operation.MessageNew(c, newMessage); err != nil {
		return text("operation.MessageNew: " + err.Error())
	}

	return nil, nil
}

func handlerFileUpload(c Context, r *http.Request) (_ any, err error) {
	var appID int
	if appID, err = headerAppID(r); err != nil {
		return nil, err
	}
	if err = r.ParseMultipartForm(21 << 20); err != nil {
		return nil, err
	}

	file := event.FileUpload{AppID: appID}

	// MsgLID
	if len(r.MultipartForm.Value["id"]) != 1 {
		return text("id.len must have equals 1")
	}
	file.LocalID = r.MultipartForm.Value["id"][0]

	// Type
	if len(r.MultipartForm.Value["type"]) != 1 {
		return text("type.len must have equals 1")
	}
	file.Type, err = strconv.Atoi(r.MultipartForm.Value["type"][0])

	// Bytes
	fileHeader := r.MultipartForm.File["file"]
	if len(fileHeader) != 1 {
		return text("file.len must have equals 1")
	}
	var mpFile multipart.File
	if mpFile, err = fileHeader[0].Open(); err != nil {
		return text("Unable to open file")
	}
	defer mpFile.Close()
	if file.Bytes, err = io.ReadAll(mpFile); err != nil {
		return text("Unable to read file")
	}

	if err = operation.FileUpload(c, file); err != nil {
		return text("operation.FileUpload: " + err.Error())
	}

	return text("ok")
}

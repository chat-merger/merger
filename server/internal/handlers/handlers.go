package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"mime/multipart"
	"net/http"
	"strconv"

	"gorm.io/gorm"

	"github.com/chat-merger/merger/server/internal/callback"
	"github.com/chat-merger/merger/server/internal/event"
	"github.com/chat-merger/merger/server/internal/event/file/upload"
	eventMessageNew "github.com/chat-merger/merger/server/internal/event/message/new"
)

type Context interface {
	CBClient() callback.Client
	DB() *gorm.DB
}

func Setup(c Context, router *http.ServeMux) {
	router.HandleFunc("POST /events/newMessage", wrap(c, handlerEventMessageNew))
	router.HandleFunc("POST /files", wrap(c, handlerFileUpload))
	router.HandleFunc("GET /", wrap(c, handlerOk))
	router.HandleFunc("POST /echo", wrap(c, handlerEcho))
	router.HandleFunc("POST /test/app", wrap(c, handlerApp))
}

func handlerApp(_ Context, r *http.Request) (_ any, err error) {
	var requestBody callback.Body
	if err = json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		return nil, err
	}
	switch {
	case requestBody.MessageNew != nil:
		return callback.MessageNewResponse{
			LocalID: strconv.Itoa(rand.Int()),
		}, nil
	default:
		return text("unexpected callback")
	}
}

func handlerEcho(_ Context, r *http.Request) (any, error) {
	var b, _ = io.ReadAll(r.Body)
	fmt.Println(string(b))
	return text("OK")
}

func handlerOk(_ Context, _ *http.Request) (any, error) {
	return text("OK")
}

func handlerEventMessageNew(c Context, r *http.Request) (_ any, err error) {
	var appID int
	if appID, err = headerAppID(r); err != nil {
		return nil, err
	}

	var newMessage eventMessageNew.Message
	if err = json.NewDecoder(r.Body).Decode(&newMessage); err != nil {
		return nil, err
	}
	newMessage.AppID = appID

	if err = eventMessageNew.Exec(c, newMessage); err != nil {
		return text("eventMessageNew.Exec: " + err.Error())
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

	if err = upload.FileUpload(c, file); err != nil {
		return text("event.FileUpload: " + err.Error())
	}

	return text("ok")
}

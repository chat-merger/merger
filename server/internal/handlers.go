package internal

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"mime/multipart"
	"net/http"
	"strconv"
	"strings"
)

type File struct {
	Bytes   []byte
	Type    int
	InAppID string
}
type Context struct {
	app *App
}

func setup(router *http.ServeMux, app *App) {
	ctx := Context{app: app}
	router.HandleFunc("GET /", wrap(ctx, handlerOk))
	router.HandleFunc("POST /events/newMessage", wrap(ctx, handlerNewMessage))
	router.HandleFunc("POST /files", wrap(ctx, handlerFilesUpload))
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

type EventNewMsg struct {
	InAppID            string
	IsSilent           bool
	ReplyInAppID       string
	Username           string
	Text               string
	AttachmentInAppIDs []string
	Forwards           []EventNewMsgForward
	Attachments        EventNewMsgAttachment
}

type EventNewMsgForward struct {
	InAppID            int
	Username           string
	Text               string
	CreateDate         string
	AttachmentInAppIDs []string
}

type EventNewMsgAttachment struct {
	InAppID    string
	HasSpoiler bool
	Type       int
	// Общедоступная ссылка для загрузки файла.
	// Если ссылка не передана и по такому FileID не найдено файлов, то клиент должен будет загрузить файлы на специальный эндпоинт
	Url string
}

type EventNewMsgResponse struct {
	ID                               int
	RequireUploadAttachmentsInAppIDs []string
}

func handlerNewMessage(c Context, r *http.Request) (_ any, err error) {
	var b []byte
	if _, err = r.Body.Read(b); err != nil {
		return text("read body err: " + err.Error())
	}

	var newMsg EventNewMsg
	if err = json.Unmarshal(b, &newMsg); err != nil {
		return nil, err
	}

	EventNewMessage(c.app, newMsg)

	return nil, nil
}

func handlerFilesUpload(c Context, r *http.Request) (_ any, err error) {
	if err = r.ParseMultipartForm(21 << 20); err != nil {
		return nil, err
	}
	var file File

	// InAppID
	if len(r.MultipartForm.Value["id"]) != 1 {
		return text("id.len must have equals 1")
	}
	file.InAppID = r.MultipartForm.Value["id"][0]

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

	FileUpload(c.app, file)

	return text("ok")
}

func intFromQuery(r *http.Request, name string) (int, error) {
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

func intsFromQuery(r *http.Request, name string) ([]int, error) {
	queryVal := r.URL.Query().Get(name)
	if queryVal == "" {
		return nil, nil
	}
	strValues := strings.Split(queryVal, ",")
	result := make([]int, len(strValues))
	for i, v := range strValues {
		if val, err := strconv.Atoi(v); err != nil {
			return nil, err
		} else {
			result[i] = val
		}
	}

	return result, nil
}

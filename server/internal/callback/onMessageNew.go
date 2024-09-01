package callback

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/chat-merger/merger/server/internal/model"
)

type MessageNewRequest struct {
	MessageNew
}

type MessageNewResponse struct {
	LocalID string `json:"local_id,omitempty"`
	MsgID   int    `json:"msg_id,omitempty"`
	AppID   int    `json:"app_id,omitempty"`
}

type MessageNew struct {
	ID          int             `json:"id,omitempty"`
	IsSilent    bool            `json:"is_silent,omitempty"`
	Reply       int             `json:"reply,omitempty"`
	ReplyLocal  string          `json:"reply_local,omitempty"`
	Username    string          `json:"username,omitempty"`
	Text        string          `json:"text,omitempty"`
	Attachments []AttachmentNew `json:"attachments,omitempty"`
	Forwards    []ForwardNew    `json:"forwards,omitempty"`

	App model.Application `json:"-"`
}

type ForwardNew struct {
	ID          int             `json:"id,omitempty"`
	LocalID     string          `json:"local_id,omitempty"`
	Username    string          `json:"username,omitempty"`
	Text        string          `json:"text,omitempty"`
	CreateDate  string          `json:"create_date,omitempty"`
	Attachments []AttachmentNew `json:"attachments,omitempty"`
}

type AttachmentNew struct {
	HasSpoiler bool   `json:"has_spoiler,omitempty"`
	Type       int    `json:"type,omitempty"`
	Url        string `json:"url,omitempty"`
	WaitUpload bool   `json:"wait_upload,omitempty"` // false - ждать ивента о том что файл стал доступен по ссылке
}

func (a *api) OnMessageNew(newMessages []MessageNew) ([]MessageNewResponse, error) {
	result := make([]MessageNewResponse, len(newMessages))
	for i, newMsg := range newMessages {
		b, err := json.Marshal(Body{MessageNew: &newMsg})
		if err != nil {
			return nil, fmt.Errorf("json.Marshal: %w ", err)
		}
		r, err := a.cl.Post(newMsg.App.Host, "application/json", bytes.NewBuffer(b))
		if err != nil {
			return nil, fmt.Errorf("a.cl.Post: %w", err)
		}
		var resp MessageNewResponse
		if err = json.NewDecoder(r.Body).Decode(&resp); err != nil {
			return nil, fmt.Errorf("json.NewDecoder.Decode: %w", err)
		}
		resp.MsgID = newMsg.ID
		resp.AppID = newMsg.App.ID

		result[i] = resp
	}

	return result, nil
}

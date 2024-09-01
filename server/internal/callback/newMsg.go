package callback

import "github.com/chat-merger/merger/server/internal/model"

type NewMsgRequest struct {
	NewMessage
}

type NewMsgResponse struct {
	LocalID string `json:"localId,omitempty"`
	MsgID   int    `json:"msgId,omitempty"`
	AppID   int    `json:"appId,omitempty"`
}

type NewMessage struct {
	ID          int             `json:"id,omitempty"`
	IsSilent    bool            `json:"isSilent,omitempty"`
	Reply       int             `json:"reply,omitempty"`
	ReplyLocal  string          `json:"replyLocal,omitempty"`
	Username    string          `json:"username,omitempty"`
	Text        string          `json:"text,omitempty"`
	Attachments []NewAttachment `json:"attachments,omitempty"`
	Forwards    []NewForward    `json:"forwards,omitempty"`

	App model.Application `json:"-"`
}

type NewForward struct {
	ID          int             `json:"id,omitempty"`
	LocalID     string          `json:"localId,omitempty"`
	Username    string          `json:"username,omitempty"`
	Text        string          `json:"text,omitempty"`
	CreateDate  string          `json:"createDate,omitempty"`
	Attachments []NewAttachment `json:"attachments,omitempty"`
}

type NewAttachment struct {
	HasSpoiler bool   `json:"hasSpoiler,omitempty"`
	Type       int    `json:"type,omitempty"`
	Url        string `json:"url,omitempty"`
	WaitUpload bool   `json:"waitUpload,omitempty"` // false - ждать ивента о том что файл стал доступен по ссылке
}

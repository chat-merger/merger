package callback

type NewMsgRequest struct {
	NewMessage
}

type NewMsgResponse struct {
	LocalID string
}

type NewMessage struct {
	ID          int             `json:"ID,omitempty"`
	IsSilent    bool            `json:"isSilent,omitempty"`
	Reply       int             `json:"reply,omitempty"`
	ReplyLocal  string          `json:"replyLocal,omitempty"`
	Username    string          `json:"username,omitempty"`
	Text        string          `json:"text,omitempty"`
	Attachments []NewAttachment `json:"attachments,omitempty"`
	Forwards    []NewForward    `json:"forwards,omitempty"`
}

type NewForward struct {
	ID          int
	LocalID     string
	Username    string
	Text        string
	CreateDate  string
	Attachments []NewAttachment
}

type NewAttachment struct {
	HasSpoiler bool
	Type       int
	Url        string
	WaitUpload bool // false - ждать ивента о том что файл стал доступен по ссылке
}

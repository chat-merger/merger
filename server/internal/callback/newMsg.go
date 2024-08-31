package callback

type NewMsgRequest struct {
	NewMessage
}

type NewMsgResponse struct {
	LocalID string
}

type NewMessage struct {
	ID          int
	IsSilent    bool
	Reply       int
	ReplyLocal  string
	Username    string
	Text        string
	Attachments []NewAttachment
	Forwards    []NewForward
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

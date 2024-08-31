package internal

type CallbackNewMsgForward struct {
	InAppID       int
	Username      string
	Text          string
	CreateDate    string
	AttachmentIDs []string
}

type CallbackNewMsg struct {
	ID            int
	IsSilent      bool
	Reply         int
	ReplyInAppID  int
	Username      string
	Text          string
	AttachmentIDs []string
	Forwards      []CallbackNewMsgForward
	Attachments   map[string]CallbackNewMsgAttachment
}

type CallbackNewMsgResponse struct {
	InAppID string
}

type CallbackNewMsgAttachment struct {
	HasSpoiler bool
	Type       AttachmentType
	Url        string
	WaitUpload bool // false - ждать ивента о том что файл стал доступен по ссылке
}

type CallbackAPI interface {
	OnNewMsg(CallbackNewMsg) CallbackNewMsgResponse
}

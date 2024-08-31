package internal

type callbackNewMsgForward struct {
	InAppID     int
	Username    string
	Text        string
	CreateDate  string
	Attachments []callbackNewMsgAttachment
}

type callbackNewMsg struct {
	ID           int
	IsSilent     bool
	Reply        int
	ReplyInAppID string
	Username     string
	Text         string
	Attachments  []callbackNewMsgAttachment
	Forwards     []callbackNewMsgForward
}

type callbackNewMsgResponse struct {
	InAppID string
}

type callbackNewMsgAttachment struct {
	HasSpoiler bool
	Type       int
	Url        string
	WaitUpload bool // false - ждать ивента о том что файл стал доступен по ссылке
}

type callbackAPI interface {
	OnNewMsg(c map[int]callbackNewMsg) callbackNewMsgResponse
}

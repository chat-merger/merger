package internal

type Application struct {
	ID   int    `gorm:"column:id;primary_key"`
	Name string `gorm:"column:name"`
	XKey string `gorm:"column:xkey"`
	Host string `gorm:"column:host"`
}

const TableApplications = "Applications"

func (*Application) TableName() string { return TableApplications }

func collectHosts(apps []Application) []string {
	hosts := make([]string, len(apps))
	for i := range apps {
		hosts[i] = apps[i].Host
	}

	return hosts
}

type Message struct {
	ID    int `gorm:"column:id;primary_key"`
	AppID int `gorm:"column:appId"`
	//IsSilent   bool   `gorm:"column:isSilent"`
	//IsForward  bool   `gorm:"column:isForward"`
	Reply int `gorm:"column:reply"`
	//Username   string `gorm:"column:username"`
	//Text       string `gorm:"column:text"`
	//CreateDate int `gorm:"column:createDate"`
}

type MessageExt struct {
	Message
	Attachments []*Attachment
}

const TableMessages = "Messages"

func (*Message) TableName() string { return TableMessages }

type MessageMap struct {
	AppID      int    `gorm:"column:appId"`
	MsgID      int    `gorm:"column:msgID"`
	MsgLocalID string `gorm:"column:msgLocalID"`
}

const TableMessagesMap = "MessagesMap"

func (*MessageMap) TableName() string { return TableMessagesMap }

type Attachment struct {
	ID         int    `gorm:"column:id"`
	LocalID    string `gorm:"column:localId"`
	AppID      int    `gorm:"column:appId"`
	Url        string `gorm:"column:url"`
	HasSpoiler bool   `gorm:"column:hasSpoiler"`
	Type       int    `gorm:"column:type"`
	//FileName   string `gorm:"column:fileName"`
	//Confirmed  bool   `gorm:"column:confirmed"`
}

func CollectAttachmentID(attachments []*Attachment) []int {
	result := make([]int, len(attachments))
	for i, fwd := range attachments {
		result[i] = fwd.ID
	}

	return result
}

const TableAttachments = "Attachments"

func (*Attachment) TableName() string { return TableAttachments }

type AttachmentType int

const (
	AtAudio   = 1
	AtVideo   = 2
	AtFile    = 3
	AtPhoto   = 4
	AtSticker = 5
)

type File struct {
	ID           int    `gorm:"column:id"`
	AttachmentID int    `gorm:"column:attachmentId"`
	FileName     string `gorm:"column:fileName"`
}

const TableFiles = "Files"

func (*File) TableName() string { return TableFiles }

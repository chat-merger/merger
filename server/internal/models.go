package internal

type Application struct {
	ID       int    `gorm:"column:id;primary_key"`
	Name     string `gorm:"column:name"`
	XKey     string `gorm:"column:xkey"`
	Callback string `gorm:"column:callback"`
}

func (*Application) TableName() string { return "Applications" }

type Message struct {
	ID         int    `gorm:"column:id;primary_key"`
	IsSilent   bool   `gorm:"column:isSilent"`
	IsForward  bool   `gorm:"column:isForward"`
	Reply      int    `gorm:"column:reply"`
	Username   string `gorm:"column:username"`
	Text       string `gorm:"column:text"`
	CreateDate int    `gorm:"column:createDate"`
}

func (*Message) TableName() string { return "Messages" }

type MessageMap struct {
	MsgID   int    `gorm:"column:msgId"`
	InAppID string `gorm:"column:inAppId"`
}

func (*MessageMap) TableName() string { return "MessagesMap" }

type Attachments struct {
	MsgID      int            `gorm:"column:msgId"`
	FileID     int            `gorm:"column:fileId"`
	HasSpoiler bool           `gorm:"column:hasSpoiler"`
	Type       AttachmentType `gorm:"column:type"`
}

func (*Attachments) TableName() string { return "Attachments" }

type AttachmentType int

const (
	AtAudio   = 1
	AtVideo   = 2
	AtFile    = 3
	AtPhoto   = 4
	AtSticker = 5
)

type Files struct {
	ID       int    `gorm:"column:id"`
	FileName string `gorm:"column:fileName"`
}

func (*Files) TableName() string { return "Files" }

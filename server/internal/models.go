package internal

type Application struct {
	ID   int    `gorm:"column:id;primary_key"`
	Name string `gorm:"column:name"`
	XKey string `gorm:"column:xkey"`
	Host string `gorm:"column:host"`
}

func (*Application) TableName() string { return "Applications" }

func collectHosts(apps []Application) []string {
	hosts := make([]string, len(apps))
	for i := range apps {
		hosts[i] = apps[i].Host
	}

	return hosts
}

type Message struct {
	ID int `gorm:"column:id;primary_key"`
	//IsSilent   bool   `gorm:"column:isSilent"`
	//IsForward  bool   `gorm:"column:isForward"`
	Reply int `gorm:"column:reply"`
	//Username   string `gorm:"column:username"`
	//Text       string `gorm:"column:text"`
	CreateDate int `gorm:"column:createDate"`
}

func (*Message) TableName() string { return "Messages" }

type MessageMap struct {
	AppID   int    `gorm:"column:appId"`
	MsgID   int    `gorm:"column:msgId"`
	InAppID string `gorm:"column:inAppId"`
}

func msgMapByAppID(msgID int) map[int]MessageMap {

}

func (*MessageMap) TableName() string { return "MessagesMap" }

type Attachment struct {
	ID         int    `gorm:"column:id"`
	AppID      int    `gorm:"column:fileId"`
	InAppID    string `gorm:"column:inAppId"`
	Url        string `gorm:"column:url"`
	HasSpoiler bool   `gorm:"column:hasSpoiler"`
	Type       int    `gorm:"column:type"`
	WaitUpload bool   `gorm:"column:waitUpload"`
}

func (*Attachment) TableName() string { return "Attachments" }

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

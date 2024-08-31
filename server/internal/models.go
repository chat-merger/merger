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
	Attachments []Attachment
}

func (*Message) TableName() string { return "Messages" }

type MessageMap struct {
	AppID      int    `gorm:"column:appId"`
	MsgID      int    `gorm:"column:msgID"`
	MsgLocalID string `gorm:"column:msgLocalID"`
}

func (*MessageMap) TableName() string { return "MessagesMap" }

type Attachment struct {
	ID         int    `gorm:"column:id"`
	LocalID    string `gorm:"column:localId"`
	AppID      int    `gorm:"column:appId"`
	Url        string `gorm:"column:url"`
	HasSpoiler bool   `gorm:"column:hasSpoiler"`
	Type       int    `gorm:"column:type"`
	FileName   string `gorm:"column:fileName"`
	Confirmed  bool   `gorm:"column:confirmed"`
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

//type Files struct {
//	ID       int    `gorm:"column:id"`
//	AppID    int    `gorm:"column:appId"`
//	FileName string `gorm:"column:fileName"`
//}
//
//func (*Files) TableName() string { return "Files" }

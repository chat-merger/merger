package model

type Message struct {
	ID    int `gorm:"column:id;primary_key"`
	AppID int `gorm:"column:appId"`
	Reply int `gorm:"column:reply"`
}

type MessageExt struct {
	Message
	Attachments []*Attachment
}

const TableMessages = "Messages"

var InstMessage = new(Message)

func (*Message) TableName() string { return TableMessages }

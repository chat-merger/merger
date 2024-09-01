package model

type MessageMap struct {
	AppID      int    `gorm:"column:appId"`
	MsgID      int    `gorm:"column:msgID"`
	MsgLocalID string `gorm:"column:msgLocalID"`
}

const TableMessagesMap = "MessagesMap"

var InstMessageMap = new(MessageMap)

func (*MessageMap) TableName() string { return TableMessagesMap }

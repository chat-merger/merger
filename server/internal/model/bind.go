package model

type Bind struct {
	AppID      int    `gorm:"column:appId"`
	MsgID      int    `gorm:"column:msgID"`
	MsgLocalID string `gorm:"column:msgLocalID"`
}

const TableBinds = "Binds"

var InstBind = new(Bind)

func (*Bind) TableName() string { return TableBinds }

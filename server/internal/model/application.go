package model

type Application struct {
	ID   int    `gorm:"column:id;primary_key"`
	Name string `gorm:"column:name"`
	XKey string `gorm:"column:xkey"`
	Host string `gorm:"column:host"`
}

const TableApplications = "Applications"

func (*Application) TableName() string { return TableApplications }

package model

type File struct {
	ID           int    `gorm:"column:id"`
	AttachmentID int    `gorm:"column:attachmentId"`
	FileName     string `gorm:"column:fileName"`
}

const TableFiles = "Files"

func (*File) TableName() string { return TableFiles }

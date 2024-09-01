package model

type Attachment struct {
	ID         int    `gorm:"column:id"`
	LocalID    string `gorm:"column:localId"`
	AppID      int    `gorm:"column:appId"`
	MsgID      int    `gorm:"column:msgId"`
	Url        string `gorm:"column:url"`
	HasSpoiler bool   `gorm:"column:hasSpoiler"`
	Type       int    `gorm:"column:type"`
}

const TableAttachments = "Attachments"

var InstAttachment = new(Attachment)

func (*Attachment) TableName() string { return TableAttachments }

type AttachmentType int

const (
	AtAudio   = 1
	AtVideo   = 2
	AtFile    = 3
	AtPhoto   = 4
	AtSticker = 5
)

func CollectAttachmentID(attachments []*Attachment) []int {
	result := make([]int, len(attachments))
	for i, fwd := range attachments {
		result[i] = fwd.ID
	}

	return result
}

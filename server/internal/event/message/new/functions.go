package new

import (
	"slices"

	"github.com/chat-merger/merger/server/internal/callback"
	"github.com/chat-merger/merger/server/internal/event"
	"github.com/chat-merger/merger/server/internal/model"
)

type ForwardExt struct {
	ID          int
	LocalID     string
	Username    string
	Text        string
	CreateDate  string
	Attachments []*model.Attachment
}

func NewMsgResponseToMsgMap(responses []callback.MessageNewResponse) []*model.MessageMap {
	mm := make([]*model.MessageMap, len(responses))
	for i, response := range responses {
		mm[i] = &model.MessageMap{
			AppID:      response.AppID,
			MsgID:      response.MsgID,
			MsgLocalID: response.LocalID,
		}
	}

	return mm
}

func ForwardExtToCbkForwards(exts []ForwardExt, attachWaitingIDs []int) []callback.ForwardNew {
	newForwards := make([]callback.ForwardNew, 0, len(exts))
	for _, ext := range exts {
		newForwards = append(newForwards, callback.ForwardNew{
			ID:          ext.ID,
			LocalID:     ext.LocalID,
			Username:    ext.Username,
			Text:        ext.Text,
			CreateDate:  ext.CreateDate,
			Attachments: AttachmentToCbkAttachs(ext.Attachments, attachWaitingIDs),
		})
	}

	return newForwards
}

func SaveMessage(c event.Context, e Message) (model.MessageExt, []ForwardExt, error) {
	// Save message
	msg := model.MessageExt{Message: model.Message{AppID: e.AppID}}
	var reply struct {
		ID int `gorm:"column:msgId"`
	}
	if err := c.DB().
		Model(model.InstMessageMap).
		Select("msgId").
		Where("appId = ?", e.AppID).
		Where("msgLocalId = ?", e.ReplyLocalID).
		Find(&reply).Error; err != nil {
		return model.MessageExt{}, nil, err
	}
	msg.Reply = reply.ID

	if err := c.DB().Create(&msg.Message).Error; err != nil {
		return model.MessageExt{}, nil, err
	}

	// Fill message attachments
	msg.Attachments = make([]*model.Attachment, len(e.Forwards))
	for i, attach := range e.Attachments {
		msg.Attachments[i] = &model.Attachment{
			LocalID:    attach.LocalID,
			AppID:      e.AppID,
			Url:        attach.Url,
			HasSpoiler: attach.HasSpoiler,
			Type:       attach.Type,
		}
	}

	var forwardsAttachments []*model.Attachment

	fwdLocalIDs := CollectMessageNewForwardsLocalIDs(e.Forwards)
	var messagesMap []model.MessageMap
	if err := c.DB().
		Table(model.TableMessagesMap).
		Where("appId = ?", e.AppID).
		Where("msgLocalId IN (?)", fwdLocalIDs).
		Find(&messagesMap).Error; err != nil {
		return model.MessageExt{}, nil, err
	}
	fwdLocalToID := make(map[string]int, len(messagesMap))
	for _, pair := range messagesMap {
		fwdLocalToID[pair.MsgLocalID] = pair.MsgID
	}
	forwards := make([]ForwardExt, len(e.Forwards))
	for i, forward := range e.Forwards {
		forwards[i] = ForwardExt{
			ID:          fwdLocalToID[forward.LocalID],
			LocalID:     forward.LocalID,
			Username:    forward.Username,
			Text:        forward.Text,
			CreateDate:  forward.CreateDate,
			Attachments: make([]*model.Attachment, len(e.Forwards)),
		}

		// Fill attachments of forwards
		for j, attach := range forward.Attachments {
			forwards[i].Attachments[j] = &model.Attachment{
				LocalID:    attach.LocalID,
				AppID:      e.AppID,
				Url:        attach.Url,
				HasSpoiler: attach.HasSpoiler,
				Type:       attach.Type,
			}
		}
		forwardsAttachments = append(forwardsAttachments, forwards[i].Attachments...)
	}

	if len(msg.Attachments) != 0 {
		// Save attachments
		if err := c.DB().Create(msg.Attachments).Error; err != nil {
			return model.MessageExt{}, nil, err
		}
	}

	return msg, forwards, nil
}

func CollectForwardExtAttachIDs(exts []ForwardExt) []int {
	var attachIDs []int
	for _, ext := range exts {
		attachIDs = append(attachIDs, model.CollectAttachmentID(ext.Attachments)...)
	}

	return attachIDs
}

func AttachIDsInWaitingUpload(c event.Context, ids []int) ([]int, error) {
	var waitingIDs []int
	if err := c.DB().
		Table(model.TableAttachments).
		Joins("LEFT JOIN "+model.TableFiles+" ON Attachments.id = Files.attachmentId").
		Where("Attachments.id IN (?)", ids).
		Where("Files.id IS NULL").
		Pluck("Attachments.id AS id", &waitingIDs).Error; err != nil {
		return nil, err
	}

	return waitingIDs, nil
}

func AttachmentToCbkAttachs(attachs []*model.Attachment, waitingIDs []int) []callback.AttachmentNew {
	newAttachments := make([]callback.AttachmentNew, 0, len(attachs))
	for _, attach := range attachs {
		newAttachments = append(newAttachments, callback.AttachmentNew{
			HasSpoiler: attach.HasSpoiler,
			Type:       attach.Type,
			Url:        attach.Url,
			WaitUpload: !slices.Contains(waitingIDs, attach.ID),
		})
	}

	return newAttachments
}

func AppIDToMsgMapByMsgID(c event.Context, msgID int) (map[int]model.MessageMap, error) {
	var msgMapSl []model.MessageMap
	if err := c.DB().
		Where("msgId = ?", msgID).
		Find(&msgMapSl).Error; err != nil {
		return nil, err
	}

	msgMap := make(map[int]model.MessageMap, len(msgMapSl))
	for _, e := range msgMapSl {
		msgMap[e.AppID] = e
	}

	return msgMap, nil
}

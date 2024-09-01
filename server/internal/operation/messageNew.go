package operation

import (
	"slices"

	"github.com/chat-merger/merger/server/internal/callback"
	"github.com/chat-merger/merger/server/internal/event"
	"github.com/chat-merger/merger/server/internal/model"
)

func MessageNew(c Context, e event.MessageNew) error {
	msg, forwards, err := SaveMsg(c, e)
	if err != nil {
		return err
	}

	attachIDs := CollectForwardExtAttachIDs(forwards)
	var attachWaitingIDs []int
	if attachWaitingIDs, err = AttachIDsInWaitingUpload(c, attachIDs); err != nil {
		return err
	}

	var replyByAppID map[int]model.MessageMap
	if replyByAppID, err = AppIDToMsgMapByMsgID(c, msg.Reply); err != nil {
		return err
	}

	callbackMsgMap := make(map[int]callback.NewMessage)

	var applications []*model.Application
	if applications, err = Applications(c, msg.AppID); err != nil {
		return err
	}

	for _, app := range applications {
		callbackMsgMap[app.ID] = callback.NewMessage{
			ID:          msg.ID,
			IsSilent:    e.IsSilent,
			Reply:       msg.Reply,
			ReplyLocal:  replyByAppID[app.ID].MsgLocalID,
			Username:    e.Username,
			Text:        e.Text,
			Forwards:    ForwardExtToCbkForwards(forwards, attachWaitingIDs),
			Attachments: AttachmentToCbkAttachs(msg.Attachments, attachWaitingIDs),
		}
	}

	c.CallbackApi().OnNewMsg(callbackMsgMap)

	return nil
}

func ForwardExtToCbkForwards(exts []ForwardExt, attachWaitingIDs []int) []callback.NewForward {
	newForwards := make([]callback.NewForward, 0, len(exts))
	for _, ext := range exts {
		newForwards = append(newForwards, callback.NewForward{
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

func SaveMsg(c Context, e event.MessageNew) (model.MessageExt, []ForwardExt, error) {
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

	fwdLocalIDs := event.CollectMessageNewForwardsLoclIDs(e.Forwards)
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
			Forward: Forward{
				ID:         fwdLocalToID[forward.LocalID],
				LocalID:    forward.LocalID,
				Username:   forward.Username,
				Text:       forward.Text,
				CreateDate: forward.CreateDate,
			},
		}

		// Fill attachments of forwards
		forwards[i].Attachments = make([]*model.Attachment, len(e.Forwards))
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

type Forward struct {
	ID          int
	LocalID     string
	Username    string
	Text        string
	CreateDate  string
	Attachments []model.Attachment
}

type ForwardExt struct {
	Forward
	Attachments []*model.Attachment
}

func CollectForwardExtAttachIDs(exts []ForwardExt) []int {
	var attachIDs []int
	for _, ext := range exts {
		attachIDs = append(attachIDs, model.CollectAttachmentID(ext.Attachments)...)
	}

	return attachIDs
}

func Applications(c Context, excludeIDs ...int) ([]*model.Application, error) {
	var apps []*model.Application
	if err := c.DB().
		Where("id NOT IN (?)", excludeIDs).
		Find(&apps).Error; err != nil {
		return nil, err
	}

	return apps, nil
}

func AttachIDsInWaitingUpload(c Context, ids []int) ([]int, error) {
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

func AttachmentToCbkAttachs(attachs []*model.Attachment, waitingIDs []int) []callback.NewAttachment {
	newAttachments := make([]callback.NewAttachment, 0, len(attachs))
	for _, attach := range attachs {
		newAttachments = append(newAttachments, callback.NewAttachment{
			HasSpoiler: attach.HasSpoiler,
			Type:       attach.Type,
			Url:        attach.Url,
			WaitUpload: !slices.Contains(waitingIDs, attach.ID),
		})
	}

	return newAttachments
}

func AppIDToMsgMapByMsgID(c Context, msgID int) (map[int]model.MessageMap, error) {
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

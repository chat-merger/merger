package internal

import (
	"slices"

	"github.com/chat-merger/merger/server/internal/callback"
	"github.com/chat-merger/merger/server/internal/event"
)

func EventNewMessage(a *App, e event.MessageNew) error {
	msg, forwards, err := SaveMsg(a, e)
	if err != nil {
		return err
	}

	attachIDs := CollectForwardExtAttachIDs(forwards)
	var attachWaitingIDs []int
	if attachWaitingIDs, err = AttachIDsInWaitingUpload(a, attachIDs); err != nil {
		return err
	}

	var replyByAppID map[int]MessageMap
	if replyByAppID, err = AppIDToMsgMapByMsgID(a, msg.Reply); err != nil {
		return err
	}

	callbackMsgMap := make(map[int]callback.NewMessage)

	var applications []*Application
	if applications, err = Applications(a, msg.AppID); err != nil {
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

	a.callbackApi.OnNewMsg(callbackMsgMap)

	return nil
}

func FileUpload(a *App, f event.FileUpload) {}

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

func attachAppIDToID() {}

func SaveMsg(a *App, e event.MessageNew) (MessageExt, []ForwardExt, error) {
	// Save message
	msg := MessageExt{Message: Message{AppID: e.AppID}}
	var reply struct {
		ID int `gorm:"column:msgId"`
	}
	if err := a.db.
		Model(MessageMap{}).
		Select("msgId").
		Where("appId = ?", e.AppID).
		Where("msgLocalId = ?", e.ReplyLocalID).
		Find(&reply).Error; err != nil {
		return MessageExt{}, nil, err
	}
	msg.Reply = reply.ID

	if err := a.db.Save(&msg).Error; err != nil {
		return MessageExt{}, nil, err
	}

	// Fill message attachments
	msg.Attachments = make([]*Attachment, len(e.Forwards))
	for i, attach := range e.Attachments {
		msg.Attachments[i] = &Attachment{
			LocalID:    attach.LocalID,
			AppID:      e.AppID,
			Url:        attach.Url,
			HasSpoiler: attach.HasSpoiler,
			Type:       attach.Type,
		}
	}

	var forwardsAttachments []*Attachment

	fwdLocalIDs := event.CollectMessageNewForwardsLoclIDs(e.Forwards)
	var messagesMap []MessageMap
	if err := a.db.
		Table(TableMessagesMap).
		Where("appId = ?", e.AppID).
		Where("msgLocalId IN (?)", fwdLocalIDs).
		Find(&messagesMap).Error; err != nil {
		return MessageExt{}, nil, err
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
		forwards[i].Attachments = make([]*Attachment, len(e.Forwards))
		for j, attach := range forward.Attachments {
			forwards[i].Attachments[j] = &Attachment{
				LocalID:    attach.LocalID,
				AppID:      e.AppID,
				Url:        attach.Url,
				HasSpoiler: attach.HasSpoiler,
				Type:       attach.Type,
			}
		}
		forwardsAttachments = append(forwardsAttachments, forwards[i].Attachments...)
	}

	// Save attachments
	if err := a.db.Create(msg.Attachments).Error; err != nil {
		return MessageExt{}, nil, err
	}

	return msg, forwards, nil
}

type Forward struct {
	ID          int
	LocalID     string
	Username    string
	Text        string
	CreateDate  string
	Attachments []Attachment
}

type ForwardExt struct {
	Forward
	Attachments []*Attachment
}

func CollectForwardExtAttachIDs(exts []ForwardExt) []int {
	var attachIDs []int
	for _, ext := range exts {
		attachIDs = append(attachIDs, CollectAttachmentID(ext.Attachments)...)
	}

	return attachIDs
}

func Applications(a *App, excludeIDs ...int) ([]*Application, error) {
	var apps []*Application
	if err := a.db.
		Where("id NOT IN (?)", excludeIDs).
		Find(&apps).Error; err != nil {
		return nil, err
	}

	return apps, nil
}

func AttachIDsInWaitingUpload(a *App, ids []int) ([]int, error) {
	var waitingIDs []int
	if err := a.db.
		Table(TableAttachments).
		Joins("LEFT JOIN "+TableFiles+" ON Attachments.id = File.attachmentId").
		Where("Attachments.id IN (?)", ids).
		Where("File.id IS NULL").
		Pluck("Attachments.id AS id", &waitingIDs).Error; err != nil {
		return nil, err
	}

	return waitingIDs, nil
}

func AttachmentToCbkAttachs(attachs []*Attachment, waitingIDs []int) []callback.NewAttachment {
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

func AppIDToMsgMapByMsgID(a *App, msgID int) (map[int]MessageMap, error) {
	var msgMapSl []MessageMap
	if err := a.db.
		Where("msgId = ?", msgID).
		Find(&msgMapSl).Error; err != nil {
		return nil, err
	}

	msgMap := make(map[int]MessageMap, len(msgMapSl))
	for _, e := range msgMapSl {
		msgMap[e.AppID] = e
	}

	return msgMap, nil
}

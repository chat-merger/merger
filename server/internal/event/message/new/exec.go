package new

import (
	"github.com/chat-merger/merger/server/internal/callback"
	"github.com/chat-merger/merger/server/internal/common"
	"github.com/chat-merger/merger/server/internal/event"
	"github.com/chat-merger/merger/server/internal/model"
)

func Exec(c event.Context, e Message) error {
	msg, forwards, err := SaveMessage(c, e)
	if err != nil {
		return err
	}

	attachIDs := CollectForwardExtAttachIDs(forwards)
	var attachWaitingIDs []int
	if attachWaitingIDs, err = AttachIDsInWaitingUpload(c, attachIDs); err != nil {
		return err
	}

	var replyByAppID map[int]model.Bind
	if replyByAppID, err = AppIDToMsgMapByMsgID(c, msg.Reply); err != nil {
		return err
	}

	var applications []*model.Application
	if applications, err = common.Applications(c.DB(), msg.AppID); err != nil {
		return err
	}

	callbackNewMessages := make([]callback.MessageNew, len(applications))
	for i, app := range applications {
		callbackNewMessages[i] = callback.MessageNew{
			App:         *app,
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

	var newBinds []*model.Bind
	if newBinds, err = c.CBClient().MessageNew(callbackNewMessages); err != nil {
		return err
	}

	if err = c.DB().
		Table(model.TableBinds).
		Create(&newBinds).Error; err != nil {
		return err
	}

	return nil
}

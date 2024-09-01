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

	var replyByAppID map[int]model.MessageMap
	if replyByAppID, err = AppIDToMsgMapByMsgID(c, msg.Reply); err != nil {
		return err
	}

	var applications []*model.Application
	if applications, err = common.Applications(c.DB(), msg.AppID); err != nil {
		return err
	}

	callbackNewMessages := make([]callback.MessageNew, 0, len(applications))
	for _, app := range applications {
		callbackNewMessages = append(callbackNewMessages, callback.MessageNew{
			App:         *app,
			ID:          msg.ID,
			IsSilent:    e.IsSilent,
			Reply:       msg.Reply,
			ReplyLocal:  replyByAppID[app.ID].MsgLocalID,
			Username:    e.Username,
			Text:        e.Text,
			Forwards:    ForwardExtToCbkForwards(forwards, attachWaitingIDs),
			Attachments: AttachmentToCbkAttachs(msg.Attachments, attachWaitingIDs),
		})
	}

	var responses []callback.MessageNewResponse
	if responses, err = c.CallbackApi().OnMessageNew(callbackNewMessages); err != nil {
		return err
	}

	messageMaps := NewMsgResponseToMsgMap(responses)
	if err = c.DB().
		Table(model.TableMessagesMap).
		Create(&messageMaps).Error; err != nil {
		return err
	}

	return nil
}

package internal

func EventNewMessage(a *App, e EventNewMsg) {
	msg, forwards := SaveMsg(a, e)
	replyByAppID := MsgMapByAppID(msg.Reply)
	callbackMsgMap := make(map[int]callbackNewMsg)
	for _, app := range Applications(a, msg.AppID) {
		callbackMsgMap[app.ID] = callbackNewMsg{
			ID:           msg.ID,
			IsSilent:     e.IsSilent,
			Reply:        msg.Reply,
			ReplyLocalID: replyByAppID[app.ID].MsgLID,
			Username:     e.Username,
			Text:         e.Text,
			Forwards:     MessageExtToCbkForwards(forwards),
			Attachments:  AttachmentToCbkAttachs(msg.Attachments),
		}
	}

	a.cbApi.OnNewMsg(callbackMsgMap)
}

func FileUpload(a *App, f File) {}

func MessageExtToCbkForwards(exts []MessageExt) []callbackNewMsgForward {}

func attachAppIDToID() {}

func SaveMsg(a *App, e EventNewMsg) (MessageExt, []MessageExt) {}

func Applications(a *App, excludeIDs ...int) []Application {}

func IDByInApp(id string) int {}
func InAppByID(id int) string {}

func IDAttachByInApp(id string) int       {}
func InAppAttachByID(id int) string       {}
func InAppAttachByIDs(ids []int) []string {}

func AttachmentToCbkAttachs(attachs []Attachment) []callbackNewMsgAttachment {}

func AttachmentsFilter(attachs []Attachment, ids []int) []Attachment {}

func MsgMapByAppID(msgID int) map[int]MessageMap {}

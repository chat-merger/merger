package internal

func EventNewMessage(a *App, e EventNewMsg) {
	msg, attachs := SaveMsg(a, e)
	replyByAppID := msgMapByAppID(msg.Reply)
	callbackMsgMap := make(map[int]callbackNewMsg)
	for i, app := range Apps(a) {
		callbackMsgMap[app.ID] = callbackNewMsg{
			ID:            msg.ID,
			IsSilent:      e.IsSilent,
			Reply:         msg.Reply,
			ReplyInAppID:  replyByAppID[app.ID].InAppID,
			Username:      e.Username,
			Text:          e.Text,
			Forwards:      ,
			Attachments:   AttachmentToCbkAttach(attachs...),
		}

	}

	a.cbApi.OnNewMsg(e, hosts)
	for i, app := range hosts {

	}
}

func attachAppIDToID() {

}

//func Files(a *App,ids []int) []File {
//
//}

func FileUpload(a *App, f File) {

}

func SaveMsg(a *App, e EventNewMsg) (Message, []Attachment, []Forward) {}

func Apps(a *App) []Application {}

func IDByInApp(id string) int {}
func InAppByID(id int) string {}

func IDAttachByInApp(id string) int {}
func InAppAttachByID(id int) string {}
func InAppAttachByIDs(ids []int) []string {}

func AttachmentToCbkAttach(attachs Attachment) []callbackNewMsgAttachment {

}

func AttachmentsFilter(attachs []Attachment, ids []int) []Attachment {

}

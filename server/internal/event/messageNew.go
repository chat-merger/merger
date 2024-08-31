package event

type MessageNew struct {
	AppID        int
	LocalID      string
	IsSilent     bool
	ReplyLocalID string
	Username     string
	Text         string
	Forwards     []ForwardNew
	Attachments  []AttachmentNew
}

func CollectMessageNewForwardsLoclIDs(fwds []ForwardNew) []string {
	result := make([]string, len(fwds))
	for i, fwd := range fwds {
		result[i] = fwd.LocalID
	}

	return result
}

type ForwardNew struct {
	LocalID     string
	Username    string
	Text        string
	CreateDate  string
	Attachments []*AttachmentNew
}

type AttachmentNew struct {
	LocalID    string
	HasSpoiler bool
	Type       int
	// Общедоступная ссылка для загрузки файла.
	// Если ссылка не передана и по такому FileID не найдено файлов, то клиент должен будет загрузить файлы на специальный эндпоинт
	Url string
}

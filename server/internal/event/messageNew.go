package event

type MessageNew struct {
	AppID        int             `json:"appID,omitempty"`
	LocalID      string          `json:"localID,omitempty"`
	IsSilent     bool            `json:"isSilent,omitempty"`
	ReplyLocalID string          `json:"replyLocalID,omitempty"`
	Username     string          `json:"username,omitempty"`
	Text         string          `json:"text,omitempty"`
	Forwards     []ForwardNew    `json:"forwards,omitempty"`
	Attachments  []AttachmentNew `json:"attachments,omitempty"`
}

func CollectMessageNewForwardsLoclIDs(fwds []ForwardNew) []string {
	result := make([]string, len(fwds))
	for i, fwd := range fwds {
		result[i] = fwd.LocalID
	}

	return result
}

type ForwardNew struct {
	LocalID     string           `json:"localID,omitempty"`
	Username    string           `json:"username,omitempty"`
	Text        string           `json:"text,omitempty"`
	CreateDate  string           `json:"createDate,omitempty"`
	Attachments []*AttachmentNew `json:"attachments,omitempty"`
}

type AttachmentNew struct {
	LocalID    string `json:"localID,omitempty"`
	HasSpoiler bool   `json:"hasSpoiler,omitempty"`
	Type       int    `json:"type,omitempty"`
	// Общедоступная ссылка для загрузки файла.
	// Если ссылка не передана и по такому FileID не найдено файлов, то клиент должен будет загрузить файлы на специальный эндпоинт
	Url string `json:"url,omitempty"`
}

package new

type Message struct {
	AppID        int          `json:"app_id,omitempty"`
	LocalID      string       `json:"local_id,omitempty"`
	IsSilent     bool         `json:"is_silent,omitempty"`
	ReplyLocalID string       `json:"reply_local_id,omitempty"`
	Username     string       `json:"username,omitempty"`
	Text         string       `json:"text,omitempty"`
	Forwards     []Forward    `json:"forwards,omitempty"`
	Attachments  []Attachment `json:"attachments,omitempty"`
}

func CollectMessageNewForwardsLocalIDs(fwds []Forward) []string {
	result := make([]string, len(fwds))
	for i, fwd := range fwds {
		result[i] = fwd.LocalID
	}

	return result
}

type Forward struct {
	LocalID     string        `json:"local_id,omitempty"`
	Username    string        `json:"username,omitempty"`
	Text        string        `json:"text,omitempty"`
	CreateDate  string        `json:"create_date,omitempty"`
	Attachments []*Attachment `json:"attachments,omitempty"`
}

type Attachment struct {
	LocalID    string `json:"local_id,omitempty"`
	HasSpoiler bool   `json:"has_spoiler,omitempty"`
	Type       int    `json:"type,omitempty"`
	// Общедоступная ссылка для загрузки файла.
	// Если ссылка не передана и по такому FileID не найдено файлов, то клиент должен будет загрузить файлы на специальный эндпоинт
	Url string `json:"url,omitempty"`
}

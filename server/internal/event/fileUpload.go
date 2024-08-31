package event

type FileUpload struct {
	AppID   int
	Bytes   []byte
	Type    int
	LocalID string
}

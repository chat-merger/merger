package callback

import (
	"net/http"
)

type API interface {
	OnMessageNew(c []MessageNew) ([]MessageNewResponse, error)
}

func NewAPI() API { return &api{cl: http.DefaultClient} }

type api struct {
	cl *http.Client
}

type Body struct {
	MessageNew *MessageNew `json:"message_new,omitempty"`
}

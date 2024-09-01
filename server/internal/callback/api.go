package callback

import (
	"net/http"

	"github.com/chat-merger/merger/server/internal/model"
)

type Client interface {
	MessageNew(c []MessageNew) ([]*model.Bind, error)
}

func NewClient() Client { return &client{cl: http.DefaultClient} }

type client struct {
	cl *http.Client
}

type Body struct {
	MessageNew *MessageNew `json:"message_new,omitempty"`
}

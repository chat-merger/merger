package callback

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type API interface {
	OnNewMsg(c []NewMessage) ([]NewMsgResponse, error)
}

type api struct {
	cl *http.Client
}

func (a *api) OnNewMsg(newMessages []NewMessage) ([]NewMsgResponse, error) {
	result := make([]NewMsgResponse, len(newMessages))
	for i, newMsg := range newMessages {
		b, err := json.Marshal(newMsg)
		if err != nil {
			return nil, fmt.Errorf("json.Marshal: %w ", err)
		}
		r, err := a.cl.Post(newMsg.App.Host, "application/json", bytes.NewBuffer(b))
		if err != nil {
			return nil, fmt.Errorf("a.cl.Post: %w", err)
		}
		var resp NewMsgResponse
		if err = json.NewDecoder(r.Body).Decode(&resp); err != nil {
			return nil, fmt.Errorf("json.NewDecoder.Decode: %w", err)
		}
		resp.MsgID = newMsg.ID
		resp.AppID = newMsg.App.ID

		result[i] = resp
	}

	return result, nil
}

func NewAPI() API {
	return &api{
		cl: http.DefaultClient,
	}
}

package callback

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type API interface {
	OnNewMsg(c map[int]NewMessage) NewMsgResponse
}

type api struct {
	cl *http.Client
}

func (a *api) OnNewMsg(c map[int]NewMessage) NewMsgResponse {
	b, err := json.Marshal(c)
	if err != nil {
		fmt.Println(err.Error())
	}
	_, err = a.cl.Post("http://localhost:43687/echo", "application/json", bytes.NewBuffer(b))
	if err != nil {
		fmt.Println(err.Error())
	}
	return NewMsgResponse{LocalID: "x"}
}

func NewAPI() API {
	return &api{
		cl: http.DefaultClient,
	}
}

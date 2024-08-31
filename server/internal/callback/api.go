package callback

type API interface {
	OnNewMsg(c map[int]NewMessage) NewMsgResponse
}

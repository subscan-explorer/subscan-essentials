package pusher

type Stat interface {
	Push(msgType, text string, extra ...string) error
}

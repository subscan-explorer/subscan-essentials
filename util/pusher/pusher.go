package pusher

// Alert interface
type Alert interface {
	Push(scene, text string, extra ...string) error
}

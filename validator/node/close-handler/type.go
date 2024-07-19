package closehandler

// CloseHandler is a struct used to manage the graceful shutdown of a node.
// It contains a function to trigger the closing process and a channel to
// wait for the completion of the shutdown.
type CloseHandler struct {
	CloseFunc func()
	CloseCh   chan struct{}
}

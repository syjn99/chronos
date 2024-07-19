package rpc

import (
	"net/http"
)

// CloseClient closes validator client gracefully.
func (s *Server) CloseClient(w http.ResponseWriter, r *http.Request) {
	log.Infof("Got interrupt through HTTP request...")

	closeFunc := s.closeHandler.CloseFunc
	closeCh := s.closeHandler.CloseCh

	go closeFunc()

	// Wait for close channel to be closed.
	<-closeCh
}

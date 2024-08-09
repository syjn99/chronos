package overnode

import (
	"net/http"

	log "github.com/sirupsen/logrus"
)

// CloseClient closes beacon chain client gracefully.
func (s *Server) CloseClient(w http.ResponseWriter, r *http.Request) {
	log.Infof("Got interrupt through HTTP request...")

	closeFunc := s.CloseHandler.CloseFunc
	closeCh := s.CloseHandler.CloseCh

	go closeFunc()

	// Wait for close channel to be closed.
	<-closeCh
}

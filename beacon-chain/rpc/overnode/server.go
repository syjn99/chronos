package overnode

import (
	closehandler "github.com/prysmaticlabs/prysm/v5/beacon-chain/node/close-handler"
)

type Server struct {
	CloseHandler *closehandler.CloseHandler
}

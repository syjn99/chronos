package over

import (
	"github.com/prysmaticlabs/prysm/v5/beacon-chain/blockchain"
	"github.com/prysmaticlabs/prysm/v5/beacon-chain/rpc/lookup"
)

// Server defines a server implementation of custom APIs for OverProtocol.
type Server struct {
	Stater                lookup.Stater
	GenesisTimeFetcher    blockchain.TimeFetcher
	HeadFetcher           blockchain.HeadFetcher
	OptimisticModeFetcher blockchain.OptimisticModeFetcher
	FinalizationFetcher   blockchain.FinalizationFetcher
}

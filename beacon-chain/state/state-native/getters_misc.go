package state_native

import (
	"github.com/prysmaticlabs/prysm/v4/consensus-types/primitives"
	ethpb "github.com/prysmaticlabs/prysm/v4/proto/prysm/v1alpha1"
	"github.com/prysmaticlabs/prysm/v4/runtime/version"
)

// GenesisTime of the beacon state as a uint64.
func (b *BeaconState) GenesisTime() uint64 {
	b.lock.RLock()
	defer b.lock.RUnlock()

	return b.genesisTime
}

// GenesisValidatorsRoot of the beacon state.
func (b *BeaconState) GenesisValidatorsRoot() []byte {
	b.lock.RLock()
	defer b.lock.RUnlock()

	return b.genesisValidatorsRoot[:]
}

// Version of the beacon state. This method
// is strictly meant to be used without a lock
// internally.
func (b *BeaconState) Version() int {
	return b.version
}

// Slot of the current beacon chain state.
func (b *BeaconState) Slot() primitives.Slot {
	b.lock.RLock()
	defer b.lock.RUnlock()

	return b.slot
}

// Fork version of the beacon chain.
func (b *BeaconState) Fork() *ethpb.Fork {
	if b.fork == nil {
		return nil
	}

	b.lock.RLock()
	defer b.lock.RUnlock()

	return b.forkVal()
}

// forkVal version of the beacon chain.
// This assumes that a lock is already held on BeaconState.
func (b *BeaconState) forkVal() *ethpb.Fork {
	if b.fork == nil {
		return nil
	}

	prevVersion := make([]byte, len(b.fork.PreviousVersion))
	copy(prevVersion, b.fork.PreviousVersion)
	currVersion := make([]byte, len(b.fork.CurrentVersion))
	copy(currVersion, b.fork.CurrentVersion)
	return &ethpb.Fork{
		PreviousVersion: prevVersion,
		CurrentVersion:  currVersion,
		Epoch:           b.fork.Epoch,
	}
}

// HistoricalRoots based on epochs stored in the beacon state.
func (b *BeaconState) HistoricalRoots() ([][]byte, error) {
	if b.historicalRoots == nil {
		return nil, nil
	}

	b.lock.RLock()
	defer b.lock.RUnlock()

	return b.historicalRoots.Slice(), nil
}

// Reserve of the beacon state as an uint64.
func (b *BeaconState) RewardAdjustmentFactor() uint64 {
	b.lock.RLock()
	defer b.lock.RUnlock()

	return b.rewardAdjustmentFactor
}

// balancesLength returns the length of the balances slice.
// This assumes that a lock is already held on BeaconState.
func (b *BeaconState) balancesLength() int {
	if b.balances == nil {
		return 0
	}

	return len(b.balances)
}

// PreviousEpochReserve of the beacon state as an uint64.
func (b *BeaconState) PreviousEpochReserve() uint64 {
	b.lock.RLock()
	defer b.lock.RUnlock()

	return b.previousEpochReserve
}

// CurrentEpochReserve of the beacon state as an uint64.
func (b *BeaconState) CurrentEpochReserve() uint64 {
	b.lock.RLock()
	defer b.lock.RUnlock()

	return b.currentEpochReserve
}

// HistoricalSummaries of the beacon state.
func (b *BeaconState) HistoricalSummaries() ([]*ethpb.HistoricalSummary, error) {
	if b.version < version.Capella {
		return nil, errNotSupported("HistoricalSummaries", b.version)
	}

	if b.historicalSummaries == nil {
		return nil, nil
	}

	b.lock.RLock()
	defer b.lock.RUnlock()

	return b.historicalSummariesVal(), nil
}

// historicalSummariesVal of the beacon state.
// This assumes that a lock is already held on BeaconState.
func (b *BeaconState) historicalSummariesVal() []*ethpb.HistoricalSummary {
	return ethpb.CopyHistoricalSummaries(b.historicalSummaries)
}

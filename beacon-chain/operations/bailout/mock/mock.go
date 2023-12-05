package mock

import (
	"github.com/prysmaticlabs/prysm/v4/beacon-chain/state"
	"github.com/prysmaticlabs/prysm/v4/consensus-types/primitives"
	eth "github.com/prysmaticlabs/prysm/v4/proto/prysm/v1alpha1"
)

// PoolMock is a fake implementation of PoolManager.
type PoolMock struct {
	Exits       []*eth.BailOut
	initialized bool
}

// Initialize --
func (m *PoolMock) Initialize(_ state.ReadOnlyBeaconState) {
	m.Exits = []*eth.BailOut{}
}

// ExitsForInclusion --
func (m *PoolMock) IsInitialized() bool {
	return m.initialized
}

// ExitsForInclusion --
func (m *PoolMock) BailoutsForInclusion(_ state.ReadOnlyBeaconState, _ primitives.Slot) ([]*eth.BailOut, error) {
	return m.Exits, nil
}

// InsertVoluntaryExit --
func (m *PoolMock) UpdateBailOuts(_ state.ReadOnlyBeaconState) {
	m.Exits = append(m.Exits, &eth.BailOut{})
}

// MarkIncluded --
func (*PoolMock) MarkIncluded(_ *eth.SignedVoluntaryExit) {
	panic("implement me")
}

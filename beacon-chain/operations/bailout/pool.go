package bailout

import (
	"math"
	"sort"
	"sync"

	"github.com/prysmaticlabs/prysm/v5/beacon-chain/core/blocks"
	"github.com/prysmaticlabs/prysm/v5/beacon-chain/state"
	"github.com/prysmaticlabs/prysm/v5/config/params"
	types "github.com/prysmaticlabs/prysm/v5/consensus-types/primitives"
	doublylinkedlist "github.com/prysmaticlabs/prysm/v5/container/doubly-linked-list"
	ethpb "github.com/prysmaticlabs/prysm/v5/proto/prysm/v1alpha1"
	"github.com/sirupsen/logrus"
)

// PoolManager maintains pending bail out queue.
// This pool is used by proposers to insert bail outs into new blocks.
type PoolManager interface {
	Initialize(state state.ReadOnlyBeaconState)
	IsInitialized() bool
	PendingBailOuts() ([]*ethpb.BailOut, error)
	BailoutsForInclusion(state state.ReadOnlyBeaconState, voluntaryLen int) ([]*ethpb.BailOut, error)
	UpdateBailOuts(state state.ReadOnlyBeaconState)
	MarkIncluded(exit *ethpb.BailOut)
}

// Pool is a concrete implementation of PoolManager.
type Pool struct {
	lock        sync.RWMutex
	initialized bool
	pending     doublylinkedlist.List[*ethpb.BailOut]
	m           map[types.ValidatorIndex]*doublylinkedlist.Node[*ethpb.BailOut]
}

// NewPool returns an initialized pool.
func NewPool() *Pool {
	return &Pool{
		initialized: false,
		pending:     doublylinkedlist.List[*ethpb.BailOut]{},
		m:           make(map[types.ValidatorIndex]*doublylinkedlist.Node[*ethpb.BailOut]),
	}
}

// Initialize initializes the pool with the current state.
func (p *Pool) Initialize(state state.ReadOnlyBeaconState) {
	if p.initialized {
		return
	}

	bailoutScores, err := state.BailOutScores()
	if err != nil {
		logrus.WithError(err).Error("could not get bail out scores from state")
		return
	}

	// Iterate to find threshold exceeded validators and sort it by bail out score.
	var filteredScores [][]uint64
	for i, s := range bailoutScores {
		if s >= params.BeaconConfig().BailOutScoreThreshold && s < math.MaxUint64 {
			filteredScores = append(filteredScores, []uint64{s, uint64(i)})
		}
	}

	if len(filteredScores) == 0 {
		p.initialized = true
		return
	}

	sort.Slice(filteredScores, func(i, j int) bool {
		return filteredScores[i][0] > filteredScores[j][0]
	})

	for _, s := range filteredScores {
		p.insertBailOut(&ethpb.BailOut{
			ValidatorIndex: types.ValidatorIndex(s[1]),
		})
	}
	p.initialized = true
}

// IsInitialized check if pool is initialized or not.
func (p *Pool) IsInitialized() bool {
	p.lock.Lock()
	defer p.lock.Unlock()
	return p.initialized
}

// PendingBailOuts returns all objects from the pool.
func (p *Pool) PendingBailOuts() ([]*ethpb.BailOut, error) {
	p.lock.RLock()
	defer p.lock.RUnlock()

	result := make([]*ethpb.BailOut, p.pending.Len())
	node := p.pending.First()
	var err error
	for i := 0; node != nil; i++ {
		result[i], err = node.Value()
		if err != nil {
			return nil, err
		}
		node, err = node.Next()
		if err != nil {
			return nil, err
		}
	}
	return result, nil
}

// BailoutsForInclusion returns objects that are ready for inclusion at the given slot. This method will not
// return more than the block enforced MaxVoluntaryExits.
func (p *Pool) BailoutsForInclusion(state state.ReadOnlyBeaconState, voluntaryLen int) ([]*ethpb.BailOut, error) {
	p.lock.RLock()
	length := int(math.Min(float64(params.BeaconConfig().MaxVoluntaryExits-uint64(voluntaryLen)), float64(p.pending.Len())))
	result := make([]*ethpb.BailOut, 0, length)
	node := p.pending.First()
	for node != nil && len(result) < length {
		exit, err := node.Value()
		if err != nil {
			p.lock.RUnlock()
			return nil, err
		}

		if err = blocks.VerifyBailOut(state, exit.ValidatorIndex); err != nil {
			p.lock.RUnlock()
			// MarkIncluded removes the invalid exit from the pool
			p.MarkIncluded(exit)
			p.lock.RLock()
		} else {
			result = append(result, exit)
		}
		node, err = node.Next()
		if err != nil {
			p.lock.RUnlock()
			return nil, err
		}
	}
	p.lock.RUnlock()
	return result, nil
}

// UpdateBailOuts updates the initialized pool with the current state.
func (p *Pool) UpdateBailOuts(state state.ReadOnlyBeaconState) {
	if !p.initialized {
		return
	}

	bailoutScores, err := state.BailOutScores()
	if err != nil {
		logrus.WithError(err).Error("could not get bail out scores from state")
		return
	}

	// Iterate to find threshold exceeded validators and add it to the pool.
	for i, s := range bailoutScores {
		if s >= params.BeaconConfig().BailOutScoreThreshold && s < params.BeaconConfig().BailOutScoreThreshold+params.BeaconConfig().BailOutScoreBias {
			p.insertBailOut(&ethpb.BailOut{
				ValidatorIndex: types.ValidatorIndex(i),
			})
		}
	}
}

// insertVoluntaryExit into the pool.
func (p *Pool) insertBailOut(exit *ethpb.BailOut) {
	p.lock.Lock()
	defer p.lock.Unlock()

	_, exists := p.m[exit.ValidatorIndex]
	if exists {
		return
	}

	p.pending.Append(doublylinkedlist.NewNode(exit))
	p.m[exit.ValidatorIndex] = p.pending.Last()
}

// MarkIncluded is used when an exit has been included in a beacon block. Every block seen by this
// node should call this method to include the exit. This will remove the exit from the pool.
func (p *Pool) MarkIncluded(exit *ethpb.BailOut) {
	p.lock.Lock()
	defer p.lock.Unlock()

	node := p.m[exit.ValidatorIndex]
	if node == nil {
		return
	}

	delete(p.m, exit.ValidatorIndex)
	p.pending.Remove(node)
}

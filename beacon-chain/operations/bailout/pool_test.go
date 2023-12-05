package bailout

import (
	"testing"

	state_native "github.com/prysmaticlabs/prysm/v4/beacon-chain/state/state-native"
	"github.com/prysmaticlabs/prysm/v4/config/params"
	types "github.com/prysmaticlabs/prysm/v4/consensus-types/primitives"
	"github.com/prysmaticlabs/prysm/v4/crypto/bls"
	"github.com/prysmaticlabs/prysm/v4/crypto/bls/common"
	"github.com/prysmaticlabs/prysm/v4/math"
	ethpb "github.com/prysmaticlabs/prysm/v4/proto/prysm/v1alpha1"
	"github.com/prysmaticlabs/prysm/v4/testing/assert"
	"github.com/prysmaticlabs/prysm/v4/testing/require"
)

func TestInitializeAndInitialized(t *testing.T) {
	spb := &ethpb.BeaconStateCapella{
		Fork: &ethpb.Fork{
			CurrentVersion:  params.BeaconConfig().GenesisForkVersion,
			PreviousVersion: params.BeaconConfig().GenesisForkVersion,
		},
	}
	stateSlot := types.Slot(uint64(params.BeaconConfig().ShardCommitteePeriod) * uint64(params.BeaconConfig().SlotsPerEpoch))
	spb.Slot = stateSlot
	numValidators := 2 * params.BeaconConfig().MaxVoluntaryExits //32
	validators := make([]*ethpb.Validator, numValidators)
	bailouts := make([]*ethpb.BailOut, numValidators)
	bailoutScores := make([]uint64, numValidators)
	privKeys := make([]common.SecretKey, numValidators)

	for i := range validators {
		v := &ethpb.Validator{}
		if i == len(validators)-2 {
			// exit for this validator is invalid
			v.ExitEpoch = 0
		} else {
			v.ExitEpoch = params.BeaconConfig().FarFutureEpoch
		}
		priv, err := bls.RandKey()
		require.NoError(t, err)
		privKeys[i] = priv
		pubkey := priv.PublicKey().Marshal()
		v.PublicKey = pubkey

		message := &ethpb.BailOut{
			ValidatorIndex: types.ValidatorIndex(i),
		}

		validators[i] = v
		bailouts[i] = message
		bailoutScores[i] = uint64(i)
	}
	spb.Validators = validators
	st, err := state_native.InitializeFromProtoCapella(spb)
	require.NoError(t, err)
	err = st.SetBailOutScores(bailoutScores)
	require.NoError(t, err)

	t.Run("no score above threshold", func(t *testing.T) {
		pool := NewPool()
		pool.Initialize(st)
		init, err := pool.PendingBailOuts()
		require.NoError(t, err)
		assert.Equal(t, 0, len(init))
	})

	for i := range validators {
		bailoutScores[i] += params.BeaconConfig().BailOutScoreThreshold
	}
	bailoutScores[len(bailoutScores)-2] = math.MaxUint64
	err = st.SetBailOutScores(bailoutScores)
	require.NoError(t, err)
	t.Run("some scores are already bail outed", func(t *testing.T) {
		pool := NewPool()
		pool.Initialize(st)
		init, err := pool.PendingBailOuts()
		require.NoError(t, err)
		assert.Equal(t, len(validators)-1, len(init))
	})
	t.Run("pool already initialized", func(t *testing.T) {
		pool := NewPool()
		pool.Initialize(st)
		initialized := pool.IsInitialized()
		assert.Equal(t, true, initialized)
	})
}

func TestPendingBailOuts(t *testing.T) {
	t.Run("empty pool", func(t *testing.T) {
		pool := NewPool()
		changes, err := pool.PendingBailOuts()
		require.NoError(t, err)
		assert.Equal(t, 0, len(changes))
	})
	t.Run("non-empty pool", func(t *testing.T) {
		pool := NewPool()
		pool.insertBailOut(&ethpb.BailOut{
			ValidatorIndex: 0,
		})
		pool.insertBailOut(&ethpb.BailOut{
			ValidatorIndex: 1,
		})
		changes, err := pool.PendingBailOuts()
		require.NoError(t, err)
		assert.Equal(t, 2, len(changes))
	})
}

func TestBailoutsForInclusion(t *testing.T) {
	spb := &ethpb.BeaconStateCapella{
		Fork: &ethpb.Fork{
			CurrentVersion:  params.BeaconConfig().GenesisForkVersion,
			PreviousVersion: params.BeaconConfig().GenesisForkVersion,
		},
	}
	stateSlot := types.Slot(uint64(params.BeaconConfig().ShardCommitteePeriod) * uint64(params.BeaconConfig().SlotsPerEpoch))
	spb.Slot = stateSlot
	numValidators := 2 * params.BeaconConfig().MaxVoluntaryExits
	validators := make([]*ethpb.Validator, numValidators)
	bailouts := make([]*ethpb.BailOut, numValidators)
	bailoutScores := make([]uint64, numValidators)
	privKeys := make([]common.SecretKey, numValidators)

	for i := range validators {
		v := &ethpb.Validator{}
		if i == len(validators)-2 {
			// exit for this validator is invalid
			v.ExitEpoch = 0
		} else {
			v.ExitEpoch = params.BeaconConfig().FarFutureEpoch
		}
		priv, err := bls.RandKey()
		require.NoError(t, err)
		privKeys[i] = priv
		pubkey := priv.PublicKey().Marshal()
		v.PublicKey = pubkey

		message := &ethpb.BailOut{
			ValidatorIndex: types.ValidatorIndex(i),
		}

		validators[i] = v
		bailouts[i] = message
		bailoutScores[i] = params.BeaconConfig().BailOutScoreThreshold
	}
	spb.Validators = validators
	st, err := state_native.InitializeFromProtoCapella(spb)
	require.NoError(t, err)
	err = st.SetBailOutScores(bailoutScores)
	require.NoError(t, err)

	t.Run("empty pool", func(t *testing.T) {
		pool := NewPool()
		exits, err := pool.BailoutsForInclusion(st, 10)
		require.NoError(t, err)
		assert.Equal(t, 0, len(exits))
	})
	t.Run("voluntaryExits less than MaxVoluntaryExits in pool", func(t *testing.T) {
		pool := NewPool()
		for i := uint64(0); i < params.BeaconConfig().MaxVoluntaryExits-1; i++ {
			pool.insertBailOut(bailouts[i])
		}
		exits, err := pool.BailoutsForInclusion(st, 14)
		require.NoError(t, err)
		assert.Equal(t, int(params.BeaconConfig().MaxVoluntaryExits)-14, len(exits))
	})
	t.Run("MaxVoluntaryExits in pool", func(t *testing.T) {
		pool := NewPool()
		for i := uint64(0); i < params.BeaconConfig().MaxVoluntaryExits; i++ {
			pool.insertBailOut(bailouts[i])
		}
		exits, err := pool.BailoutsForInclusion(st, 16)
		require.NoError(t, err)
		assert.Equal(t, 0, len(exits))
	})
	t.Run("more than MaxVoluntaryExits in pool", func(t *testing.T) {
		pool := NewPool()
		for i := uint64(0); i < numValidators; i++ {
			pool.insertBailOut(bailouts[i])
		}
		exits, err := pool.BailoutsForInclusion(st, 0)
		require.NoError(t, err)
		assert.Equal(t, int(params.BeaconConfig().MaxVoluntaryExits), len(exits))
		for _, ch := range exits {
			assert.NotEqual(t, types.ValidatorIndex(params.BeaconConfig().MaxVoluntaryExits), ch.ValidatorIndex)
		}
	})
	t.Run("invalid exit not returned", func(t *testing.T) {
		pool := NewPool()
		pool.insertBailOut(bailouts[len(bailouts)-2])
		exits, err := pool.BailoutsForInclusion(st, 15)
		require.NoError(t, err)
		assert.Equal(t, 0, len(exits))
	})
}

func TestInsertBailOut(t *testing.T) {
	t.Run("empty pool", func(t *testing.T) {
		pool := NewPool()
		exit := &ethpb.BailOut{
			ValidatorIndex: 0,
		}
		pool.insertBailOut(exit)
		require.Equal(t, 1, pool.pending.Len())
		require.Equal(t, 1, len(pool.m))
		n, ok := pool.m[0]
		require.Equal(t, true, ok)
		v, err := n.Value()
		require.NoError(t, err)
		assert.DeepEqual(t, exit, v)
	})
	t.Run("item in pool", func(t *testing.T) {
		pool := NewPool()
		old := &ethpb.BailOut{
			ValidatorIndex: 0,
		}
		exit := &ethpb.BailOut{
			ValidatorIndex: 1,
		}
		pool.insertBailOut(old)
		pool.insertBailOut(exit)
		require.Equal(t, 2, pool.pending.Len())
		require.Equal(t, 2, len(pool.m))
		n, ok := pool.m[0]
		require.Equal(t, true, ok)
		v, err := n.Value()
		require.NoError(t, err)
		assert.DeepEqual(t, old, v)
		n, ok = pool.m[1]
		require.Equal(t, true, ok)
		v, err = n.Value()
		require.NoError(t, err)
		assert.DeepEqual(t, exit, v)
	})
	t.Run("validator index already exists", func(t *testing.T) {
		pool := NewPool()
		old := &ethpb.BailOut{
			ValidatorIndex: 0,
		}
		exit := &ethpb.BailOut{
			ValidatorIndex: 0,
		}
		pool.insertBailOut(old)
		pool.insertBailOut(exit)
		assert.Equal(t, 1, pool.pending.Len())
		require.Equal(t, 1, len(pool.m))
		n, ok := pool.m[0]
		require.Equal(t, true, ok)
		v, err := n.Value()
		require.NoError(t, err)
		assert.DeepEqual(t, old, v)
	})
}

func TestUpdateBailOuts(t *testing.T) {
	spb := &ethpb.BeaconStateCapella{
		Fork: &ethpb.Fork{
			CurrentVersion:  params.BeaconConfig().GenesisForkVersion,
			PreviousVersion: params.BeaconConfig().GenesisForkVersion,
		},
	}
	stateSlot := types.Slot(uint64(params.BeaconConfig().ShardCommitteePeriod) * uint64(params.BeaconConfig().SlotsPerEpoch))
	spb.Slot = stateSlot
	numValidators := 2 * params.BeaconConfig().MaxVoluntaryExits //32
	validators := make([]*ethpb.Validator, numValidators)
	bailouts := make([]*ethpb.BailOut, numValidators)
	bailoutScores := make([]uint64, numValidators)
	privKeys := make([]common.SecretKey, numValidators)

	for i := range validators {
		v := &ethpb.Validator{}
		if i == len(validators)-2 {
			// exit for this validator is invalid
			v.ExitEpoch = 0
		} else {
			v.ExitEpoch = params.BeaconConfig().FarFutureEpoch
		}
		priv, err := bls.RandKey()
		require.NoError(t, err)
		privKeys[i] = priv
		pubkey := priv.PublicKey().Marshal()
		v.PublicKey = pubkey

		message := &ethpb.BailOut{
			ValidatorIndex: types.ValidatorIndex(i),
		}

		validators[i] = v
		bailouts[i] = message
		bailoutScores[i] = params.BeaconConfig().BailOutScoreThreshold - uint64(len(validators)) + uint64(i) + 2
	}
	spb.Validators = validators
	st, err := state_native.InitializeFromProtoCapella(spb)
	require.NoError(t, err)
	err = st.SetBailOutScores(bailoutScores)
	require.NoError(t, err)

	t.Run("pool is not initialized", func(t *testing.T) {
		pool := NewPool()
		pool.UpdateBailOuts(st)
		init, err := pool.PendingBailOuts()
		require.NoError(t, err)
		assert.Equal(t, 0, len(init))
	})

	t.Run("two score in range for update", func(t *testing.T) {
		pool := NewPool()
		pool.Initialize(st)
		init, err := pool.PendingBailOuts()
		require.NoError(t, err)
		assert.Equal(t, 2, len(init))

		bailoutScores[0] = math.MaxUint64
		bailoutScores[1] = params.BeaconConfig().BailOutScoreThreshold + 1
		bailoutScores[2] = params.BeaconConfig().BailOutScoreThreshold + params.BeaconConfig().BailOutScoreBias - 1
		bailoutScores[3] = params.BeaconConfig().BailOutScoreThreshold + params.BeaconConfig().BailOutScoreBias
		bailoutScores[len(validators)-2] += params.BeaconConfig().BailOutScoreBias
		bailoutScores[len(validators)-1] += params.BeaconConfig().BailOutScoreBias
		err = st.SetBailOutScores(bailoutScores)
		require.NoError(t, err)

		pool.UpdateBailOuts(st)
		updated, err := pool.PendingBailOuts()
		require.NoError(t, err)
		assert.Equal(t, 4, len(updated))
	})
}

func TestMarkIncluded(t *testing.T) {
	t.Run("one element in pool", func(t *testing.T) {
		pool := NewPool()
		exit := &ethpb.BailOut{
			ValidatorIndex: 0,
		}
		pool.insertBailOut(exit)
		pool.MarkIncluded(exit)
		assert.Equal(t, 0, pool.pending.Len())
		_, ok := pool.m[0]
		assert.Equal(t, false, ok)
	})
	t.Run("first of multiple elements", func(t *testing.T) {
		pool := NewPool()
		first := &ethpb.BailOut{
			ValidatorIndex: 0,
		}
		second := &ethpb.BailOut{
			ValidatorIndex: 1,
		}
		third := &ethpb.BailOut{
			ValidatorIndex: 2,
		}
		pool.insertBailOut(first)
		pool.insertBailOut(second)
		pool.insertBailOut(third)
		pool.MarkIncluded(first)
		require.Equal(t, 2, pool.pending.Len())
		_, ok := pool.m[0]
		assert.Equal(t, false, ok)
	})
	t.Run("last of multiple elements", func(t *testing.T) {
		pool := NewPool()
		first := &ethpb.BailOut{
			ValidatorIndex: 0,
		}
		second := &ethpb.BailOut{
			ValidatorIndex: 1,
		}
		third := &ethpb.BailOut{
			ValidatorIndex: 2,
		}
		pool.insertBailOut(first)
		pool.insertBailOut(second)
		pool.insertBailOut(third)
		pool.MarkIncluded(third)
		require.Equal(t, 2, pool.pending.Len())
		_, ok := pool.m[2]
		assert.Equal(t, false, ok)
	})
	t.Run("in the middle of multiple elements", func(t *testing.T) {
		pool := NewPool()
		first := &ethpb.BailOut{
			ValidatorIndex: 0,
		}
		second := &ethpb.BailOut{
			ValidatorIndex: 1,
		}
		third := &ethpb.BailOut{
			ValidatorIndex: 2,
		}
		pool.insertBailOut(first)
		pool.insertBailOut(second)
		pool.insertBailOut(third)
		pool.MarkIncluded(second)
		require.Equal(t, 2, pool.pending.Len())
		_, ok := pool.m[1]
		assert.Equal(t, false, ok)
	})
	t.Run("not in pool", func(t *testing.T) {
		pool := NewPool()
		first := &ethpb.BailOut{
			ValidatorIndex: 0,
		}
		second := &ethpb.BailOut{
			ValidatorIndex: 1,
		}
		exit := &ethpb.BailOut{
			ValidatorIndex: 2,
		}
		pool.insertBailOut(first)
		pool.insertBailOut(second)
		pool.MarkIncluded(exit)
		require.Equal(t, 2, pool.pending.Len())
		_, ok := pool.m[0]
		require.Equal(t, true, ok)
		assert.NotNil(t, pool.m[0])
		_, ok = pool.m[1]
		require.Equal(t, true, ok)
		assert.NotNil(t, pool.m[1])
	})
}

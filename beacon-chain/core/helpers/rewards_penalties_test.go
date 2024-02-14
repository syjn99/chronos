package helpers

import (
	"math"
	"testing"

	"github.com/prysmaticlabs/prysm/v4/beacon-chain/core/time"
	state_native "github.com/prysmaticlabs/prysm/v4/beacon-chain/state/state-native"
	"github.com/prysmaticlabs/prysm/v4/config/params"
	"github.com/prysmaticlabs/prysm/v4/consensus-types/primitives"
	ethpb "github.com/prysmaticlabs/prysm/v4/proto/prysm/v1alpha1"
	"github.com/prysmaticlabs/prysm/v4/testing/assert"
	"github.com/prysmaticlabs/prysm/v4/testing/require"
)

func TestTotalBalance_OK(t *testing.T) {
	state, err := state_native.InitializeFromProtoPhase0(&ethpb.BeaconState{Validators: []*ethpb.Validator{
		{EffectiveBalance: 27 * 1e9},
		{EffectiveBalance: 28 * 1e9},
		{EffectiveBalance: 32 * 1e9},
		{EffectiveBalance: 40 * 1e9},
	}})
	require.NoError(t, err)

	balance := TotalBalance(state, []primitives.ValidatorIndex{0, 1, 2, 3})
	wanted := state.Validators()[0].EffectiveBalance + state.Validators()[1].EffectiveBalance +
		state.Validators()[2].EffectiveBalance + state.Validators()[3].EffectiveBalance
	assert.Equal(t, wanted, balance, "Incorrect TotalBalance")
}

func TestTotalBalance_ReturnsEffectiveBalanceIncrement(t *testing.T) {
	state, err := state_native.InitializeFromProtoPhase0(&ethpb.BeaconState{Validators: []*ethpb.Validator{}})
	require.NoError(t, err)

	balance := TotalBalance(state, []primitives.ValidatorIndex{})
	wanted := params.BeaconConfig().EffectiveBalanceIncrement
	assert.Equal(t, wanted, balance, "Incorrect TotalBalance")
}

func TestGetBalance_OK(t *testing.T) {
	tests := []struct {
		i uint64
		b []uint64
	}{
		{i: 0, b: []uint64{27 * 1e9, 28 * 1e9, 32 * 1e9}},
		{i: 1, b: []uint64{27 * 1e9, 28 * 1e9, 32 * 1e9}},
		{i: 2, b: []uint64{27 * 1e9, 28 * 1e9, 32 * 1e9}},
		{i: 0, b: []uint64{0, 0, 0}},
		{i: 2, b: []uint64{0, 0, 0}},
	}
	for _, test := range tests {
		state, err := state_native.InitializeFromProtoPhase0(&ethpb.BeaconState{Balances: test.b})
		require.NoError(t, err)
		assert.Equal(t, test.b[test.i], state.Balances()[test.i], "Incorrect Validator balance")
	}
}

func TestTotalActiveBalance(t *testing.T) {
	tests := []struct {
		vCount int
	}{
		{1},
		{10},
		{10000},
	}
	for _, test := range tests {
		validators := make([]*ethpb.Validator, 0)
		for i := 0; i < test.vCount; i++ {
			validators = append(validators, &ethpb.Validator{EffectiveBalance: params.BeaconConfig().MaxEffectiveBalance, ExitEpoch: 1})
		}
		state, err := state_native.InitializeFromProtoPhase0(&ethpb.BeaconState{Validators: validators})
		require.NoError(t, err)
		bal, err := TotalActiveBalance(state)
		require.NoError(t, err)
		require.Equal(t, uint64(test.vCount)*params.BeaconConfig().MaxEffectiveBalance, bal)
	}
}

func TestTotalActiveBal_ReturnMin(t *testing.T) {
	ClearCache()
	defer ClearCache()
	tests := []struct {
		vCount int
	}{
		{1},
		{10},
		{10000},
	}
	for _, test := range tests {
		validators := make([]*ethpb.Validator, 0)
		for i := 0; i < test.vCount; i++ {
			validators = append(validators, &ethpb.Validator{EffectiveBalance: 1, ExitEpoch: 1})
		}
		state, err := state_native.InitializeFromProtoPhase0(&ethpb.BeaconState{Validators: validators})
		require.NoError(t, err)
		bal, err := TotalActiveBalance(state)
		require.NoError(t, err)
		require.Equal(t, params.BeaconConfig().EffectiveBalanceIncrement, bal)
	}
}

func TestTotalActiveBalance_WithCache(t *testing.T) {
	ClearCache()
	defer ClearCache()
	tests := []struct {
		vCount    int
		wantCount int
	}{
		{vCount: 1, wantCount: 1},
		{vCount: 10, wantCount: 10},
		{vCount: 10000, wantCount: 10000},
	}
	for _, test := range tests {
		validators := make([]*ethpb.Validator, 0)
		for i := 0; i < test.vCount; i++ {
			validators = append(validators, &ethpb.Validator{EffectiveBalance: params.BeaconConfig().MaxEffectiveBalance, ExitEpoch: 1})
		}
		state, err := state_native.InitializeFromProtoPhase0(&ethpb.BeaconState{Validators: validators})
		require.NoError(t, err)
		bal, err := TotalActiveBalance(state)
		require.NoError(t, err)
		require.Equal(t, uint64(test.wantCount)*params.BeaconConfig().MaxEffectiveBalance, bal)
	}
}

func TestIncreaseBalance_OK(t *testing.T) {
	tests := []struct {
		i  primitives.ValidatorIndex
		b  []uint64
		nb uint64
		eb uint64
	}{
		{i: 0, b: []uint64{27 * 1e9, 28 * 1e9, 32 * 1e9}, nb: 1, eb: 27*1e9 + 1},
		{i: 1, b: []uint64{27 * 1e9, 28 * 1e9, 32 * 1e9}, nb: 0, eb: 28 * 1e9},
		{i: 2, b: []uint64{27 * 1e9, 28 * 1e9, 32 * 1e9}, nb: 33 * 1e9, eb: 65 * 1e9},
	}
	for _, test := range tests {
		state, err := state_native.InitializeFromProtoPhase0(&ethpb.BeaconState{
			Validators: []*ethpb.Validator{
				{EffectiveBalance: 4}, {EffectiveBalance: 4}, {EffectiveBalance: 4},
			},
			Balances: test.b,
		})
		require.NoError(t, err)
		require.NoError(t, IncreaseBalance(state, test.i, test.nb))
		assert.Equal(t, test.eb, state.Balances()[test.i], "Incorrect Validator balance")
	}
}

func TestDecreaseBalance_OK(t *testing.T) {
	tests := []struct {
		i  primitives.ValidatorIndex
		b  []uint64
		nb uint64
		eb uint64
	}{
		{i: 0, b: []uint64{2, 28 * 1e9, 32 * 1e9}, nb: 1, eb: 1},
		{i: 1, b: []uint64{27 * 1e9, 28 * 1e9, 32 * 1e9}, nb: 0, eb: 28 * 1e9},
		{i: 2, b: []uint64{27 * 1e9, 28 * 1e9, 1}, nb: 2, eb: 0},
		{i: 3, b: []uint64{27 * 1e9, 28 * 1e9, 1, 28 * 1e9}, nb: 28 * 1e9, eb: 0},
	}
	for _, test := range tests {
		state, err := state_native.InitializeFromProtoPhase0(&ethpb.BeaconState{
			Validators: []*ethpb.Validator{
				{EffectiveBalance: 4}, {EffectiveBalance: 4}, {EffectiveBalance: 4}, {EffectiveBalance: 3},
			},
			Balances: test.b,
		})
		require.NoError(t, err)
		require.NoError(t, DecreaseBalance(state, test.i, test.nb))
		assert.Equal(t, test.eb, state.Balances()[test.i], "Incorrect Validator balance")
	}
}

func TestFinalityDelay(t *testing.T) {
	base := buildState(params.BeaconConfig().SlotsPerEpoch*10, 1)
	base.FinalizedCheckpoint = &ethpb.Checkpoint{Epoch: 3}
	beaconState, err := state_native.InitializeFromProtoPhase0(base)
	require.NoError(t, err)
	prevEpoch := primitives.Epoch(0)
	finalizedEpoch := primitives.Epoch(0)
	// Set values for each test case
	setVal := func() {
		prevEpoch = time.PrevEpoch(beaconState)
		finalizedEpoch = beaconState.FinalizedCheckpointEpoch()
	}
	setVal()
	d := FinalityDelay(prevEpoch, finalizedEpoch)
	w := time.PrevEpoch(beaconState) - beaconState.FinalizedCheckpointEpoch()
	assert.Equal(t, w, d, "Did not get wanted finality delay")

	require.NoError(t, beaconState.SetFinalizedCheckpoint(&ethpb.Checkpoint{Epoch: 4}))
	setVal()
	d = FinalityDelay(prevEpoch, finalizedEpoch)
	w = time.PrevEpoch(beaconState) - beaconState.FinalizedCheckpointEpoch()
	assert.Equal(t, w, d, "Did not get wanted finality delay")

	require.NoError(t, beaconState.SetFinalizedCheckpoint(&ethpb.Checkpoint{Epoch: 5}))
	setVal()
	d = FinalityDelay(prevEpoch, finalizedEpoch)
	w = time.PrevEpoch(beaconState) - beaconState.FinalizedCheckpointEpoch()
	assert.Equal(t, w, d, "Did not get wanted finality delay")
}

func TestIsInInactivityLeak(t *testing.T) {
	base := buildState(params.BeaconConfig().SlotsPerEpoch*10, 1)
	base.FinalizedCheckpoint = &ethpb.Checkpoint{Epoch: 3}
	beaconState, err := state_native.InitializeFromProtoPhase0(base)
	require.NoError(t, err)
	prevEpoch := primitives.Epoch(0)
	finalizedEpoch := primitives.Epoch(0)
	// Set values for each test case
	setVal := func() {
		prevEpoch = time.PrevEpoch(beaconState)
		finalizedEpoch = beaconState.FinalizedCheckpointEpoch()
	}
	setVal()
	assert.Equal(t, true, IsInInactivityLeak(prevEpoch, finalizedEpoch), "Wanted inactivity leak true")
	require.NoError(t, beaconState.SetFinalizedCheckpoint(&ethpb.Checkpoint{Epoch: 4}))
	setVal()
	assert.Equal(t, true, IsInInactivityLeak(prevEpoch, finalizedEpoch), "Wanted inactivity leak true")
	require.NoError(t, beaconState.SetFinalizedCheckpoint(&ethpb.Checkpoint{Epoch: 5}))
	setVal()
	assert.Equal(t, false, IsInInactivityLeak(prevEpoch, finalizedEpoch), "Wanted inactivity leak false")
}

func buildState(slot primitives.Slot, validatorCount uint64) *ethpb.BeaconState {
	validators := make([]*ethpb.Validator, validatorCount)
	for i := 0; i < len(validators); i++ {
		validators[i] = &ethpb.Validator{
			ExitEpoch:        params.BeaconConfig().FarFutureEpoch,
			EffectiveBalance: params.BeaconConfig().MaxEffectiveBalance,
		}
	}
	validatorBalances := make([]uint64, len(validators))
	for i := 0; i < len(validatorBalances); i++ {
		validatorBalances[i] = params.BeaconConfig().MaxEffectiveBalance
	}
	latestActiveIndexRoots := make(
		[][]byte,
		params.BeaconConfig().EpochsPerHistoricalVector,
	)
	for i := 0; i < len(latestActiveIndexRoots); i++ {
		latestActiveIndexRoots[i] = params.BeaconConfig().ZeroHash[:]
	}
	latestRandaoMixes := make(
		[][]byte,
		params.BeaconConfig().EpochsPerHistoricalVector,
	)
	for i := 0; i < len(latestRandaoMixes); i++ {
		latestRandaoMixes[i] = params.BeaconConfig().ZeroHash[:]
	}
	return &ethpb.BeaconState{
		Slot:                        slot,
		Balances:                    validatorBalances,
		Validators:                  validators,
		RandaoMixes:                 make([][]byte, params.BeaconConfig().EpochsPerHistoricalVector),
		Slashings:                   make([]uint64, params.BeaconConfig().EpochsPerSlashingsVector),
		BlockRoots:                  make([][]byte, params.BeaconConfig().SlotsPerEpoch*10),
		FinalizedCheckpoint:         &ethpb.Checkpoint{Root: make([]byte, 32)},
		PreviousJustifiedCheckpoint: &ethpb.Checkpoint{Root: make([]byte, 32)},
		CurrentJustifiedCheckpoint:  &ethpb.Checkpoint{Root: make([]byte, 32)},
	}
}

func TestIncreaseBadBalance_NotOK(t *testing.T) {
	tests := []struct {
		i  primitives.ValidatorIndex
		b  []uint64
		nb uint64
	}{
		{i: 0, b: []uint64{math.MaxUint64, math.MaxUint64, math.MaxUint64}, nb: 1},
		{i: 2, b: []uint64{math.MaxUint64, math.MaxUint64, math.MaxUint64}, nb: 33 * 1e9},
	}
	for _, test := range tests {
		state, err := state_native.InitializeFromProtoPhase0(&ethpb.BeaconState{
			Validators: []*ethpb.Validator{
				{EffectiveBalance: 4}, {EffectiveBalance: 4}, {EffectiveBalance: 4},
			},
			Balances: test.b,
		})
		require.NoError(t, err)
		require.ErrorContains(t, "addition overflows", IncreaseBalance(state, test.i, test.nb))
	}
}

func TestEpochIssuance(t *testing.T) {
	tests := []struct {
		name string
		e    primitives.Epoch
		want uint64
	}{
		{name: "Issuance of Year 1", e: primitives.Epoch(0), want: 243531202435},
		{name: "Issuance of Year 2", e: primitives.Epoch(params.BeaconConfig().EpochsPerYear + 1), want: 669710806697},
		{name: "Issuance of Year 3", e: primitives.Epoch(params.BeaconConfig().EpochsPerYear*2 + 1), want: 913242009132},
		{name: "Issuance of Year 4", e: primitives.Epoch(params.BeaconConfig().EpochsPerYear*3 + 1), want: 791476407914},
		{name: "Issuance of Year 5", e: primitives.Epoch(params.BeaconConfig().EpochsPerYear*4 + 1), want: 608828006088},
		{name: "Issuance of Year 6", e: primitives.Epoch(params.BeaconConfig().EpochsPerYear*5 + 1), want: 487062404870},
		{name: "Issuance of Year 7", e: primitives.Epoch(params.BeaconConfig().EpochsPerYear*6 + 1), want: 365296803652},
		{name: "Issuance of Year 8", e: primitives.Epoch(params.BeaconConfig().EpochsPerYear*7 + 1), want: 304414003044},
		{name: "Issuance of Year 9", e: primitives.Epoch(params.BeaconConfig().EpochsPerYear*8 + 1), want: 280060882800},
		{name: "Issuance of Year 10", e: primitives.Epoch(params.BeaconConfig().EpochsPerYear*9 + 1), want: 207001522070},
		{name: "Issuance of Year 11", e: primitives.Epoch(params.BeaconConfig().EpochsPerYear*10 + 1), want: 182648401826},
		{name: "Issuance of Year 12", e: primitives.Epoch(params.BeaconConfig().EpochsPerYear*11 + 1), want: 182648401826},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := EpochIssuance(tt.e)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestTargetDepositPlan(t *testing.T) {
	tests := []struct {
		name string
		e    primitives.Epoch
		want uint64
	}{
		{name: "TargetDepositPlan of Epoch 0 (year 1)", e: primitives.Epoch(0), want: 5120000000000000},
		{name: "TargetDepositPlan of Epoch 41063 (year 1)", e: primitives.Epoch((params.BeaconConfig().EpochsPerYear + 1) / 2), want: 48640529923876874},
		{name: "TargetDepositPlan of Epoch 82125 (Year 2)", e: primitives.Epoch(params.BeaconConfig().EpochsPerYear), want: 92159999999960750},
		{name: "TargetDepositPlan of Epoch 123188 (Year 2)", e: primitives.Epoch((params.BeaconConfig().EpochsPerYear*3 + 1) / 2), want: 135680529923837624},
		{name: "TargetDepositPlan of Epoch 164250 (Year 3)", e: primitives.Epoch(params.BeaconConfig().EpochsPerYear * 2), want: 179199999999852250},
		{name: "TargetDepositPlan of Epoch 205313 (Year 3)", e: primitives.Epoch((params.BeaconConfig().EpochsPerYear*5 + 1) / 2), want: 188800116894792481},
		{name: "TargetDepositPlan of Epoch 246375 (Year 4)", e: primitives.Epoch(params.BeaconConfig().EpochsPerYear * 3), want: 198399999999778375},
		{name: "TargetDepositPlan of Epoch 287438 (Year 4)", e: primitives.Epoch((params.BeaconConfig().EpochsPerYear*7 + 1) / 2), want: 208000116894718606},
		{name: "TargetDepositPlan of Epoch 328500 (Year 5)", e: primitives.Epoch(params.BeaconConfig().EpochsPerYear * 4), want: 217599999999704500},
		{name: "TargetDepositPlan of Epoch 369563 (Year 5)", e: primitives.Epoch((params.BeaconConfig().EpochsPerYear*9 + 1) / 2), want: 227200116894644731},
		{name: "TargetDepositPlan of Epoch 410625 (Year 6)", e: primitives.Epoch(params.BeaconConfig().EpochsPerYear * 5), want: 236799999999630625},
		{name: "TargetDepositPlan of Epoch 451688 (Year 6)", e: primitives.Epoch((params.BeaconConfig().EpochsPerYear*11 + 1) / 2), want: 246400116894570856},
		{name: "TargetDepositPlan of Epoch 492750 (Year 7)", e: primitives.Epoch(params.BeaconConfig().EpochsPerYear * 6), want: 256000000 * 1e9},
		{name: "TargetDepositPlan of Epoch 533813 (Year 7)", e: primitives.Epoch((params.BeaconConfig().EpochsPerYear*13 + 1) / 2), want: 256000000 * 1e9},
		{name: "TargetDepositPlan of Epoch 574875 (Year 8)", e: primitives.Epoch(params.BeaconConfig().EpochsPerYear * 7), want: 256000000 * 1e9},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := TargetDepositPlan(tt.e)
			require.Equal(t, tt.want, got)
		})
	}
}

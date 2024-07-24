package blocks_test

import (
	"context"
	"testing"

	"github.com/prysmaticlabs/prysm/v5/beacon-chain/core/blocks"
	state_native "github.com/prysmaticlabs/prysm/v5/beacon-chain/state/state-native"
	"github.com/prysmaticlabs/prysm/v5/config/params"
	types "github.com/prysmaticlabs/prysm/v5/consensus-types/primitives"
	ethpb "github.com/prysmaticlabs/prysm/v5/proto/prysm/v1alpha1"
	"github.com/prysmaticlabs/prysm/v5/testing/require"
)

func TestProcessBailOut(t *testing.T) {
	ffe := params.BeaconConfig().FarFutureEpoch
	meb := params.BeaconConfig().MaxEffectiveBalance
	s, err := state_native.InitializeFromProtoAltair(&ethpb.BeaconStateAltair{
		Slot: types.Slot(uint64(params.BeaconConfig().ShardCommitteePeriod) * uint64(params.BeaconConfig().SlotsPerEpoch)),
		Validators: []*ethpb.Validator{
			{WithdrawableEpoch: ffe, ExitEpoch: ffe, EffectiveBalance: meb},
			{WithdrawableEpoch: ffe, ExitEpoch: ffe, EffectiveBalance: meb},
		},
		InactivityScores: []uint64{0, 1},
		BailOutScores:    []uint64{0, params.BeaconConfig().BailOutScoreThreshold + 1},
	})
	require.NoError(t, err)
	out0 := []*ethpb.BailOut{
		{ValidatorIndex: 0},
	}
	st, err := blocks.ProcessBailOuts(context.Background(), s, out0)
	require.ErrorContains(t, "validator bailout score is below bail out threshold", err)
	out1 := []*ethpb.BailOut{
		{ValidatorIndex: 1},
	}
	st, err = blocks.ProcessBailOuts(context.Background(), s, out1)
	require.NoError(t, err)
	require.DeepEqual(t, &ethpb.Validator{
		EffectiveBalance:  meb,
		WithdrawableEpoch: ffe,
		ExitEpoch:         ffe,
	}, st.Validators()[0])
	require.DeepEqual(t, &ethpb.Validator{
		EffectiveBalance:  meb,
		WithdrawableEpoch: params.MainnetConfig().ShardCommitteePeriod + params.MainnetConfig().MinValidatorWithdrawabilityDelay + 5,
		ExitEpoch:         params.MainnetConfig().ShardCommitteePeriod + 5,
	}, st.Validators()[1])
}

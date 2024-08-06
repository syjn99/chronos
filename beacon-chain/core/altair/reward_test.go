package altair_test

import (
	"math"
	"testing"

	"github.com/prysmaticlabs/prysm/v4/beacon-chain/core/altair"
	"github.com/prysmaticlabs/prysm/v4/beacon-chain/core/helpers"
	"github.com/prysmaticlabs/prysm/v4/beacon-chain/state"
	"github.com/prysmaticlabs/prysm/v4/config/params"
	"github.com/prysmaticlabs/prysm/v4/consensus-types/primitives"
	"github.com/prysmaticlabs/prysm/v4/testing/require"
	"github.com/prysmaticlabs/prysm/v4/testing/util"
)

func Test_BaseReward(t *testing.T) {
	helpers.ClearCache()
	genState := func(valCount uint64) state.ReadOnlyBeaconState {
		s, _ := util.DeterministicGenesisStateAltair(t, valCount)
		return s
	}
	tests := []struct {
		name      string
		valIdx    primitives.ValidatorIndex
		st        state.ReadOnlyBeaconState
		want      uint64
		errString string
	}{
		{
			name:      "unknown validator",
			valIdx:    2,
			st:        genState(1),
			want:      0,
			errString: "index 2 out of range",
		},
		{
			name:      "active balance is 256eth",
			valIdx:    0,
			st:        genState(1),
			want:      243531202432,
			errString: "",
		},
		{
			name:      "active balance is 256eth * target committee size",
			valIdx:    0,
			st:        genState(params.BeaconConfig().TargetCommitteeSize),
			want:      1902587488,
			errString: "",
		},
		{
			name:      "active balance is 256eth * max validator per committee size",
			valIdx:    0,
			st:        genState(params.BeaconConfig().MaxValidatorsPerCommittee),
			want:      118911712,
			errString: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := altair.BaseReward(tt.st, tt.valIdx)
			if (err != nil) && (tt.errString != "") {
				require.ErrorContains(t, tt.errString, err)
				return
			}
			require.Equal(t, tt.want, got)
		})
	}
}

func Test_BaseRewardWithTotalBalance(t *testing.T) {
	helpers.ClearCache()
	s, _ := util.DeterministicGenesisStateAltair(t, 1)
	tests := []struct {
		name          string
		valIdx        primitives.ValidatorIndex
		activeBalance uint64
		want          uint64
		wantReserve   uint64
		errString     string
	}{
		{
			name:          "active balance is 0",
			valIdx:        0,
			activeBalance: 0,
			want:          0,
			wantReserve:   0,
			errString:     "active balance can't be 0",
		},
		{
			name:          "unknown validator",
			valIdx:        2,
			activeBalance: 1,
			want:          0,
			wantReserve:   0,
			errString:     "index 2 out of range",
		},
		{
			name:          "active balance is 1",
			valIdx:        0,
			activeBalance: 1,
			want:          0,
			wantReserve:   0,
			errString:     "active balance can't be lower than effective balance increment",
		},
		{
			name:          "active balance is 1eth",
			valIdx:        0,
			activeBalance: params.BeaconConfig().EffectiveBalanceIncrement,
			want:          7792998477920,
			wantReserve:   0,
			errString:     "",
		},
		{
			name:          "active balance is 256eth",
			valIdx:        0,
			activeBalance: params.BeaconConfig().MaxEffectiveBalance,
			want:          243531202432,
			wantReserve:   0,
			errString:     "",
		},
		{
			name:          "active balance is 256eth * 1m validators",
			valIdx:        0,
			activeBalance: params.BeaconConfig().MaxEffectiveBalance * 1e6,
			want:          243520,
			wantReserve:   0,
			errString:     "",
		},
		{
			name:          "active balance is max uint64",
			valIdx:        0,
			activeBalance: math.MaxUint64,
			want:          3360,
			wantReserve:   0,
			errString:     "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, reserve, err := altair.BaseRewardWithTotalBalance(s, tt.valIdx, tt.activeBalance)
			if (err != nil) && (tt.errString != "") {
				require.ErrorContains(t, tt.errString, err)
				return
			}
			require.Equal(t, tt.want, got)
			require.Equal(t, tt.wantReserve, reserve)
		})
	}
}

func Test_BaseRewardPerIncrement(t *testing.T) {
	helpers.ClearCache()
	genState := func(valCount, rewardFactor uint64) state.ReadOnlyBeaconState {
		s, _ := util.DeterministicGenesisStateAltair(t, valCount)
		err := s.SetRewardAdjustmentFactor(rewardFactor)
		require.NoError(t, err)
		err = s.SetPreviousEpochReserve(10000000000)
		require.NoError(t, err)
		return s
	}
	tests := []struct {
		name          string
		epoch         primitives.Epoch
		activeBalance uint64
		want          uint64
		rewardFactor  uint64
		wantReserve   uint64
		errString     string
	}{
		{
			name:          "active balance is 0",
			epoch:         primitives.Epoch(0),
			activeBalance: 0,
			want:          0,
			rewardFactor:  0,
			wantReserve:   0,
			errString:     "active balance can't be 0",
		},
		{
			name:          "active balance is 1",
			epoch:         primitives.Epoch(0),
			activeBalance: 1,
			want:          0,
			rewardFactor:  0,
			wantReserve:   0,
			errString:     "active balance can't be lower than effective balance increment",
		},
		{
			name:          "active balance is 8 over",
			epoch:         primitives.Epoch(0),
			activeBalance: params.BeaconConfig().EffectiveBalanceIncrement,
			want:          243531202435,
			rewardFactor:  0,
			wantReserve:   0,
			errString:     "",
		},
		{
			name:          "active balance is 8 over and reward factor is 10000000000",
			epoch:         primitives.Epoch(0),
			activeBalance: params.BeaconConfig().EffectiveBalanceIncrement,
			want:          253531202435,
			rewardFactor:  10000000000,
			wantReserve:   10000000001,
			errString:     "",
		},
		{
			name:          "active balance is 256over",
			epoch:         primitives.Epoch(0),
			activeBalance: params.BeaconConfig().MaxEffectiveBalance,
			want:          7610350076,
			rewardFactor:  0,
			wantReserve:   0,
			errString:     "",
		},
		{
			name:          "active balance is 256over and reward factor is 10000000000",
			epoch:         primitives.Epoch(0),
			activeBalance: params.BeaconConfig().MaxEffectiveBalance,
			want:          7922850076,
			rewardFactor:  10000000000,
			wantReserve:   312500001,
			errString:     "",
		},
		{
			name:          "active balance is 256eth * 1m validators",
			epoch:         primitives.Epoch(0),
			activeBalance: params.BeaconConfig().MaxEffectiveBalance * 1e9,
			want:          120,
			rewardFactor:  0,
			wantReserve:   0,
			errString:     "",
		},
		{
			name:          "active balance is 256eth * 1m validators",
			epoch:         primitives.Epoch(0),
			activeBalance: params.BeaconConfig().MaxEffectiveBalance * 1e9,
			want:          125,
			rewardFactor:  10000000000,
			wantReserve:   5,
			errString:     "",
		},
		{
			name:          "active balance is max uint64",
			epoch:         primitives.Epoch(0),
			activeBalance: math.MaxUint64,
			want:          105,
			rewardFactor:  0,
			wantReserve:   0,
			errString:     "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, reserve, err := altair.BaseRewardPerIncrement(genState(20000, tt.rewardFactor), tt.activeBalance)
			if (err != nil) && (tt.errString != "") {
				require.ErrorContains(t, tt.errString, err)
				return
			}
			require.Equal(t, tt.want, got)
			require.Equal(t, tt.wantReserve, reserve)
		})
	}
}

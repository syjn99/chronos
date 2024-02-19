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
			want:      155859969536,
			errString: "",
		},
		{
			name:      "active balance is 256eth * target committee size",
			valIdx:    0,
			st:        genState(params.BeaconConfig().TargetCommitteeSize),
			want:      1217655808,
			errString: "",
		},
		{
			name:      "active balance is 256eth * max validator per committee size",
			valIdx:    0,
			st:        genState(params.BeaconConfig().MaxValidatorsPerCommittee),
			want:      76103424,
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
			want:          39900152206848,
			wantReserve:   0,
			errString:     "",
		},
		{
			name:          "active balance is 256eth",
			valIdx:        0,
			activeBalance: params.BeaconConfig().MaxEffectiveBalance,
			want:          155859969536,
			wantReserve:   0,
			errString:     "",
		},
		{
			name:          "active balance is 256eth * 1m validators",
			valIdx:        0,
			activeBalance: params.BeaconConfig().MaxEffectiveBalance * 1e6,
			want:          155648,
			wantReserve:   0,
			errString:     "",
		},
		{
			name:          "active balance is max uint64",
			valIdx:        0,
			activeBalance: math.MaxUint64,
			want:          2048,
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
	genState := func(valCount uint64) state.ReadOnlyBeaconState {
		s, _ := util.DeterministicGenesisStateAltair(t, valCount)
		return s
	}
	tests := []struct {
		name          string
		epoch         primitives.Epoch
		activeBalance uint64
		want          uint64
		wantReserve   uint64
		errString     string
	}{
		{
			name:          "active balance is 0",
			epoch:         primitives.Epoch(0),
			activeBalance: 0,
			want:          0,
			wantReserve:   0,
			errString:     "active balance can't be 0",
		},
		{
			name:          "active balance is 1",
			epoch:         primitives.Epoch(0),
			activeBalance: 1,
			want:          0,
			wantReserve:   0,
			errString:     "active balance can't be lower than effective balance increment",
		},
		{
			name:          "active balance is 1eth",
			epoch:         primitives.Epoch(0),
			activeBalance: params.BeaconConfig().EffectiveBalanceIncrement,
			want:          155859969558,
			wantReserve:   0,
			errString:     "",
		},
		{
			name:          "active balance is 256eth",
			epoch:         primitives.Epoch(0),
			activeBalance: params.BeaconConfig().MaxEffectiveBalance,
			want:          608828006,
			wantReserve:   0,
			errString:     "",
		},
		{
			name:          "active balance is 256eth * 1m validators",
			epoch:         primitives.Epoch(0),
			activeBalance: params.BeaconConfig().MaxEffectiveBalance * 1e9,
			want:          9,
			wantReserve:   0,
			errString:     "",
		},
		{
			name:          "active balance is max uint64",
			epoch:         primitives.Epoch(0),
			activeBalance: math.MaxUint64,
			want:          8,
			wantReserve:   0,
			errString:     "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, reserve, err := altair.BaseRewardPerIncrement(genState(20000), tt.activeBalance)
			if (err != nil) && (tt.errString != "") {
				require.ErrorContains(t, tt.errString, err)
				return
			}
			require.Equal(t, tt.want, got)
			require.Equal(t, tt.wantReserve, reserve)
		})
	}
}

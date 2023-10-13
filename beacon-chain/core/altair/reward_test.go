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
			name:      "active balance is 32eth",
			valIdx:    0,
			st:        genState(1),
			want:      19025875168, // 761035007584,
			errString: "",
		},
		{
			name:      "active balance is 32eth * target committee size",
			valIdx:    0,
			st:        genState(params.BeaconConfig().TargetCommitteeSize),
			want:      148639648, // 5945585984,
			errString: "",
		},
		{
			name:      "active balance is 32eth * max validator per committee size",
			valIdx:    0,
			st:        genState(params.BeaconConfig().MaxValidatorsPerCommittee),
			want:      9289952, // 371599104,
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
		errString     string
	}{
		{
			name:          "active balance is 0",
			valIdx:        0,
			activeBalance: 0,
			want:          0,
			errString:     "active balance can't be 0",
		},
		{
			name:          "unknown validator",
			valIdx:        2,
			activeBalance: 1,
			want:          0,
			errString:     "index 2 out of range",
		},
		{
			name:          "active balance is 1",
			valIdx:        0,
			activeBalance: 1,
			want:          0,
			errString:     "active balance can't be lower than effective balance increment",
		},
		{
			name:          "active balance is 1eth",
			valIdx:        0,
			activeBalance: params.BeaconConfig().EffectiveBalanceIncrement,
			want:          608828006080, // 24353120243520,
			errString:     "",
		},
		{
			name:          "active balance is 32eth",
			valIdx:        0,
			activeBalance: params.BeaconConfig().MaxEffectiveBalance,
			want:          19025875168, // 761035007584,
			errString:     "",
		},
		{
			name:          "active balance is 32eth * 1m validators",
			valIdx:        0,
			activeBalance: params.BeaconConfig().MaxEffectiveBalance * 1e9,
			want:          32, // 1792,
			errString:     "",
		},
		{
			name:          "active balance is max uint64",
			valIdx:        0,
			activeBalance: math.MaxUint64,
			want:          32, // 1312,
			errString:     "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := altair.BaseRewardWithTotalBalance(s, tt.valIdx, tt.activeBalance)
			if (err != nil) && (tt.errString != "") {
				require.ErrorContains(t, tt.errString, err)
				return
			}
			require.Equal(t, tt.want, got)
		})
	}
}

func Test_BaseRewardPerIncrement(t *testing.T) {
	helpers.ClearCache()
	tests := []struct {
		name          string
		epoch         primitives.Epoch
		activeBalance uint64
		want          uint64
		errString     string
	}{
		{
			name:          "active balance is 0",
			epoch:         primitives.Epoch(0),
			activeBalance: 0,
			want:          0,
			errString:     "active balance can't be 0",
		},
		{
			name:          "active balance is 1",
			epoch:         primitives.Epoch(0),
			activeBalance: 1,
			want:          0,
			errString:     "active balance can't be lower than effective balance increment",
		},
		{
			name:          "active balance is 1eth",
			epoch:         primitives.Epoch(0),
			activeBalance: params.BeaconConfig().EffectiveBalanceIncrement,
			want:          19025875190, // 761035007610,
			errString:     "",
		},
		{
			name:          "active balance is 32eth",
			epoch:         primitives.Epoch(0),
			activeBalance: params.BeaconConfig().MaxEffectiveBalance,
			want:          594558599, // 23782343987,
			errString:     "",
		},
		{
			name:          "active balance is 32eth * 1m validators",
			epoch:         primitives.Epoch(0),
			activeBalance: params.BeaconConfig().MaxEffectiveBalance * 1e9,
			want:          1, // 56,
			errString:     "",
		},
		{
			name:          "active balance is max uint64",
			epoch:         primitives.Epoch(0),
			activeBalance: math.MaxUint64,
			want:          1, // 41,
			errString:     "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := altair.BaseRewardPerIncrement(tt.epoch, tt.activeBalance)
			if (err != nil) && (tt.errString != "") {
				require.ErrorContains(t, tt.errString, err)
				return
			}
			require.Equal(t, tt.want, got)
		})
	}
}

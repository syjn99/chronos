package altair_test

import (
	"context"
	"testing"

	"github.com/prysmaticlabs/go-bitfield"
	"github.com/prysmaticlabs/prysm/v4/beacon-chain/core/altair"
	"github.com/prysmaticlabs/prysm/v4/beacon-chain/core/helpers"
	"github.com/prysmaticlabs/prysm/v4/beacon-chain/core/signing"
	"github.com/prysmaticlabs/prysm/v4/beacon-chain/core/time"
	p2pType "github.com/prysmaticlabs/prysm/v4/beacon-chain/p2p/types"
	fieldparams "github.com/prysmaticlabs/prysm/v4/config/fieldparams"
	"github.com/prysmaticlabs/prysm/v4/config/params"
	"github.com/prysmaticlabs/prysm/v4/consensus-types/primitives"
	"github.com/prysmaticlabs/prysm/v4/crypto/bls"
	"github.com/prysmaticlabs/prysm/v4/encoding/bytesutil"
	ethpb "github.com/prysmaticlabs/prysm/v4/proto/prysm/v1alpha1"
	"github.com/prysmaticlabs/prysm/v4/testing/assert"
	"github.com/prysmaticlabs/prysm/v4/testing/require"
	"github.com/prysmaticlabs/prysm/v4/testing/util"
	"github.com/prysmaticlabs/prysm/v4/time/slots"
)

func TestProcessSyncCommittee_PerfectParticipation(t *testing.T) {
	beaconState, privKeys := util.DeterministicGenesisStateAltair(t, params.BeaconConfig().MaxValidatorsPerCommittee)
	require.NoError(t, beaconState.SetSlot(1))
	committee, err := altair.NextSyncCommittee(context.Background(), beaconState)
	require.NoError(t, err)
	require.NoError(t, beaconState.SetCurrentSyncCommittee(committee))

	syncBits := bitfield.NewBitvector512()
	for i := range syncBits {
		syncBits[i] = 0xff
	}
	indices, err := altair.NextSyncCommitteeIndices(context.Background(), beaconState)
	require.NoError(t, err)
	ps := slots.PrevSlot(beaconState.Slot())
	pbr, err := helpers.BlockRootAtSlot(beaconState, ps)
	require.NoError(t, err)
	sigs := make([]bls.Signature, len(indices))
	for i, indice := range indices {
		b := p2pType.SSZBytes(pbr)
		sb, err := signing.ComputeDomainAndSign(beaconState, time.CurrentEpoch(beaconState), &b, params.BeaconConfig().DomainSyncCommittee, privKeys[indice])
		require.NoError(t, err)
		sig, err := bls.SignatureFromBytes(sb)
		require.NoError(t, err)
		sigs[i] = sig
	}
	aggregatedSig := bls.AggregateSignatures(sigs).Marshal()
	syncAggregate := &ethpb.SyncAggregate{
		SyncCommitteeBits:      syncBits,
		SyncCommitteeSignature: aggregatedSig,
	}

	var reward uint64
	beaconState, reward, err = altair.ProcessSyncAggregate(context.Background(), beaconState, syncAggregate)
	require.NoError(t, err)
	//assert.Equal(t, uint64(2654208), reward) // 106170880 // TODO(john): Fix this test after sync committee rewards are enabled.
	assert.Equal(t, uint64(0), reward) // 106170880

	// Use a non-sync committee index to compare profitability.
	syncCommittee := make(map[primitives.ValidatorIndex]bool)
	for _, index := range indices {
		syncCommittee[index] = true
	}
	nonSyncIndex := primitives.ValidatorIndex(params.BeaconConfig().MaxValidatorsPerCommittee + 1)
	for i := primitives.ValidatorIndex(0); uint64(i) < params.BeaconConfig().MaxValidatorsPerCommittee; i++ {
		if !syncCommittee[i] {
			nonSyncIndex = i
			break
		}
	}

	// Sync committee should be more profitable than non sync committee
	balances := beaconState.Balances()
	//require.Equal(t, true, balances[indices[0]] > balances[nonSyncIndex]) // TODO(john): Fix this test after sync committee rewards are enabled.
	require.Equal(t, true, balances[indices[0]] == balances[nonSyncIndex])

	// Proposer should be more profitable than rest of the sync committee
	proposerIndex, err := helpers.BeaconProposerIndex(context.Background(), beaconState)
	require.NoError(t, err)
	//require.Equal(t, true, balances[proposerIndex] > balances[indices[0]]) // TODO(john): Fix this test after sync committee rewards are enabled.
	require.Equal(t, true, balances[proposerIndex] == balances[indices[0]])

	// Sync committee should have the same profits, except you are a proposer
	for i := 1; i < len(indices); i++ {
		if proposerIndex == indices[i-1] || proposerIndex == indices[i] {
			continue
		}
		require.Equal(t, balances[indices[i-1]], balances[indices[i]])
	}

	// Increased balance validator count should equal to sync committee count
	// TODO(john): Fix this test after sync committee rewards are enabled.
	//increased := uint64(0)
	//for _, balance := range balances {
	//	if balance > params.BeaconConfig().MaxEffectiveBalance {
	//		increased++
	//	}
	//}
	//require.Equal(t, params.BeaconConfig().SyncCommitteeSize, increased)
}

func TestProcessSyncCommittee_MixParticipation_BadSignature(t *testing.T) {
	beaconState, privKeys := util.DeterministicGenesisStateAltair(t, params.BeaconConfig().MaxValidatorsPerCommittee)
	require.NoError(t, beaconState.SetSlot(1))
	committee, err := altair.NextSyncCommittee(context.Background(), beaconState)
	require.NoError(t, err)
	require.NoError(t, beaconState.SetCurrentSyncCommittee(committee))

	syncBits := bitfield.NewBitvector512()
	for i := range syncBits {
		syncBits[i] = 0xAA
	}
	indices, err := altair.NextSyncCommitteeIndices(context.Background(), beaconState)
	require.NoError(t, err)
	ps := slots.PrevSlot(beaconState.Slot())
	pbr, err := helpers.BlockRootAtSlot(beaconState, ps)
	require.NoError(t, err)
	sigs := make([]bls.Signature, len(indices))
	for i, indice := range indices {
		b := p2pType.SSZBytes(pbr)
		sb, err := signing.ComputeDomainAndSign(beaconState, time.CurrentEpoch(beaconState), &b, params.BeaconConfig().DomainSyncCommittee, privKeys[indice])
		require.NoError(t, err)
		sig, err := bls.SignatureFromBytes(sb)
		require.NoError(t, err)
		sigs[i] = sig
	}
	aggregatedSig := bls.AggregateSignatures(sigs).Marshal()
	syncAggregate := &ethpb.SyncAggregate{
		SyncCommitteeBits:      syncBits,
		SyncCommitteeSignature: aggregatedSig,
	}

	_, _, err = altair.ProcessSyncAggregate(context.Background(), beaconState, syncAggregate)
	require.ErrorContains(t, "invalid sync committee signature", err)
}

func TestProcessSyncCommittee_MixParticipation_GoodSignature(t *testing.T) {
	beaconState, privKeys := util.DeterministicGenesisStateAltair(t, params.BeaconConfig().MaxValidatorsPerCommittee)
	require.NoError(t, beaconState.SetSlot(1))
	committee, err := altair.NextSyncCommittee(context.Background(), beaconState)
	require.NoError(t, err)
	require.NoError(t, beaconState.SetCurrentSyncCommittee(committee))

	syncBits := bitfield.NewBitvector512()
	for i := range syncBits {
		syncBits[i] = 0xAA
	}
	indices, err := altair.NextSyncCommitteeIndices(context.Background(), beaconState)
	require.NoError(t, err)
	ps := slots.PrevSlot(beaconState.Slot())
	pbr, err := helpers.BlockRootAtSlot(beaconState, ps)
	require.NoError(t, err)
	sigs := make([]bls.Signature, 0, len(indices))
	for i, indice := range indices {
		if syncBits.BitAt(uint64(i)) {
			b := p2pType.SSZBytes(pbr)
			sb, err := signing.ComputeDomainAndSign(beaconState, time.CurrentEpoch(beaconState), &b, params.BeaconConfig().DomainSyncCommittee, privKeys[indice])
			require.NoError(t, err)
			sig, err := bls.SignatureFromBytes(sb)
			require.NoError(t, err)
			sigs = append(sigs, sig)
		}
	}
	aggregatedSig := bls.AggregateSignatures(sigs).Marshal()
	syncAggregate := &ethpb.SyncAggregate{
		SyncCommitteeBits:      syncBits,
		SyncCommitteeSignature: aggregatedSig,
	}

	_, _, err = altair.ProcessSyncAggregate(context.Background(), beaconState, syncAggregate)
	require.NoError(t, err)
}

// This is a regression test #11696
func TestProcessSyncCommittee_DontPrecompute(t *testing.T) {
	beaconState, _ := util.DeterministicGenesisStateAltair(t, params.BeaconConfig().MaxValidatorsPerCommittee)
	require.NoError(t, beaconState.SetSlot(1))
	committee, err := altair.NextSyncCommittee(context.Background(), beaconState)
	require.NoError(t, err)
	committeeKeys := committee.Pubkeys
	committeeKeys[1] = committeeKeys[0]
	require.NoError(t, beaconState.SetCurrentSyncCommittee(committee))
	idx, ok := beaconState.ValidatorIndexByPubkey(bytesutil.ToBytes48(committeeKeys[0]))
	require.Equal(t, true, ok)

	syncBits := bitfield.NewBitvector512()
	for i := range syncBits {
		syncBits[i] = 0xFF
	}
	syncBits.SetBitAt(0, false)
	syncAggregate := &ethpb.SyncAggregate{
		SyncCommitteeBits: syncBits,
	}
	require.NoError(t, beaconState.UpdateBalancesAtIndex(idx, 0))
	st, votedKeys, _, err := altair.ProcessSyncAggregateEported(context.Background(), beaconState, syncAggregate)
	require.NoError(t, err)
	require.Equal(t, 511, len(votedKeys))
	require.DeepEqual(t, committeeKeys[0], votedKeys[0].Marshal())
	balances := st.Balances()
	//require.Equal(t, uint64(36288), balances[idx]) // 1451559
	require.Equal(t, uint64(0), balances[idx]) // TODO(john): Fix this test after sync committee rewards are enabled.
}

func TestProcessSyncCommittee_processSyncAggregate(t *testing.T) {
	beaconState, _ := util.DeterministicGenesisStateAltair(t, params.BeaconConfig().MaxValidatorsPerCommittee)
	require.NoError(t, beaconState.SetSlot(1))
	committee, err := altair.NextSyncCommittee(context.Background(), beaconState)
	require.NoError(t, err)
	require.NoError(t, beaconState.SetCurrentSyncCommittee(committee))

	syncBits := bitfield.NewBitvector512()
	for i := range syncBits {
		syncBits[i] = 0xAA
	}
	syncAggregate := &ethpb.SyncAggregate{
		SyncCommitteeBits: syncBits,
	}

	st, votedKeys, _, err := altair.ProcessSyncAggregateEported(context.Background(), beaconState, syncAggregate)
	require.NoError(t, err)
	votedMap := make(map[[fieldparams.BLSPubkeyLength]byte]bool)
	for _, key := range votedKeys {
		votedMap[bytesutil.ToBytes48(key.Marshal())] = true
	}
	require.Equal(t, int(syncBits.Len()/2), len(votedKeys))

	currentSyncCommittee, err := st.CurrentSyncCommittee()
	require.NoError(t, err)
	committeeKeys := currentSyncCommittee.Pubkeys
	balances := st.Balances()

	proposerIndex, err := helpers.BeaconProposerIndex(context.Background(), beaconState)
	require.NoError(t, err)

	for i := 0; i < len(syncBits); i++ {
		if syncBits.BitAt(uint64(i)) {
			pk := bytesutil.ToBytes48(committeeKeys[i])
			require.DeepEqual(t, true, votedMap[pk])
			idx, ok := st.ValidatorIndexByPubkey(pk)
			require.Equal(t, true, ok)
			require.Equal(t, uint64(32000000000), balances[idx]) // 32001451559
		} else {
			pk := bytesutil.ToBytes48(committeeKeys[i])
			require.DeepEqual(t, false, votedMap[pk])
			idx, ok := st.ValidatorIndexByPubkey(pk)
			require.Equal(t, true, ok)
			if idx != proposerIndex {
				// TODO(john): Fix this test after sync committee rewards are enabled.
				// require.Equal(t, uint64(31999963712), balances[idx]) // 31998548441
				require.Equal(t, uint64(32000000000), balances[idx]) // 31998548441
			}
		}
	}
	require.Equal(t, uint64(32000000000), balances[proposerIndex]) // 32051633881
}

func Test_VerifySyncCommitteeSig(t *testing.T) {
	beaconState, privKeys := util.DeterministicGenesisStateAltair(t, params.BeaconConfig().MaxValidatorsPerCommittee)
	require.NoError(t, beaconState.SetSlot(1))
	committee, err := altair.NextSyncCommittee(context.Background(), beaconState)
	require.NoError(t, err)
	require.NoError(t, beaconState.SetCurrentSyncCommittee(committee))

	syncBits := bitfield.NewBitvector512()
	for i := range syncBits {
		syncBits[i] = 0xff
	}
	indices, err := altair.NextSyncCommitteeIndices(context.Background(), beaconState)
	require.NoError(t, err)
	ps := slots.PrevSlot(beaconState.Slot())
	pbr, err := helpers.BlockRootAtSlot(beaconState, ps)
	require.NoError(t, err)
	sigs := make([]bls.Signature, len(indices))
	pks := make([]bls.PublicKey, len(indices))
	for i, indice := range indices {
		b := p2pType.SSZBytes(pbr)
		sb, err := signing.ComputeDomainAndSign(beaconState, time.CurrentEpoch(beaconState), &b, params.BeaconConfig().DomainSyncCommittee, privKeys[indice])
		require.NoError(t, err)
		sig, err := bls.SignatureFromBytes(sb)
		require.NoError(t, err)
		sigs[i] = sig
		pks[i] = privKeys[indice].PublicKey()
	}
	aggregatedSig := bls.AggregateSignatures(sigs).Marshal()

	blsKey, err := bls.RandKey()
	require.NoError(t, err)
	require.ErrorContains(t, "invalid sync committee signature", altair.VerifySyncCommitteeSig(beaconState, pks, blsKey.Sign([]byte{'m', 'e', 'o', 'w'}).Marshal()))

	require.NoError(t, altair.VerifySyncCommitteeSig(beaconState, pks, aggregatedSig))
}

func Test_SyncRewards(t *testing.T) {
	tests := []struct {
		name                  string
		epoch                 primitives.Epoch
		wantProposerReward    uint64
		wantParticipantReward uint64
		errString             string
	}{
		{
			name:                  "epoch is 1 (year 1) ",
			epoch:                 primitives.Epoch(0),
			wantProposerReward:    0, // 5184,  // 207365,
			wantParticipantReward: 0, // 36288, // 1451559,
			errString:             "active balance can't be 0",
		},
		{
			name:                  "epoch is 82126 (year 2)",
			epoch:                 primitives.Epoch(params.BeaconConfig().EpochsPerYear + 1),
			wantProposerReward:    0, // 4769,  // 190776,
			wantParticipantReward: 0, // 33385, // 1335434,
			errString:             "",
		},
		{
			name:                  "epoch is 164251 (year 3)",
			epoch:                 primitives.Epoch(params.BeaconConfig().EpochsPerYear*2 + 1),
			wantProposerReward:    0, // 4354,  // 174187,
			wantParticipantReward: 0, // 30482, // 1219309,
			errString:             "",
		},
		{
			name:                  "epoch is 82126 (year 4)",
			epoch:                 primitives.Epoch(params.BeaconConfig().EpochsPerYear*3 + 1),
			wantProposerReward:    0, // 3939,  // 157597,
			wantParticipantReward: 0, // 27579, // 1103184,
			errString:             "",
		},
		{
			name:                  "epoch is 82126 (year 5)",
			epoch:                 primitives.Epoch(params.BeaconConfig().EpochsPerYear*4 + 1),
			wantProposerReward:    0, // 3525,  // 141008,
			wantParticipantReward: 0, // 24676, // 987060,
			errString:             "",
		},
		{
			name:                  "epoch is 82126 (year 6)",
			epoch:                 primitives.Epoch(params.BeaconConfig().EpochsPerYear*5 + 1),
			wantProposerReward:    0, // 3110,  // 124419,
			wantParticipantReward: 0, // 21773, // 870935,
			errString:             "",
		},
		{
			name:                  "epoch is 82126 (year 7)",
			epoch:                 primitives.Epoch(params.BeaconConfig().EpochsPerYear*6 + 1),
			wantProposerReward:    0, // 2695,  // 107830,
			wantParticipantReward: 0, // 18870, // 754810,
			errString:             "",
		},
		{
			name:                  "epoch is 82126 (year 8)",
			epoch:                 primitives.Epoch(params.BeaconConfig().EpochsPerYear*7 + 1),
			wantProposerReward:    0, // 2281,  // 91240,
			wantParticipantReward: 0, // 15967, // 638685,
			errString:             "",
		},
		{
			name:                  "epoch is 82126 (year 9)",
			epoch:                 primitives.Epoch(params.BeaconConfig().EpochsPerYear*8 + 1),
			wantProposerReward:    0, // 1866,  // 74651,
			wantParticipantReward: 0, // 13064, // 522561,
			errString:             "",
		},
		{
			name:                  "epoch is 82126 (year 10)",
			epoch:                 primitives.Epoch(params.BeaconConfig().EpochsPerYear*9 + 1),
			wantProposerReward:    0, // 1451,  // 58062,
			wantParticipantReward: 0, // 10160, // 406436,
			errString:             "",
		},
		{
			name:                  "epoch is 82126 (year 11)",
			epoch:                 primitives.Epoch(params.BeaconConfig().EpochsPerYear*10 + 1),
			wantProposerReward:    0, // 1244, // 49767,
			wantParticipantReward: 0, // 8709, // 348374,
			errString:             "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			proposerReward, participantReward, err := altair.SyncRewards(tt.epoch)
			if (err != nil) && (tt.errString != "") {
				require.ErrorContains(t, tt.errString, err)
				return
			}
			require.Equal(t, tt.wantProposerReward, proposerReward)
			require.Equal(t, tt.wantParticipantReward, participantReward)
		})
	}
}

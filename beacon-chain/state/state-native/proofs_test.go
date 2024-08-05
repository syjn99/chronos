package state_native_test

import (
	"context"
	"testing"

	"github.com/ethereum/go-ethereum/common/hexutil"
	statenative "github.com/prysmaticlabs/prysm/v4/beacon-chain/state/state-native"
	"github.com/prysmaticlabs/prysm/v4/container/trie"
	"github.com/prysmaticlabs/prysm/v4/testing/require"
	"github.com/prysmaticlabs/prysm/v4/testing/util"
)

func TestBeaconStateMerkleProofs_phase0_notsupported(t *testing.T) {
	ctx := context.Background()
	st, _ := util.DeterministicGenesisState(t, 256)
	t.Run("current sync committee", func(t *testing.T) {
		_, err := st.CurrentSyncCommitteeProof(ctx)
		require.ErrorContains(t, "not supported", err)
	})
	t.Run("next sync committee", func(t *testing.T) {
		_, err := st.NextSyncCommitteeProof(ctx)
		require.ErrorContains(t, "not supported", err)
	})
	t.Run("finalized root", func(t *testing.T) {
		_, err := st.FinalizedRootProof(ctx)
		require.ErrorContains(t, "not supported", err)
	})
}

func TestBeaconStateMerkleProofs_altair(t *testing.T) {
	ctx := context.Background()
	altair, err := util.NewBeaconStateAltair()
	require.NoError(t, err)
	htr, err := altair.HashTreeRoot(ctx)
	require.NoError(t, err)
	t.Run("current sync committee", func(t *testing.T) {
		currentResults := []string{
			"0xacff3e632bf8ff27b783ac48086a544d1e920512add91817790d355e09846cd0",
			"0x966474d982a6695e57f9cd7dbcd1d646b817acc468c6a9c8fa9da4bd25428d9d",
			"0xdb56114e00fdd4c1f85c892bf35ac9a89289aaecb1ebd0a96cde606a748b5d71",
			"0xd96747cbf167ae33a66b7c32dd372eaba88f4b0d0669f319f06a1f1b193be155",
			"0x6e84be54558f6577ea6770acf4062accf291b086abdd724025eb92b99095e266",
		}

		cscp, err := altair.CurrentSyncCommitteeProof(ctx)
		require.NoError(t, err)
		require.Equal(t, 5, len(cscp))
		for i, bytes := range cscp {
			require.Equal(t, currentResults[i], hexutil.Encode(bytes))
		}
	})
	t.Run("next sync committee", func(t *testing.T) {
		nextResults := []string{
			"0xacff3e632bf8ff27b783ac48086a544d1e920512add91817790d355e09846cd0",
			"0x1b1897375f88a804a22d60a8855ef68068011015741b2bd91238e7c93fbc9ad9",
			"0xdb56114e00fdd4c1f85c892bf35ac9a89289aaecb1ebd0a96cde606a748b5d71",
			"0xd96747cbf167ae33a66b7c32dd372eaba88f4b0d0669f319f06a1f1b193be155",
			"0x6e84be54558f6577ea6770acf4062accf291b086abdd724025eb92b99095e266",
		}

		nscp, err := altair.NextSyncCommitteeProof(ctx)
		require.NoError(t, err)
		require.Equal(t, 5, len(nscp))
		for i, bytes := range nscp {
			require.Equal(t, nextResults[i], hexutil.Encode(bytes))
		}
	})
	t.Run("finalized root", func(t *testing.T) {
		finalizedRoot := altair.FinalizedCheckpoint().Root
		proof, err := altair.FinalizedRootProof(ctx)
		require.NoError(t, err)
		gIndex := statenative.FinalizedRootGeneralizedIndex()
		valid := trie.VerifyMerkleProof(htr[:], finalizedRoot, gIndex, proof)
		require.Equal(t, true, valid)
	})
	t.Run("recomputes root on dirty fields", func(t *testing.T) {
		currentRoot, err := altair.HashTreeRoot(ctx)
		require.NoError(t, err)
		cpt := altair.FinalizedCheckpoint()
		require.NoError(t, err)

		// Edit the checkpoint.
		cpt.Epoch = 100
		require.NoError(t, altair.SetFinalizedCheckpoint(cpt))

		// Produce a proof for the finalized root.
		proof, err := altair.FinalizedRootProof(ctx)
		require.NoError(t, err)

		// We expect the previous step to have triggered
		// a recomputation of dirty fields in the beacon state, resulting
		// in a new hash tree root as the finalized checkpoint had previously
		// changed and should have been marked as a dirty state field.
		// The proof validity should be false for the old root, but true for the new.
		finalizedRoot := altair.FinalizedCheckpoint().Root
		gIndex := statenative.FinalizedRootGeneralizedIndex()
		valid := trie.VerifyMerkleProof(currentRoot[:], finalizedRoot, gIndex, proof)
		require.Equal(t, false, valid)

		newRoot, err := altair.HashTreeRoot(ctx)
		require.NoError(t, err)

		valid = trie.VerifyMerkleProof(newRoot[:], finalizedRoot, gIndex, proof)
		require.Equal(t, true, valid)
	})
}

func TestBeaconStateMerkleProofs_bellatrix(t *testing.T) {
	ctx := context.Background()
	bellatrix, err := util.NewBeaconStateBellatrix()
	require.NoError(t, err)
	htr, err := bellatrix.HashTreeRoot(ctx)
	require.NoError(t, err)
	t.Run("current sync committee", func(t *testing.T) {
		currentResults := []string{
			"0xacff3e632bf8ff27b783ac48086a544d1e920512add91817790d355e09846cd0",
			"0x966474d982a6695e57f9cd7dbcd1d646b817acc468c6a9c8fa9da4bd25428d9d",
			"0x220339fd428c8d07195b94016bffe0bf9137fd9122dfde141fe4006e03d80c4a",
			"0xd96747cbf167ae33a66b7c32dd372eaba88f4b0d0669f319f06a1f1b193be155",
			"0x6e84be54558f6577ea6770acf4062accf291b086abdd724025eb92b99095e266",
		}

		cscp, err := bellatrix.CurrentSyncCommitteeProof(ctx)
		require.NoError(t, err)
		require.Equal(t, 5, len(cscp))
		for i, bytes := range cscp {
			require.Equal(t, currentResults[i], hexutil.Encode(bytes))
		}
	})
	t.Run("next sync committee", func(t *testing.T) {
		nextResults := []string{
			"0xacff3e632bf8ff27b783ac48086a544d1e920512add91817790d355e09846cd0",
			"0x1b1897375f88a804a22d60a8855ef68068011015741b2bd91238e7c93fbc9ad9",
			"0x220339fd428c8d07195b94016bffe0bf9137fd9122dfde141fe4006e03d80c4a",
			"0xd96747cbf167ae33a66b7c32dd372eaba88f4b0d0669f319f06a1f1b193be155",
			"0x6e84be54558f6577ea6770acf4062accf291b086abdd724025eb92b99095e266",
		}

		nscp, err := bellatrix.NextSyncCommitteeProof(ctx)
		require.NoError(t, err)
		require.Equal(t, 5, len(nscp))
		for i, bytes := range nscp {
			require.Equal(t, nextResults[i], hexutil.Encode(bytes))
		}
	})
	t.Run("finalized root", func(t *testing.T) {
		finalizedRoot := bellatrix.FinalizedCheckpoint().Root
		proof, err := bellatrix.FinalizedRootProof(ctx)
		require.NoError(t, err)
		gIndex := statenative.FinalizedRootGeneralizedIndex()
		valid := trie.VerifyMerkleProof(htr[:], finalizedRoot, gIndex, proof)
		require.Equal(t, true, valid)
	})
	t.Run("recomputes root on dirty fields", func(t *testing.T) {
		currentRoot, err := bellatrix.HashTreeRoot(ctx)
		require.NoError(t, err)
		cpt := bellatrix.FinalizedCheckpoint()
		require.NoError(t, err)

		// Edit the checkpoint.
		cpt.Epoch = 100
		require.NoError(t, bellatrix.SetFinalizedCheckpoint(cpt))

		// Produce a proof for the finalized root.
		proof, err := bellatrix.FinalizedRootProof(ctx)
		require.NoError(t, err)

		// We expect the previous step to have triggered
		// a recomputation of dirty fields in the beacon state, resulting
		// in a new hash tree root as the finalized checkpoint had previously
		// changed and should have been marked as a dirty state field.
		// The proof validity should be false for the old root, but true for the new.
		finalizedRoot := bellatrix.FinalizedCheckpoint().Root
		gIndex := statenative.FinalizedRootGeneralizedIndex()
		valid := trie.VerifyMerkleProof(currentRoot[:], finalizedRoot, gIndex, proof)
		require.Equal(t, false, valid)

		newRoot, err := bellatrix.HashTreeRoot(ctx)
		require.NoError(t, err)

		valid = trie.VerifyMerkleProof(newRoot[:], finalizedRoot, gIndex, proof)
		require.Equal(t, true, valid)
	})
}

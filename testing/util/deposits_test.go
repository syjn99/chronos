package util

import (
	"bytes"
	"context"
	"encoding/hex"
	"testing"

	"github.com/prysmaticlabs/prysm/v4/beacon-chain/core/helpers"
	"github.com/prysmaticlabs/prysm/v4/config/params"
	"github.com/prysmaticlabs/prysm/v4/testing/require"
	"google.golang.org/protobuf/proto"
)

func TestSetupInitialDeposits_1024Entries(t *testing.T) {
	entries := 1
	resetCache()
	deposits, privKeys, err := DeterministicDepositsAndKeys(uint64(entries))
	require.NoError(t, err)
	_, depositDataRoots, err := DeterministicDepositTrie(len(deposits))
	require.NoError(t, err)

	if len(deposits) != entries {
		t.Fatalf("incorrect number of deposits returned, wanted %d but received %d", entries, len(deposits))
	}
	if len(privKeys) != entries {
		t.Fatalf("incorrect number of private keys returned, wanted %d but received %d", entries, len(privKeys))
	}
	expectedPublicKeyAt0 := []byte{0xa9, 0x9a, 0x76, 0xed, 0x77, 0x96, 0xf7, 0xbe, 0x22, 0xd5, 0xb7, 0xe8, 0x5d, 0xee, 0xb7, 0xc5, 0x67, 0x7e, 0x88, 0xe5, 0x11, 0xe0, 0xb3, 0x37, 0x61, 0x8f, 0x8c, 0x4e, 0xb6, 0x13, 0x49, 0xb4, 0xbf, 0x2d, 0x15, 0x3f, 0x64, 0x9f, 0x7b, 0x53, 0x35, 0x9f, 0xe8, 0xb9, 0x4a, 0x38, 0xe4, 0x4c}
	if !bytes.Equal(deposits[0].Data.PublicKey, expectedPublicKeyAt0) {
		t.Fatalf("incorrect public key, wanted %x but received %x", expectedPublicKeyAt0, deposits[0].Data.PublicKey)
	}
	expectedWithdrawalCredentialsAt0 := []byte{0x00, 0xec, 0x7e, 0xf7, 0x78, 0x0c, 0x9d, 0x15, 0x15, 0x97, 0x92, 0x40, 0x36, 0x26, 0x2d, 0xd2, 0x8d, 0xc6, 0x0e, 0x12, 0x28, 0xf4, 0xda, 0x6f, 0xec, 0xf9, 0xd4, 0x02, 0xcb, 0x3f, 0x35, 0x94}
	if !bytes.Equal(deposits[0].Data.WithdrawalCredentials, expectedWithdrawalCredentialsAt0) {
		t.Fatalf("incorrect withdrawal credentials, wanted %x but received %x", expectedWithdrawalCredentialsAt0, deposits[0].Data.WithdrawalCredentials)
	}

	dRootAt0 := []byte("1a233634bd6443d9de59c28967a2cd901b0866c9a6c35f2f18ee1339a9bae4cc")
	dRootAt0B := make([]byte, hex.DecodedLen(len(dRootAt0)))
	_, err = hex.Decode(dRootAt0B, dRootAt0)
	require.NoError(t, err)
	if !bytes.Equal(depositDataRoots[0][:], dRootAt0B) {
		t.Fatalf("incorrect deposit data root, wanted %x but received %x", dRootAt0B, depositDataRoots[0])
	}

	sigAt0 := []byte("a35ee91bc5759bbef6e10d3916fe65e3ddf4f5ab735e641712b9cf60b73e17e71c26c2f5316aecb54453f09bb43032c503ab78f0ae23a97a70db69e2330b5345b5e607ee48534c755528cfdd75b8150ef7e79896e7c0d242d462216663df8a65")
	sigAt0B := make([]byte, hex.DecodedLen(len(sigAt0)))
	_, err = hex.Decode(sigAt0B, sigAt0)
	require.NoError(t, err)
	if !bytes.Equal(deposits[0].Data.Signature, sigAt0B) {
		t.Fatalf("incorrect signature, wanted %x but received %x", sigAt0B, deposits[0].Data.Signature)
	}

	entries = 1024
	resetCache()
	deposits, privKeys, err = DeterministicDepositsAndKeys(uint64(entries))
	require.NoError(t, err)
	_, depositDataRoots, err = DeterministicDepositTrie(len(deposits))
	require.NoError(t, err)
	if len(deposits) != entries {
		t.Fatalf("incorrect number of deposits returned, wanted %d but received %d", entries, len(deposits))
	}
	if len(privKeys) != entries {
		t.Fatalf("incorrect number of private keys returned, wanted %d but received %d", entries, len(privKeys))
	}
	// Ensure 0  has not changed
	if !bytes.Equal(deposits[0].Data.PublicKey, expectedPublicKeyAt0) {
		t.Fatalf("incorrect public key, wanted %x but received %x", expectedPublicKeyAt0, deposits[0].Data.PublicKey)
	}
	if !bytes.Equal(deposits[0].Data.WithdrawalCredentials, expectedWithdrawalCredentialsAt0) {
		t.Fatalf("incorrect withdrawal credentials, wanted %x but received %x", expectedWithdrawalCredentialsAt0, deposits[0].Data.WithdrawalCredentials)
	}
	if !bytes.Equal(depositDataRoots[0][:], dRootAt0B) {
		t.Fatalf("incorrect deposit data root, wanted %x but received %x", dRootAt0B, depositDataRoots[0])
	}
	if !bytes.Equal(deposits[0].Data.Signature, sigAt0B) {
		t.Fatalf("incorrect signature, wanted %x but received %x", sigAt0B, deposits[0].Data.Signature)
	}
	expectedPublicKeyAt1023 := []byte{0x81, 0x2b, 0x93, 0x5e, 0xc8, 0x4b, 0x0e, 0x9a, 0x83, 0x95, 0x55, 0xaf, 0x33, 0x60, 0xca, 0xfb, 0x83, 0x1b, 0xd6, 0x12, 0xcf, 0xa2, 0x2e, 0x25, 0xea, 0xb0, 0x3c, 0xf5, 0xfd, 0xb0, 0x2a, 0xf5, 0x2b, 0xa4, 0x01, 0x7a, 0xee, 0xa8, 0x8a, 0x2f, 0x62, 0x2c, 0x78, 0x6e, 0x7f, 0x47, 0x6f, 0x4b}
	if !bytes.Equal(deposits[1023].Data.PublicKey, expectedPublicKeyAt1023) {
		t.Fatalf("incorrect public key, wanted %x but received %x", expectedPublicKeyAt1023, deposits[1023].Data.PublicKey)
	}
	expectedWithdrawalCredentialsAt1023 := []byte{0x00, 0x23, 0xd5, 0x76, 0xbc, 0x6c, 0x15, 0xdb, 0xc4, 0x34, 0x70, 0x1f, 0x3f, 0x41, 0xfd, 0x3e, 0x67, 0x59, 0xd2, 0xea, 0x7c, 0xdc, 0x64, 0x71, 0x0e, 0xe2, 0x8d, 0xde, 0xf7, 0xd2, 0xda, 0x28}
	if !bytes.Equal(deposits[1023].Data.WithdrawalCredentials, expectedWithdrawalCredentialsAt1023) {
		t.Fatalf("incorrect withdrawal credentials, wanted %x but received %x", expectedWithdrawalCredentialsAt1023, deposits[1023].Data.WithdrawalCredentials)
	}
	dRootAt1023 := []byte("26c644a199d7c609d0668e8e7c2eeb6cbcb4428da32fa8367c30890e68340782")
	dRootAt1023B := make([]byte, hex.DecodedLen(len(dRootAt1023)))
	_, err = hex.Decode(dRootAt1023B, dRootAt1023)
	require.NoError(t, err)
	if !bytes.Equal(depositDataRoots[1023][:], dRootAt1023B) {
		t.Fatalf("incorrect deposit data root, wanted %x but received %x", dRootAt1023B, depositDataRoots[1023])
	}
	sigAt1023 := []byte("a2803bba22be7cc3055bb4769e30526204b47992d0ce2d227ba500df87ee706e5efc8cf22995d127b1ad8ae090daa8d90b662e6b30a0c496a0ab3c331d13a52fd2fcae2769c05fd01bb30e618682c7526637e7cccd500cfeba9b7e5371064ede")
	sigAt1023B := make([]byte, hex.DecodedLen(len(sigAt1023)))
	_, err = hex.Decode(sigAt1023B, sigAt1023)
	require.NoError(t, err)
	if !bytes.Equal(deposits[1023].Data.Signature, sigAt1023B) {
		t.Fatalf("incorrect signature, wanted %x but received %x", sigAt1023B, deposits[1023].Data.Signature)
	}
}

func TestDepositsWithBalance_MatchesDeterministic(t *testing.T) {
	entries := 64
	resetCache()
	balances := make([]uint64, entries)
	for i := 0; i < entries; i++ {
		balances[i] = params.BeaconConfig().MaxEffectiveBalance
	}
	deposits, depositTrie, err := DepositsWithBalance(balances)
	require.NoError(t, err)
	_, depositDataRoots, err := DepositTrieSubset(depositTrie, entries)
	require.NoError(t, err)

	determDeposits, _, err := DeterministicDepositsAndKeys(uint64(entries))
	require.NoError(t, err)
	_, determDepositDataRoots, err := DeterministicDepositTrie(entries)
	require.NoError(t, err)

	for i := 0; i < entries; i++ {
		if !proto.Equal(deposits[i], determDeposits[i]) {
			t.Errorf("Expected deposit %d to match", i)
		}
		if !bytes.Equal(depositDataRoots[i][:], determDepositDataRoots[i][:]) {
			t.Errorf("Expected deposit root %d to match", i)
		}
	}
}

func TestDepositsWithBalance_MatchesDeterministic_Cached(t *testing.T) {
	entries := 32
	resetCache()
	// Cache half of the deposit cache.
	_, _, err := DeterministicDepositsAndKeys(uint64(entries))
	require.NoError(t, err)
	_, _, err = DeterministicDepositTrie(entries)
	require.NoError(t, err)

	// Generate balanced deposits with half cache.
	entries = 64
	balances := make([]uint64, entries)
	for i := 0; i < entries; i++ {
		balances[i] = params.BeaconConfig().MaxEffectiveBalance
	}
	deposits, depositTrie, err := DepositsWithBalance(balances)
	require.NoError(t, err)
	_, depositDataRoots, err := DepositTrieSubset(depositTrie, entries)
	require.NoError(t, err)

	// Get 64 standard deposits.
	determDeposits, _, err := DeterministicDepositsAndKeys(uint64(entries))
	require.NoError(t, err)
	_, determDepositDataRoots, err := DeterministicDepositTrie(entries)
	require.NoError(t, err)

	for i := 0; i < entries; i++ {
		if !proto.Equal(deposits[i], determDeposits[i]) {
			t.Errorf("Expected deposit %d to match", i)
		}
		if !bytes.Equal(depositDataRoots[i][:], determDepositDataRoots[i][:]) {
			t.Errorf("Expected deposit root %d to match", i)
		}
	}
}

func TestSetupInitialDeposits_1024Entries_PartialDeposits(t *testing.T) {
	entries := 1
	resetCache()
	balances := make([]uint64, entries)
	for i := 0; i < entries; i++ {
		balances[i] = params.BeaconConfig().MaxEffectiveBalance / 2
	}
	deposits, depositTrie, err := DepositsWithBalance(balances)
	require.NoError(t, err)
	_, depositDataRoots, err := DepositTrieSubset(depositTrie, entries)
	require.NoError(t, err)

	if len(deposits) != entries {
		t.Fatalf("incorrect number of deposits returned, wanted %d but received %d", entries, len(deposits))
	}
	expectedPublicKeyAt0 := []byte{0xa9, 0x9a, 0x76, 0xed, 0x77, 0x96, 0xf7, 0xbe, 0x22, 0xd5, 0xb7, 0xe8, 0x5d, 0xee, 0xb7, 0xc5, 0x67, 0x7e, 0x88, 0xe5, 0x11, 0xe0, 0xb3, 0x37, 0x61, 0x8f, 0x8c, 0x4e, 0xb6, 0x13, 0x49, 0xb4, 0xbf, 0x2d, 0x15, 0x3f, 0x64, 0x9f, 0x7b, 0x53, 0x35, 0x9f, 0xe8, 0xb9, 0x4a, 0x38, 0xe4, 0x4c}
	if !bytes.Equal(deposits[0].Data.PublicKey, expectedPublicKeyAt0) {
		t.Fatalf("incorrect public key, wanted %x but received %x", expectedPublicKeyAt0, deposits[0].Data.PublicKey)
	}
	expectedWithdrawalCredentialsAt0 := []byte{0x00, 0xec, 0x7e, 0xf7, 0x78, 0x0c, 0x9d, 0x15, 0x15, 0x97, 0x92, 0x40, 0x36, 0x26, 0x2d, 0xd2, 0x8d, 0xc6, 0x0e, 0x12, 0x28, 0xf4, 0xda, 0x6f, 0xec, 0xf9, 0xd4, 0x02, 0xcb, 0x3f, 0x35, 0x94}
	if !bytes.Equal(deposits[0].Data.WithdrawalCredentials, expectedWithdrawalCredentialsAt0) {
		t.Fatalf("incorrect withdrawal credentials, wanted %x but received %x", expectedWithdrawalCredentialsAt0, deposits[0].Data.WithdrawalCredentials)
	}
	dRootAt0 := []byte("d6f1724f93965872fc89f18341106bba6c77bea0e87d04281ffdc0045923f011")
	dRootAt0B := make([]byte, hex.DecodedLen(len(dRootAt0)))
	_, err = hex.Decode(dRootAt0B, dRootAt0)
	require.NoError(t, err)
	if !bytes.Equal(depositDataRoots[0][:], dRootAt0B) {
		t.Fatalf("incorrect deposit data root, wanted %#x but received %#x", dRootAt0B, depositDataRoots[0])
	}

	sigAt0 := []byte("9753ac49adf4bad5a97f3b8dfa6a61a0bad1dd911470e2ea8eab4accf1e20b92427e0266d8654caca8312fe8bbafa72911e835f4807e6c4c40be04dfe09146b1bc96fee6214c4ffe0a9d3bb609f63635a173520d878591e1af8454aaccaa4b36")
	sigAt0B := make([]byte, hex.DecodedLen(len(sigAt0)))
	_, err = hex.Decode(sigAt0B, sigAt0)
	require.NoError(t, err)
	if !bytes.Equal(deposits[0].Data.Signature, sigAt0B) {
		t.Fatalf("incorrect signature, wanted %#x but received %#x", sigAt0B, deposits[0].Data.Signature)
	}

	entries = 1024
	resetCache()
	balances = make([]uint64, entries)
	for i := 0; i < entries; i++ {
		balances[i] = params.BeaconConfig().MaxEffectiveBalance / 2
	}
	deposits, depositTrie, err = DepositsWithBalance(balances)
	require.NoError(t, err)
	_, depositDataRoots, err = DepositTrieSubset(depositTrie, entries)
	require.NoError(t, err)
	if len(deposits) != entries {
		t.Fatalf("incorrect number of deposits returned, wanted %d but received %d", entries, len(deposits))
	}
	// Ensure 0  has not changed
	if !bytes.Equal(deposits[0].Data.PublicKey, expectedPublicKeyAt0) {
		t.Fatalf("incorrect public key, wanted %x but received %x", expectedPublicKeyAt0, deposits[0].Data.PublicKey)
	}
	if !bytes.Equal(deposits[0].Data.WithdrawalCredentials, expectedWithdrawalCredentialsAt0) {
		t.Fatalf("incorrect withdrawal credentials, wanted %x but received %x", expectedWithdrawalCredentialsAt0, deposits[0].Data.WithdrawalCredentials)
	}
	if !bytes.Equal(depositDataRoots[0][:], dRootAt0B) {
		t.Fatalf("incorrect deposit data root, wanted %x but received %x", dRootAt0B, depositDataRoots[0])
	}
	if !bytes.Equal(deposits[0].Data.Signature, sigAt0B) {
		t.Fatalf("incorrect signature, wanted %x but received %x", sigAt0B, deposits[0].Data.Signature)
	}
	expectedPublicKeyAt1023 := []byte{0x81, 0x2b, 0x93, 0x5e, 0xc8, 0x4b, 0x0e, 0x9a, 0x83, 0x95, 0x55, 0xaf, 0x33, 0x60, 0xca, 0xfb, 0x83, 0x1b, 0xd6, 0x12, 0xcf, 0xa2, 0x2e, 0x25, 0xea, 0xb0, 0x3c, 0xf5, 0xfd, 0xb0, 0x2a, 0xf5, 0x2b, 0xa4, 0x01, 0x7a, 0xee, 0xa8, 0x8a, 0x2f, 0x62, 0x2c, 0x78, 0x6e, 0x7f, 0x47, 0x6f, 0x4b}
	if !bytes.Equal(deposits[1023].Data.PublicKey, expectedPublicKeyAt1023) {
		t.Fatalf("incorrect public key, wanted %x but received %x", expectedPublicKeyAt1023, deposits[1023].Data.PublicKey)
	}
	expectedWithdrawalCredentialsAt1023 := []byte{0x00, 0x23, 0xd5, 0x76, 0xbc, 0x6c, 0x15, 0xdb, 0xc4, 0x34, 0x70, 0x1f, 0x3f, 0x41, 0xfd, 0x3e, 0x67, 0x59, 0xd2, 0xea, 0x7c, 0xdc, 0x64, 0x71, 0x0e, 0xe2, 0x8d, 0xde, 0xf7, 0xd2, 0xda, 0x28}
	if !bytes.Equal(deposits[1023].Data.WithdrawalCredentials, expectedWithdrawalCredentialsAt1023) {
		t.Fatalf("incorrect withdrawal credentials, wanted %x but received %x", expectedWithdrawalCredentialsAt1023, deposits[1023].Data.WithdrawalCredentials)
	}
	dRootAt1023 := []byte("16ea6bbce0b12765917734f9001219bccf54d6524742149475d2d6a78963ecf9")
	dRootAt1023B := make([]byte, hex.DecodedLen(len(dRootAt1023)))
	_, err = hex.Decode(dRootAt1023B, dRootAt1023)
	require.NoError(t, err)
	if !bytes.Equal(depositDataRoots[1023][:], dRootAt1023B) {
		t.Fatalf("incorrect deposit data root, wanted %#x but received %#x", dRootAt1023B, depositDataRoots[1023])
	}
	sigAt1023 := []byte("8805d4280f7be9dc0d8b3cb485b73566aecef5250bc82f3073f5baadd7ac8f81ce4f0783844170ca41c1ec50eacb979a0af7b0fce87a9a4610fb448b83d16e3e16b862339964dd7c6559ed3a414e9baf1d8344c1a5578b597c31dc8598906939")
	sigAt1023B := make([]byte, hex.DecodedLen(len(sigAt1023)))
	_, err = hex.Decode(sigAt1023B, sigAt1023)
	require.NoError(t, err)
	if !bytes.Equal(deposits[1023].Data.Signature, sigAt1023B) {
		t.Fatalf("incorrect signature, wanted %#x but received %#x", sigAt1023B, deposits[1023].Data.Signature)
	}
}

func TestDeterministicGenesisState_100Validators(t *testing.T) {
	validatorCount := uint64(100)
	beaconState, privKeys := DeterministicGenesisState(t, validatorCount)
	activeValidators, err := helpers.ActiveValidatorCount(context.Background(), beaconState, 0)
	require.NoError(t, err)

	// lint:ignore uintcast -- test code
	if len(privKeys) != int(validatorCount) {
		t.Fatalf("expected amount of private keys %d to match requested amount of validators %d", len(privKeys), validatorCount)
	}
	if activeValidators != validatorCount {
		t.Fatalf("expected validators in state %d to match requested amount %d", activeValidators, validatorCount)
	}
}

func TestDepositTrieFromDeposits(t *testing.T) {
	deposits, _, err := DeterministicDepositsAndKeys(100)
	require.NoError(t, err)
	eth1Data, err := DeterministicEth1Data(len(deposits))
	require.NoError(t, err)

	depositTrie, _, err := DepositTrieFromDeposits(deposits)
	require.NoError(t, err)

	root, err := depositTrie.HashTreeRoot()
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(root[:], eth1Data.DepositRoot) {
		t.Fatal("expected deposit trie root to equal eth1data deposit root")
	}
}

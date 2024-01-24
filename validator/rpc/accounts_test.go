package rpc

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/prysmaticlabs/prysm/v4/cmd/validator/flags"
	"github.com/prysmaticlabs/prysm/v4/crypto/aes"
	ethpbservice "github.com/prysmaticlabs/prysm/v4/proto/eth/service"
	pb "github.com/prysmaticlabs/prysm/v4/proto/prysm/v1alpha1/validator-client"
	"github.com/prysmaticlabs/prysm/v4/testing/assert"
	"github.com/prysmaticlabs/prysm/v4/testing/require"
	"github.com/prysmaticlabs/prysm/v4/validator/accounts"
	"github.com/prysmaticlabs/prysm/v4/validator/accounts/iface"
	mock "github.com/prysmaticlabs/prysm/v4/validator/accounts/testing"
	"github.com/prysmaticlabs/prysm/v4/validator/client"
	"github.com/prysmaticlabs/prysm/v4/validator/keymanager"
	"github.com/prysmaticlabs/prysm/v4/validator/keymanager/derived"
	constant "github.com/prysmaticlabs/prysm/v4/validator/testing"
)

var defaultWalletPath = filepath.Join(flags.DefaultValidatorDir(), flags.WalletDefaultDirName)

func TestServer_ListAccounts(t *testing.T) {
	ctx := context.Background()
	localWalletDir := setupWalletDir(t)
	defaultWalletPath = localWalletDir
	// We attempt to create the wallet.
	opts := []accounts.Option{
		accounts.WithWalletDir(defaultWalletPath),
		accounts.WithKeymanagerType(keymanager.Derived),
		accounts.WithWalletPassword(strongPass),
		accounts.WithSkipMnemonicConfirm(true),
	}
	acc, err := accounts.NewCLIManager(opts...)
	require.NoError(t, err)
	w, err := acc.WalletCreate(ctx)
	require.NoError(t, err)
	km, err := w.InitializeKeymanager(ctx, iface.InitKeymanagerConfig{ListenForChanges: false})
	require.NoError(t, err)
	vs, err := client.NewValidatorService(ctx, &client.Config{
		Wallet: w,
		Validator: &mock.MockValidator{
			Km: km,
		},
	})
	require.NoError(t, err)
	s := &Server{
		walletInitialized: true,
		wallet:            w,
		validatorService:  vs,
	}
	numAccounts := 50
	dr, ok := km.(*derived.Keymanager)
	require.Equal(t, true, ok)
	err = dr.RecoverAccountsFromMnemonic(ctx, constant.TestMnemonic, derived.DefaultMnemonicLanguage, "", numAccounts)
	require.NoError(t, err)
	resp, err := s.ListAccounts(ctx, &pb.ListAccountsRequest{
		PageSize: int32(numAccounts),
	})
	require.NoError(t, err)
	require.Equal(t, len(resp.Accounts), numAccounts)

	tests := []struct {
		req *pb.ListAccountsRequest
		res *pb.ListAccountsResponse
	}{
		{
			req: &pb.ListAccountsRequest{
				PageSize: 5,
			},
			res: &pb.ListAccountsResponse{
				Accounts:      resp.Accounts[0:5],
				NextPageToken: "1",
				TotalSize:     int32(numAccounts),
			},
		},
		{
			req: &pb.ListAccountsRequest{
				PageSize:  5,
				PageToken: "1",
			},
			res: &pb.ListAccountsResponse{
				Accounts:      resp.Accounts[5:10],
				NextPageToken: "2",
				TotalSize:     int32(numAccounts),
			},
		},
	}
	for _, test := range tests {
		res, err := s.ListAccounts(context.Background(), test.req)
		require.NoError(t, err)
		assert.DeepEqual(t, res, test.res)
	}
}

// func Test_CreateAccountsAndDepositData(t *testing.T) {
// 	ctx := context.Background()
// 	s := &Server{
// 		walletInitializedFeed: new(event.Feed),
// 	}
// 	// 1. Check Wallet is Opened

// 	req1 := &pb.RecoverAccountsFromWalletRequest{
// 		Password:    "testpassword",
// 		NumAccounts: 10,
// 	}

// 	_, err1 := s.CreateAccountsAndDepositData(ctx, req1)
// 	require.ErrorContains(t, "Wallet is Not Opened", err1)
// 	// 2. Normal Case
// 	testPath := "./testpath"
// 	req2 := &pb.InitializeDerivedWalletRequest{
// 		WalletDir:    testPath,
// 		Password:     "testpassword",
// 		MnemonicLang: "english",
// 	}
// 	_, err2 := s.InitializeDerivedWallet(ctx, req2)
// 	require.NoError(t, err2)

// 	req3 := &pb.RecoverAccountsFromWalletRequest{
// 		Password:    "testpassword",
// 		NumAccounts: 10,
// 	}
// 	_, err3 := s.RecoverAccountsFromWallet(ctx, req3)
// 	require.NoError(t, err3)

// 	// 3. Wrong Password

// 	req5 := &pb.RecoverAccountsFromWalletRequest{
// 		Password:    "wrongpassword",
// 		NumAccounts: 10,
// 	}

// 	_, err5 := s.RecoverAccountsFromWallet(ctx, req5)
// 	require.ErrorContains(t, "Could not recover accounts from wallet", err5)
// }

// func Test_GetDepositData(t *testing.T) {
// 	ctx := context.Background()
// 	s := &Server{
// 		walletInitializedFeed: new(event.Feed),
// 	}

// 	// 1. Test Normal

// 	// create wallet
// 	testPath := "./testpath"
// 	// new path
// 	req1 := &pb.InitializeDerivedWalletRequest{
// 		WalletDir:    testPath,
// 		Password:     "testpassword",
// 		MnemonicLang: "english",
// 	}
// 	res1, err1 := s.InitializeDerivedWallet(ctx, req1)
// 	require.NoError(t, err1)
// 	assert.Equal(t, testPath, res1.WalletDir)
// 	km, err := s.wallet.InitializeKeymanager(ctx, iface.InitKeymanagerConfig{ListenForChanges: false})
// 	require.NoError(t, err)
// 	vs, err := client.NewValidatorService(ctx, &client.Config{
// 		Wallet: s.wallet,
// 		Validator: &mock.MockValidator{
// 			Km: km,
// 		},
// 	})
// 	require.NoError(t, err)
// 	s.validatorService = vs
// 	req2 := &pb.RecoverAccountsFromWalletRequest{
// 		Password:    "testpassword",
// 		NumAccounts: 10,
// 	}
// 	_, err2 := s.RecoverAccountsFromWallet(ctx, req2)
// 	require.NoError(t, err2)

// 	req3 := &pb.ListAccountsRequest{
// 		All: true,
// 	}
// 	res3, err3 := s.ListAccounts(ctx, req3)
// 	require.NoError(t, err3)
// 	fmt.Println(len(res3.Accounts))

// 	keys := make([]*pb.DepositDataRequest, 1)
// 	keys[0] = &pb.DepositDataRequest{
// 		PublicKey:   res3.Accounts[0].ValidatingPublicKey,
// 		WithdrawKey: []byte("0xAb5801a7D398351b8bE11C439e05C5B3259aeC9B"),
// 		AmountGwei:  32000000000,
// 	}

// 	req4 := &pb.GetDepositDataRequest{
// 		DepositMessages: keys,
// 	}

// 	res4, err4 := s.GetDepositData(ctx, req4)
// 	require.NoError(t, err4)
// 	fmt.Println(res4.DepositDatas[0])

// 	// 2. Test When validator service is not initialized

// 	// 3. Test When wallet not initialized

// 	// 4. Test When pubkey and withdraw key length is not same

// 	// 5. Test When pubkey is not in our keymanager
// }

func TestImportAccount(t *testing.T) {
	// initialize setting
	ctx := context.Background()
	localWalletDir := setupWalletDir(t)
	defaultWalletPath = localWalletDir

	opts := []accounts.Option{
		accounts.WithWalletDir(defaultWalletPath),
		accounts.WithKeymanagerType(keymanager.Derived),
		accounts.WithWalletPassword(strongPass),
		accounts.WithSkipMnemonicConfirm(true),
	}

	acc, err := accounts.NewCLIManager(opts...)
	require.NoError(t, err)
	w, err := acc.WalletCreate(ctx)
	require.NoError(t, err)
	km, err := w.InitializeKeymanager(ctx, iface.InitKeymanagerConfig{ListenForChanges: false})
	require.NoError(t, err)
	vs, err := client.NewValidatorService(ctx, &client.Config{
		Wallet: w,
		Validator: &mock.MockValidator{
			Km: km,
		},
	})
	require.NoError(t, err)

	// Generate a cipher key
	cipher, err := generateRandomKey()
	require.NoError(t, err)

	// Encrypt the private key using the cipher key
	privateKey, err := hexutil.Decode("0x13f5347b11740bfaf3625d766d6f55762b68fcffa79f857e5ee8b731831cb4d3")
	require.NoError(t, err)

	encryptedPrivateKey, err := aes.Encrypt(cipher, privateKey)
	require.NoError(t, err)

	s := Server{
		walletInitialized: true,
		isOverNode:        true,
		validatorService:  vs,
		wallet:            w,
		cipherKey:         cipher,
	}

	// Create a mock request
	encryptedPrivateKeys := []string{hexutil.Encode(encryptedPrivateKey)} // Use the encrypted private key

	req := &pb.ImportAccountsRequest{
		PrivateKeys: encryptedPrivateKeys,
	}
	expected := &pb.ImportAccountsResponse{
		Data: []*pb.ImportKeystoreStatus{
			&pb.ImportKeystoreStatus{
				PublicKey: "0xabbca58352c3ca208aa72a0231aa16b22015f3309d22d5b82ea1fa27137a00499025d7162bdd822c85856338120ce052",
				Status: &ethpbservice.ImportedKeystoreStatus{
					Status:  ethpbservice.ImportedKeystoreStatus_IMPORTED,
					Message: "",
				},
			},
		},
	}

	actual, err := s.ImportAccounts(ctx, req)
	require.NoError(t, err)
	assert.DeepEqual(t, expected, actual, "handler returned unexpected body")
}

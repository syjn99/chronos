package rpc

import (
	"context"
	rd "crypto/rand"
	"path/filepath"
	"sync"
	"testing"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/prysmaticlabs/prysm/v4/async/event"
	"github.com/prysmaticlabs/prysm/v4/crypto/aes"
	pb "github.com/prysmaticlabs/prysm/v4/proto/prysm/v1alpha1/validator-client"
	"github.com/prysmaticlabs/prysm/v4/testing/assert"
	"github.com/prysmaticlabs/prysm/v4/testing/require"
	mock "github.com/prysmaticlabs/prysm/v4/validator/accounts/testing"
	"github.com/prysmaticlabs/prysm/v4/validator/accounts/wallet"
	"github.com/prysmaticlabs/prysm/v4/validator/client"
	"github.com/prysmaticlabs/prysm/v4/validator/keymanager"
	"github.com/prysmaticlabs/prysm/v4/validator/keymanager/local"
)

const strongPass = "29384283xasjasd32%%&*@*#*"

//func TestServer_CreateWallet_Local(t *testing.T) {
//	ctx := context.Background()
//	localWalletDir := setupWalletDir(t)
//	defaultWalletPath = localWalletDir
//	opts := []accounts.Option{
//		accounts.WithWalletDir(defaultWalletPath),
//		accounts.WithKeymanagerType(keymanager.Derived),
//		accounts.WithWalletPassword(strongPass),
//		accounts.WithSkipMnemonicConfirm(true),
//	}
//	acc, err := accounts.NewCLIManager(opts...)
//	require.NoError(t, err)
//	w, err := acc.WalletCreate(ctx)
//	require.NoError(t, err)
//	km, err := w.InitializeKeymanager(ctx, iface.InitKeymanagerConfig{ListenForChanges: false})
//	require.NoError(t, err)
//	vs, err := client.NewValidatorService(ctx, &client.Config{
//		Wallet: w,
//		Validator: &mock.MockValidator{
//			Km: km,
//		},
//	})
//	require.NoError(t, err)
//	s := &Server{
//		walletInitializedFeed: new(event.Feed),
//		walletDir:             defaultWalletPath,
//		validatorService:      vs,
//	}
//	req := &pb.CreateWalletRequest{
//		Keymanager:     pb.KeymanagerKind_IMPORTED,
//		WalletPassword: strongPass,
//	}
//	_, err = s.CreateWallet(ctx, req)
//	require.NoError(t, err)
//
//	numKeystores := 5
//	password := "12345678"
//	encodedKeystores := make([]string, numKeystores)
//	passwords := make([]string, numKeystores)
//	for i := 0; i < numKeystores; i++ {
//		enc, err := json.Marshal(createRandomKeystore(t, password))
//		encodedKeystores[i] = string(enc)
//		require.NoError(t, err)
//		passwords[i] = password
//	}
//
//	importReq := &ethpbservice.ImportKeystoresRequest{
//		Keystores: encodedKeystores,
//		Passwords: passwords,
//	}
//
//	encryptor := keystorev4.New()
//	keystores := make([]string, 3)
//	for i := 0; i < len(keystores); i++ {
//		privKey, err := bls.RandKey()
//		require.NoError(t, err)
//		pubKey := fmt.Sprintf("%x", privKey.PublicKey().Marshal())
//		id, err := uuid.NewRandom()
//		require.NoError(t, err)
//		cryptoFields, err := encryptor.Encrypt(privKey.Marshal(), strongPass)
//		require.NoError(t, err)
//		item := &keymanager.Keystore{
//			Crypto:      cryptoFields,
//			ID:          id.String(),
//			Version:     encryptor.Version(),
//			Pubkey:      pubKey,
//			Description: encryptor.Name(),
//		}
//		encodedFile, err := json.MarshalIndent(item, "", "\t")
//		require.NoError(t, err)
//		keystores[i] = string(encodedFile)
//	}
//	importReq.Keystores = keystores
//	_, err = s.ImportKeystores(ctx, importReq)
//	require.NoError(t, err)
//}
//
//func TestServer_CreateWallet_Local_PasswordTooWeak(t *testing.T) {
//	localWalletDir := setupWalletDir(t)
//	defaultWalletPath = localWalletDir
//	ctx := context.Background()
//	s := &Server{
//		walletInitializedFeed: new(event.Feed),
//		walletDir:             defaultWalletPath,
//	}
//	req := &pb.CreateWalletRequest{
//		Keymanager:     pb.KeymanagerKind_IMPORTED,
//		WalletPassword: "", // Weak password, empty string
//	}
//	_, err := s.CreateWallet(ctx, req)
//	require.ErrorContains(t, "Password too weak", err)
//
//	req = &pb.CreateWalletRequest{
//		Keymanager:     pb.KeymanagerKind_IMPORTED,
//		WalletPassword: "a", // Weak password, too short
//	}
//	_, err = s.CreateWallet(ctx, req)
//	require.ErrorContains(t, "Password too weak", err)
//}
//
//func TestServer_RecoverWallet_Derived(t *testing.T) {
//	localWalletDir := setupWalletDir(t)
//	ctx := context.Background()
//	s := &Server{
//		walletInitializedFeed: new(event.Feed),
//		walletDir:             localWalletDir,
//	}
//	req := &pb.RecoverWalletRequest{
//		WalletPassword: strongPass,
//		NumAccounts:    0,
//	}
//	_, err := s.RecoverWallet(ctx, req)
//	require.ErrorContains(t, "Must create at least 1 validator account", err)
//
//	req.NumAccounts = 2
//	req.Language = "Swahili"
//	_, err = s.RecoverWallet(ctx, req)
//	require.ErrorContains(t, "input not in the list of supported languages", err)
//
//	req.Language = "ENglish"
//	_, err = s.RecoverWallet(ctx, req)
//	require.ErrorContains(t, "invalid mnemonic in request", err)
//
//	mnemonicRandomness := make([]byte, 32)
//	_, err = rand.NewGenerator().Read(mnemonicRandomness)
//	require.NoError(t, err)
//	mnemonic, err := bip39.NewMnemonic(mnemonicRandomness)
//	require.NoError(t, err)
//	req.Mnemonic = mnemonic
//
//	req.Mnemonic25ThWord = " "
//	_, err = s.RecoverWallet(ctx, req)
//	require.ErrorContains(t, "mnemonic 25th word cannot be empty", err)
//	req.Mnemonic25ThWord = "outer"
//
//	// Test weak password.
//	req.WalletPassword = "123qwe"
//	_, err = s.RecoverWallet(ctx, req)
//	require.ErrorContains(t, "password did not pass validation", err)
//
//	req.WalletPassword = strongPass
//	// Create(derived) should fail then test recover.
//	reqCreate := &pb.CreateWalletRequest{
//		Keymanager:     pb.KeymanagerKind_DERIVED,
//		WalletPassword: strongPass,
//		NumAccounts:    2,
//		Mnemonic:       mnemonic,
//	}
//	_, err = s.CreateWallet(ctx, reqCreate)
//	require.ErrorContains(t, "create wallet not supported through web", err, "Create wallet for DERIVED or REMOTE types not supported through web, either import keystore or recover")
//
//	// This defer will be the last to execute in this func.
//	resetCfgFalse := features.InitWithReset(&features.Flags{
//		WriteWalletPasswordOnWebOnboarding: false,
//	})
//	defer resetCfgFalse()
//
//	resetCfgTrue := features.InitWithReset(&features.Flags{
//		WriteWalletPasswordOnWebOnboarding: true,
//	})
//	defer resetCfgTrue()
//
//	// Finally test recover.
//	_, err = s.RecoverWallet(ctx, req)
//	require.NoError(t, err)
//
//	// Password File should have been written.
//	passwordFilePath := filepath.Join(localWalletDir, wallet.DefaultWalletPasswordFile)
//	assert.Equal(t, true, file.FileExists(passwordFilePath))
//
//	// Attempting to write again should trigger an error.
//	err = writeWalletPasswordToDisk(localWalletDir, "somepassword")
//	require.ErrorContains(t, "cannot write wallet password file as it already exists", err)
//}
//
//func TestServer_ValidateKeystores_FailedPreconditions(t *testing.T) {
//	ctx := context.Background()
//	strongPass := "29384283xasjasd32%%&*@*#*"
//	ss := &Server{}
//	_, err := ss.ValidateKeystores(ctx, &pb.ValidateKeystoresRequest{})
//	assert.ErrorContains(t, "Password required for keystores", err)
//	_, err = ss.ValidateKeystores(ctx, &pb.ValidateKeystoresRequest{
//		KeystoresPassword: strongPass,
//	})
//	assert.ErrorContains(t, "No keystores included in request", err)
//	_, err = ss.ValidateKeystores(ctx, &pb.ValidateKeystoresRequest{
//		KeystoresPassword: strongPass,
//		Keystores:         []string{"badjson"},
//	})
//	assert.ErrorContains(t, "Not a valid EIP-2335 keystore", err)
//}
//
//func TestServer_ValidateKeystores_OK(t *testing.T) {
//	ctx := context.Background()
//	strongPass := "29384283xasjasd32%%&*@*#*"
//	ss := &Server{}
//
//	// Create 3 keystores with the strong password.
//	encryptor := keystorev4.New()
//	keystores := make([]string, 3)
//	pubKeys := make([][]byte, 3)
//	for i := 0; i < len(keystores); i++ {
//		privKey, err := bls.RandKey()
//		require.NoError(t, err)
//		pubKey := fmt.Sprintf("%x", privKey.PublicKey().Marshal())
//		id, err := uuid.NewRandom()
//		require.NoError(t, err)
//		cryptoFields, err := encryptor.Encrypt(privKey.Marshal(), strongPass)
//		require.NoError(t, err)
//		item := &keymanager.Keystore{
//			Crypto:      cryptoFields,
//			ID:          id.String(),
//			Version:     encryptor.Version(),
//			Pubkey:      pubKey,
//			Description: encryptor.Name(),
//		}
//		encodedFile, err := json.MarshalIndent(item, "", "\t")
//		require.NoError(t, err)
//		keystores[i] = string(encodedFile)
//		pubKeys[i] = privKey.PublicKey().Marshal()
//	}
//
//	// Validate the keystores and ensure no error.
//	_, err := ss.ValidateKeystores(ctx, &pb.ValidateKeystoresRequest{
//		KeystoresPassword: strongPass,
//		Keystores:         keystores,
//	})
//	require.NoError(t, err)
//
//	// Check that using a different password will return an error.
//	_, err = ss.ValidateKeystores(ctx, &pb.ValidateKeystoresRequest{
//		KeystoresPassword: "badpassword",
//		Keystores:         keystores,
//	})
//	require.ErrorContains(t, "is incorrect", err)
//
//	// Add a new keystore that was encrypted with a different password and expect
//	// a failure from the function.
//	differentPassword := "differentkeystorepass"
//	privKey, err := bls.RandKey()
//	require.NoError(t, err)
//	pubKey := "somepubkey"
//	id, err := uuid.NewRandom()
//	require.NoError(t, err)
//	cryptoFields, err := encryptor.Encrypt(privKey.Marshal(), differentPassword)
//	require.NoError(t, err)
//	item := &keymanager.Keystore{
//		Crypto:      cryptoFields,
//		ID:          id.String(),
//		Version:     encryptor.Version(),
//		Pubkey:      pubKey,
//		Description: encryptor.Name(),
//	}
//	encodedFile, err := json.MarshalIndent(item, "", "\t")
//	keystores = append(keystores, string(encodedFile))
//	require.NoError(t, err)
//	_, err = ss.ValidateKeystores(ctx, &pb.ValidateKeystoresRequest{
//		KeystoresPassword: strongPass,
//		Keystores:         keystores,
//	})
//	require.ErrorContains(t, "Password for keystore with public key somepubkey is incorrect", err)
//}
//
//func TestServer_WalletConfig_NoWalletFound(t *testing.T) {
//	s := &Server{}
//	resp, err := s.WalletConfig(context.Background(), &empty.Empty{})
//	require.NoError(t, err)
//	assert.DeepEqual(t, resp, &pb.WalletResponse{})
//}
//
//func TestServer_WalletConfig(t *testing.T) {
//	localWalletDir := setupWalletDir(t)
//	defaultWalletPath = localWalletDir
//	ctx := context.Background()
//	s := &Server{
//		walletInitializedFeed: new(event.Feed),
//		walletDir:             defaultWalletPath,
//	}
//	// We attempt to create the wallet.
//	opts := []accounts.Option{
//		accounts.WithWalletDir(defaultWalletPath),
//		accounts.WithKeymanagerType(keymanager.Local),
//		accounts.WithWalletPassword(strongPass),
//		accounts.WithSkipMnemonicConfirm(true),
//	}
//	acc, err := accounts.NewCLIManager(opts...)
//	require.NoError(t, err)
//	w, err := acc.WalletCreate(ctx)
//	require.NoError(t, err)
//	km, err := w.InitializeKeymanager(ctx, iface.InitKeymanagerConfig{ListenForChanges: false})
//	require.NoError(t, err)
//	s.wallet = w
//	vs, err := client.NewValidatorService(ctx, &client.Config{
//		Wallet: w,
//		Validator: &mock.MockValidator{
//			Km: km,
//		},
//	})
//	require.NoError(t, err)
//	s.validatorService = vs
//	resp, err := s.WalletConfig(ctx, &empty.Empty{})
//	require.NoError(t, err)
//
//	assert.DeepEqual(t, resp, &pb.WalletResponse{
//		WalletPath:     localWalletDir,
//		KeymanagerKind: pb.KeymanagerKind_IMPORTED,
//	})
//}
//
//func Test_writeWalletPasswordToDisk(t *testing.T) {
//	walletDir := setupWalletDir(t)
//	resetCfg := features.InitWithReset(&features.Flags{
//		WriteWalletPasswordOnWebOnboarding: false,
//	})
//	defer resetCfg()
//	err := writeWalletPasswordToDisk(walletDir, "somepassword")
//	require.NoError(t, err)
//
//	// Expected a silent failure if the feature flag is not enabled.
//	passwordFilePath := filepath.Join(walletDir, wallet.DefaultWalletPasswordFile)
//	assert.Equal(t, false, file.FileExists(passwordFilePath))
//	resetCfg = features.InitWithReset(&features.Flags{
//		WriteWalletPasswordOnWebOnboarding: true,
//	})
//	defer resetCfg()
//	err = writeWalletPasswordToDisk(walletDir, "somepassword")
//	require.NoError(t, err)
//
//	// File should have been written.
//	assert.Equal(t, true, file.FileExists(passwordFilePath))
//
//	// Attempting to write again should trigger an error.
//	err = writeWalletPasswordToDisk(walletDir, "somepassword")
//	require.NotNil(t, err)
//}

func Test_InitializeWallet(t *testing.T) {
	ctx := context.Background()
	cipher, err := generateRandomKey()
	require.NoError(t, err)
	vs, err := client.NewValidatorService(ctx, &client.Config{
		Validator: &mock.MockValidator{},
	})

	if err != nil {
		t.Fatal(err)
	}

	s := &Server{
		isOverNode:            true,
		walletInitializedFeed: new(event.Feed),
		cipherKey:             cipher,
		validatorService:      vs,
	}
	password := "testpassword"
	encryptedPassword, err := aes.Encrypt(s.cipherKey, []byte(password))

	require.NoError(t, err)

	// Test case 1. Working case.
	testPath := "./testpath"
	// new path
	req1 := &pb.InitializeWalletRequest{
		WalletDir: testPath,
		Password:  hexutil.Encode(encryptedPassword),
	}
	res1, err1 := s.InitializeWallet(ctx, req1)
	require.NoError(t, err1)
	assert.Equal(t, testPath, res1.WalletDir)

	s.wallet = nil
	s.walletInitialized = false

	// exist and normal path
	req2 := &pb.InitializeWalletRequest{
		WalletDir: testPath,
		Password:  hexutil.Encode(encryptedPassword),
	}
	res2, err2 := s.InitializeWallet(ctx, req2)
	require.NoError(t, err2)
	assert.Equal(t, testPath, res2.WalletDir)

	// // Test case 2. Wallet already opened
	testPath = "./testpath2"
	req3 := &pb.InitializeWalletRequest{
		WalletDir: testPath,
		Password:  string(encryptedPassword),
	}
	_, err3 := s.InitializeWallet(ctx, req3)
	require.ErrorContains(t, "Wallet is Already Opened", err3)

	// Test case 3. Invalid key
	s.wallet = nil
	s.walletInitialized = false

	wrongCipher, err := generateRandomKey()
	require.NoError(t, err)
	wrongEncryptedPassword, err := aes.Encrypt(wrongCipher, []byte(password))
	require.NoError(t, err)
	testPath = "./testpath"
	req4 := &pb.InitializeWalletRequest{
		WalletDir: testPath,
		Password:  hexutil.Encode(wrongEncryptedPassword),
	}
	_, err4 := s.InitializeWallet(ctx, req4)
	require.ErrorContains(t, "Could not decrypt password", err4)
}

func TestInitializeWallet_ThreadSafe(t *testing.T) {
	ctx := context.Background()
	cipher, err := generateRandomKey()
	require.NoError(t, err)
	vs, err := client.NewValidatorService(ctx, &client.Config{
		Validator: &mock.MockValidator{},
	})

	if err != nil {
		t.Fatal(err)
	}

	// Create a Server instance
	s := &Server{
		isOverNode:            true,
		walletInitializedFeed: new(event.Feed),
		cipherKey:             cipher,
		validatorService:      vs,
	}

	password := "testpassword"
	encryptedPassword, err := aes.Encrypt(s.cipherKey, []byte(password))

	require.NoError(t, err)

	// Test case 1. Working case.
	testPath := "./testpath"
	// new path
	req := &pb.InitializeWalletRequest{
		WalletDir: testPath,
		Password:  hexutil.Encode(encryptedPassword),
	}

	// Create a WaitGroup to wait for all goroutines to finish
	var wg sync.WaitGroup

	// Create a number of goroutines that will call InitializeWallet concurrently
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, err := s.InitializeWallet(context.Background(), req)
			if err != nil {
				assert.ErrorContains(t, "Wallet is Already Opened", err)
			}
		}()
	}

	// Wait for all goroutines to finish
	wg.Wait()

	// Check if the wallet, walletInitialized, and walletDir fields have the expected values
	assert.Equal(t, true, s.walletInitialized)
	assert.Equal(t, testPath, s.walletDir)
}

func TestChangeWalletPassword(t *testing.T) {
	ctx := context.Background()
	cipher, err := hexutil.Decode("0x877d4074dc2eb53f9d67548700159bdde16d673937415fffea94583f56984ef7")
	require.NoError(t, err)
	password := "testpassword"
	encryptedPassword, err := aes.Encrypt(cipher, []byte(password))
	require.NoError(t, err)

	testPath := filepath.Join(t.TempDir(), "wallet")
	wallet := wallet.New(&wallet.Config{
		WalletDir:      testPath,
		KeymanagerKind: keymanager.Local,
		WalletPassword: password,
	})

	km, err := local.NewKeymanager(ctx, &local.SetupConfig{
		Wallet:           wallet,
		ListenForChanges: true,
	})
	require.NoError(t, err)
	keystores := createRandomKeystore(t, password)
	_, err = km.ImportKeystores(ctx, []*keymanager.Keystore{keystores}, []string{password})

	require.NoError(t, err)
	vs, err := client.NewValidatorService(ctx, &client.Config{
		Validator: &mock.MockValidator{
			Km: km,
		},
	})
	require.NoError(t, err)

	if err != nil {
		t.Fatal(err)
	}

	// Create a Server instance
	s := &Server{
		isOverNode:            true,
		walletInitializedFeed: new(event.Feed),
		cipherKey:             cipher,
		validatorService:      vs,
		wallet:                wallet,
		walletInitialized:     true,
	}

	newPassword := "newPassword"
	encryptedNewPassword, err := aes.Encrypt(cipher, []byte(newPassword))
	require.NoError(t, err)

	req3 := &pb.ChangePasswordRequest{
		Password:    hexutil.Encode(encryptedPassword),
		NewPassword: hexutil.Encode(encryptedNewPassword),
	}

	_, err3 := s.ChangeWalletPassword(ctx, req3)
	require.NoError(t, err3)

	// Wrong password
	req4 := &pb.ChangePasswordRequest{
		Password:    hexutil.Encode(encryptedPassword),
		NewPassword: hexutil.Encode(encryptedNewPassword),
	}

	_, err4 := s.ChangeWalletPassword(ctx, req4)
	require.ErrorContains(t, "Old password is not correct", err4)

	// Check password changed
	req5 := &pb.ChangePasswordRequest{
		Password:    hexutil.Encode(encryptedNewPassword),
		NewPassword: hexutil.Encode(encryptedPassword),
	}

	_, err5 := s.ChangeWalletPassword(ctx, req5)
	require.NoError(t, err5)

}

func generateRandomKey() ([]byte, error) {
	key := make([]byte, 32)
	_, err := rd.Read(key)
	if err != nil {
		return nil, err
	}
	return key, nil
}

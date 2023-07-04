package rpc

import (
	"context"
	"encoding/json"
	"path/filepath"
	"strings"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/pkg/errors"
	"github.com/prysmaticlabs/prysm/v4/beacon-chain/core/signing"
	"github.com/prysmaticlabs/prysm/v4/config/features"
	"github.com/prysmaticlabs/prysm/v4/config/params"
	"github.com/prysmaticlabs/prysm/v4/crypto/bls"
	"github.com/prysmaticlabs/prysm/v4/io/file"
	"github.com/prysmaticlabs/prysm/v4/io/prompt"
	ethpb "github.com/prysmaticlabs/prysm/v4/proto/prysm/v1alpha1"
	pb "github.com/prysmaticlabs/prysm/v4/proto/prysm/v1alpha1/validator-client"
	"github.com/prysmaticlabs/prysm/v4/validator/accounts"
	"github.com/prysmaticlabs/prysm/v4/validator/accounts/wallet"
	iface "github.com/prysmaticlabs/prysm/v4/validator/client/iface"
	"github.com/prysmaticlabs/prysm/v4/validator/keymanager"
	"github.com/prysmaticlabs/prysm/v4/validator/keymanager/derived"
	"github.com/tyler-smith/go-bip39"
	"github.com/tyler-smith/go-bip39/wordlists"
	keystorev4 "github.com/wealdtech/go-eth2-wallet-encryptor-keystorev4"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

const (
	checkExistsErrMsg   = "Could not check if wallet exists"
	checkValidityErrMsg = "Could not check if wallet is valid"
	invalidWalletMsg    = "Directory does not contain a valid wallet"
)

// CreateWallet via an API request, allowing a user to save a new
// imported wallet via RPC.
// DEPRECATE: Prysm Web UI and associated endpoints will be fully removed in a future hard fork.
func (s *Server) CreateWallet(ctx context.Context, req *pb.CreateWalletRequest) (*pb.CreateWalletResponse, error) {
	walletDir := s.walletDir
	exists, err := wallet.Exists(walletDir)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Could not check for existing wallet: %v", err)
	}
	if exists {
		if err := s.initializeWallet(ctx, &wallet.Config{
			WalletDir:      walletDir,
			WalletPassword: req.WalletPassword,
		}); err != nil {
			return nil, err
		}
		keymanagerKind := pb.KeymanagerKind_IMPORTED
		switch s.wallet.KeymanagerKind() {
		case keymanager.Derived:
			keymanagerKind = pb.KeymanagerKind_DERIVED
		case keymanager.Web3Signer:
			keymanagerKind = pb.KeymanagerKind_WEB3SIGNER
		}
		return &pb.CreateWalletResponse{
			Wallet: &pb.WalletResponse{
				WalletPath:     walletDir,
				KeymanagerKind: keymanagerKind,
			},
		}, nil
	}
	if err := prompt.ValidatePasswordInput(req.WalletPassword); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Password too weak: %v", err)
	}
	if req.Keymanager == pb.KeymanagerKind_IMPORTED {
		opts := []accounts.Option{
			accounts.WithWalletDir(walletDir),
			accounts.WithKeymanagerType(keymanager.Local),
			accounts.WithWalletPassword(req.WalletPassword),
			accounts.WithSkipMnemonicConfirm(true),
		}
		acc, err := accounts.NewCLIManager(opts...)
		if err != nil {
			return nil, err
		}
		_, err = acc.WalletCreate(ctx)
		if err != nil {
			return nil, err
		}
		if err := s.initializeWallet(ctx, &wallet.Config{
			WalletDir:      walletDir,
			KeymanagerKind: keymanager.Local,
			WalletPassword: req.WalletPassword,
		}); err != nil {
			return nil, err
		}
		if err := writeWalletPasswordToDisk(walletDir, req.WalletPassword); err != nil {
			return nil, status.Error(codes.Internal, "Could not write wallet password to disk")
		}
		return &pb.CreateWalletResponse{
			Wallet: &pb.WalletResponse{
				WalletPath:     walletDir,
				KeymanagerKind: pb.KeymanagerKind_IMPORTED,
			},
		}, nil
	}
	return nil, status.Errorf(codes.InvalidArgument, "Keymanager type %T create wallet not supported through web", req.Keymanager)
}

// WalletConfig returns the wallet's configuration. If no wallet exists, we return an empty response.
// DEPRECATE: Prysm Web UI and associated endpoints will be fully removed in a future hard fork.
func (s *Server) WalletConfig(_ context.Context, _ *empty.Empty) (*pb.WalletResponse, error) {
	exists, err := wallet.Exists(s.walletDir)
	if err != nil {
		return nil, status.Errorf(codes.Internal, checkExistsErrMsg)
	}
	if !exists {
		// If no wallet is found, we simply return an empty response.
		return &pb.WalletResponse{}, nil
	}
	valid, err := wallet.IsValid(s.walletDir)
	if errors.Is(err, wallet.ErrNoWalletFound) {
		return &pb.WalletResponse{}, nil
	}
	if err != nil {
		return nil, status.Errorf(codes.Internal, checkValidityErrMsg)
	}
	if !valid {
		return nil, status.Errorf(codes.FailedPrecondition, invalidWalletMsg)
	}

	if s.wallet == nil || s.validatorService == nil {
		// If no wallet is found, we simply return an empty response.
		return &pb.WalletResponse{}, nil
	}
	var keymanagerKind pb.KeymanagerKind
	switch s.wallet.KeymanagerKind() {
	case keymanager.Derived:
		keymanagerKind = pb.KeymanagerKind_DERIVED
	case keymanager.Local:
		keymanagerKind = pb.KeymanagerKind_IMPORTED
	case keymanager.Web3Signer:
		keymanagerKind = pb.KeymanagerKind_WEB3SIGNER
	}

	return &pb.WalletResponse{
		WalletPath:     s.walletDir,
		KeymanagerKind: keymanagerKind,
	}, nil
}

// RecoverWallet via an API request, allowing a user to recover a derived.
// Generate the seed from the mnemonic + language + 25th passphrase(optional).
// Create N validator keystores from the seed specified by req.NumAccounts.
// Set the wallet password to req.WalletPassword, then create the wallet from
// the provided Mnemonic and return CreateWalletResponse.
// DEPRECATE: Prysm Web UI and associated endpoints will be fully removed in a future hard fork.
func (s *Server) RecoverWallet(ctx context.Context, req *pb.RecoverWalletRequest) (*pb.CreateWalletResponse, error) {
	numAccounts := int(req.NumAccounts)
	if numAccounts == 0 {
		return nil, status.Error(codes.InvalidArgument, "Must create at least 1 validator account")
	}

	// Check validate mnemonic with chosen language
	language := strings.ToLower(req.Language)
	allowedLanguages := map[string][]string{
		"chinese_simplified":  wordlists.ChineseSimplified,
		"chinese_traditional": wordlists.ChineseTraditional,
		"czech":               wordlists.Czech,
		"english":             wordlists.English,
		"french":              wordlists.French,
		"japanese":            wordlists.Japanese,
		"korean":              wordlists.Korean,
		"italian":             wordlists.Italian,
		"spanish":             wordlists.Spanish,
	}
	if _, ok := allowedLanguages[language]; !ok {
		return nil, status.Error(codes.InvalidArgument, "input not in the list of supported languages")
	}
	bip39.SetWordList(allowedLanguages[language])
	mnemonic := req.Mnemonic
	if err := accounts.ValidateMnemonic(mnemonic); err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid mnemonic in request")
	}

	// Check it is not null and not an empty string.
	if req.Mnemonic25ThWord != "" && strings.TrimSpace(req.Mnemonic25ThWord) == "" {
		return nil, status.Error(codes.InvalidArgument, "mnemonic 25th word cannot be empty")
	}

	// Web UI is structured to only write to the default wallet directory
	// accounts.Recoverwallet checks if wallet already exists.
	walletDir := s.walletDir

	// Web UI should check the new and confirmed password are equal.
	walletPassword := req.WalletPassword
	if err := prompt.ValidatePasswordInput(walletPassword); err != nil {
		return nil, status.Error(codes.InvalidArgument, "password did not pass validation")
	}

	opts := []accounts.Option{
		accounts.WithWalletDir(walletDir),
		accounts.WithWalletPassword(walletPassword),
		accounts.WithMnemonic(mnemonic),
		accounts.WithMnemonic25thWord(req.Mnemonic25ThWord),
		accounts.WithNumAccounts(numAccounts),
	}
	acc, err := accounts.NewCLIManager(opts...)
	if err != nil {
		return nil, err
	}
	if _, err := acc.WalletRecover(ctx); err != nil {
		return nil, err
	}
	if err := s.initializeWallet(ctx, &wallet.Config{
		WalletDir:      walletDir,
		KeymanagerKind: keymanager.Derived,
		WalletPassword: walletPassword,
	}); err != nil {
		return nil, err
	}
	if err := writeWalletPasswordToDisk(walletDir, walletPassword); err != nil {
		return nil, status.Error(codes.Internal, "Could not write wallet password to disk")
	}
	return &pb.CreateWalletResponse{
		Wallet: &pb.WalletResponse{
			WalletPath:     walletDir,
			KeymanagerKind: pb.KeymanagerKind_DERIVED,
		},
	}, nil
}

// ValidateKeystores checks whether a set of EIP-2335 keystores in the request
// can indeed be decrypted using a password in the request. If there is no issue,
// we return an empty response with no error. If the password is incorrect for a single keystore,
// we return an appropriate error.
// DEPRECATE: Prysm Web UI and associated endpoints will be fully removed in a future hard fork.
func (*Server) ValidateKeystores(
	_ context.Context, req *pb.ValidateKeystoresRequest,
) (*emptypb.Empty, error) {
	if req.KeystoresPassword == "" {
		return nil, status.Error(codes.InvalidArgument, "Password required for keystores")
	}
	// Needs to unmarshal the keystores from the requests.
	if req.Keystores == nil || len(req.Keystores) < 1 {
		return nil, status.Error(codes.InvalidArgument, "No keystores included in request")
	}
	decryptor := keystorev4.New()
	for i := 0; i < len(req.Keystores); i++ {
		encoded := req.Keystores[i]
		keystore := &keymanager.Keystore{}
		if err := json.Unmarshal([]byte(encoded), &keystore); err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "Not a valid EIP-2335 keystore JSON file: %v", err)
		}
		if _, err := decryptor.Decrypt(keystore.Crypto, req.KeystoresPassword); err != nil {
			doesNotDecrypt := strings.Contains(err.Error(), keymanager.IncorrectPasswordErrMsg)
			if doesNotDecrypt {
				return nil, status.Errorf(
					codes.InvalidArgument,
					"Password for keystore with public key %s is incorrect. "+
						"Prysm web only supports importing batches of keystores with the same password for all of them",
					keystore.Pubkey,
				)
			} else {
				return nil, status.Errorf(codes.Internal, "Unexpected error decrypting keystore: %v", err)
			}
		}
	}

	return &emptypb.Empty{}, nil
}

// Initialize a wallet and send it over a global feed.
// DEPRECATE: Prysm Web UI and associated endpoints will be fully removed in a future hard fork.
func (s *Server) initializeWallet(ctx context.Context, cfg *wallet.Config) error {
	// We first ensure the user has a wallet.
	exists, err := wallet.Exists(cfg.WalletDir)
	if err != nil {
		return errors.Wrap(err, wallet.CheckExistsErrMsg)
	}
	if !exists {
		return wallet.ErrNoWalletFound
	}
	valid, err := wallet.IsValid(cfg.WalletDir)
	if errors.Is(err, wallet.ErrNoWalletFound) {
		return wallet.ErrNoWalletFound
	}
	if err != nil {
		return errors.Wrap(err, wallet.CheckValidityErrMsg)
	}
	if !valid {
		return errors.New(wallet.InvalidWalletErrMsg)
	}

	// We fire an event with the opened wallet over
	// a global feed signifying wallet initialization.
	w, err := wallet.OpenWallet(ctx, &wallet.Config{
		WalletDir:      cfg.WalletDir,
		WalletPassword: cfg.WalletPassword,
	})
	if err != nil {
		return errors.Wrap(err, "could not open wallet")
	}

	s.walletInitialized = true
	s.wallet = w
	s.walletDir = cfg.WalletDir

	s.walletInitializedFeed.Send(w)

	return nil
}

func writeWalletPasswordToDisk(walletDir, password string) error {
	if !features.Get().WriteWalletPasswordOnWebOnboarding {
		return nil
	}
	passwordFilePath := filepath.Join(walletDir, wallet.DefaultWalletPasswordFile)
	if file.FileExists(passwordFilePath) {
	}
	return file.WriteFile(passwordFilePath, []byte(password))
}

//////////////////// PVER

func (s *Server) OpenOrCreateWallet(
	ctx context.Context, req *pb.OpenOrCreateWalletRequest,
) (*pb.OpenOrCreateWalletResponse, error) {
	// Check Wallet is Already Opened
	if s.wallet != nil {
		// Wallet is Already Exist
		return nil, status.Error(codes.AlreadyExists, "Wallet is Already Opened")
	}

	// Check Dir Wallet is Exist
	walletDir := req.WalletDir
	exists, err := wallet.Exists(walletDir)
	if err != nil {
		return nil, status.Error(codes.Internal, "Could not check if wallet exists")
	}

	valid := false
	if exists {
		valid, err = wallet.IsValid(walletDir)
		if err != nil {
			return nil, status.Error(codes.Internal, "Could not check if wallet is valid")
		}
	}

	if exists && valid {
		w, err := wallet.OpenWallet(ctx, &wallet.Config{
			WalletDir:      walletDir,
			WalletPassword: req.Password,
		})
		if err != nil {
			return nil, status.Error(codes.Internal, "Could not open wallet")
		}
		if w.KeymanagerKind() != keymanager.Derived {
			return nil, status.Error(codes.Internal, "Wallet is not a derived keymanager wallet")
		}
		s.wallet = w
		s.walletInitializedFeed.Send(w)
	} else {
		// Create Wallet And Open
		w, err := createDerivedKeymanagerWallet(ctx, req.WalletDir, req.Password, req.MnemonicLaunguage)
		if err != nil {
			return nil, err
		}
		s.wallet = w
		s.walletInitializedFeed.Send(w)
	}
	s.walletInitialized = true
	s.walletDir = req.WalletDir

	return &pb.OpenOrCreateWalletResponse{
		WalletDir: s.walletDir,
	}, nil
}

// Pver
func (s *Server) RecoverAccountsFromWallet(
	ctx context.Context, req *pb.RecoverAccountsFromWalletRequest,
) (*emptypb.Empty, error) {
	// Check Wallet is Opened
	if s.wallet == nil {
		return nil, status.Error(codes.NotFound, "Wallet is Not Opened")
	}

	err := recoverAccountsFromWallet(ctx, s.wallet, req.Password, req.NumAccounts)
	if err != nil {
		return nil, status.Error(codes.Internal, "Could not recover accounts from wallet")
	}

	return &emptypb.Empty{}, nil
}

func (s *Server) GetDepositData(ctx context.Context, req *pb.GetDepositDataRequest) (*pb.GetDepositDataResponse, error) {
	if s.validatorService == nil {
		return nil, status.Error(codes.NotFound, "Validator Service is Not Opened")
	}
	if s.wallet == nil {
		return nil, status.Error(codes.NotFound, "Wallet is Not Opened")
	}
	if len(req.DepositMessages) == 0 {
		return nil, status.Error(codes.InvalidArgument, "Deposit Data Keys is Empty")
	}
	km, err := s.validatorService.Keymanager()
	if err != nil {
		return nil, status.Error(codes.Internal, "Could not get keymanager")
	}
	datas := make([]*pb.DepositData, len(req.DepositMessages))

	for i, key := range req.DepositMessages {
		pubKey, err := bls.PublicKeyFromBytes(key.PublicKey)
		if err != nil {
			return nil, status.Error(codes.Internal, "Could not parse public key")
		}
		dd, err := createDepositData(ctx, pubKey, key.WithdrawKey, key.AmountGwei, km.Sign)
		if err != nil {
			return nil, status.Error(codes.Internal, "Could not create deposit data")
		}
		datas[i] = dd
	}

	return &pb.GetDepositDataResponse{
		DepositDatas: datas,
	}, nil
}

func createDerivedKeymanagerWallet(
	ctx context.Context,
	walletDir string,
	mnemonicPassphrase string,
	mnemonicLanguage string,
) (*wallet.Wallet, error) {
	w := wallet.New(&wallet.Config{
		WalletDir:      walletDir,
		KeymanagerKind: keymanager.Derived,
		WalletPassword: mnemonicPassphrase,
	})

	if err := w.SaveWallet(); err != nil {
		return nil, errors.Wrap(err, "could not save wallet to disk")
	}

	_, err := derived.NewKeymanager(ctx, &derived.SetupConfig{
		Wallet:           w,
		ListenForChanges: true,
	})
	if err != nil {
		return nil, errors.Wrap(err, "could not initialize HD keymanager")
	}

	err = derived.GenerateAndSaveMnemonic(mnemonicLanguage, mnemonicPassphrase, w.AccountsDir())
	if err != nil {
		return nil, errors.Wrap(err, "could not generage and save mnemonic")
	}

	return w, nil
}

func recoverAccountsFromWallet(
	ctx context.Context,
	w *wallet.Wallet,
	mnemonicPassphrase string,
	numAccounts uint64,
) error {
	km, err := derived.NewKeymanager(ctx, &derived.SetupConfig{
		Wallet:           w,
		ListenForChanges: true,
	})
	if err != nil {
		return errors.Wrap(err, "could not make keymanager for given phrase")
	}
	mnemonicStore, err := derived.LoadMnemonic(w.AccountsDir(), mnemonicPassphrase) // TODO : need to change
	if err != nil {
		return errors.Wrap(err, "could not load mnemonic")
	}

	mnemonicLanguage := "english" // TODO : FIXIT

	index := int(mnemonicStore.LatestIndex + numAccounts)

	err = km.RecoverAccountsFromMnemonic(ctx, mnemonicStore.Mnemonic, mnemonicLanguage, mnemonicPassphrase, index)
	if err != nil {
		return err
	}

	err = derived.SaveMnemonicStore(mnemonicStore.Mnemonic, mnemonicPassphrase, w.AccountsDir(), uint64(index))
	if err != nil {
		return err
	}

	return nil
}

func createDepositData(
	ctx context.Context,
	depositPubkey bls.PublicKey,
	eth1WithdrawlAddress []byte,
	amountInGwei uint64,
	signer iface.SigningFunc,
) (*pb.DepositData, error) {
	depositMessage := &ethpb.DepositMessage{
		PublicKey:             depositPubkey.Marshal(),
		WithdrawalCredentials: withdrawalCredentialsHash(eth1WithdrawlAddress),
		Amount:                amountInGwei,
	}

	sr, err := depositMessage.HashTreeRoot()
	if err != nil {
		return nil, err
	}
	domain, err := signing.ComputeDomain(
		params.BeaconConfig().DomainDeposit,
		nil, /*forkVersion*/
		nil, /*genesisValidatorsRoot*/
	)
	if err != nil {
		return nil, err
	}
	root, err := (&ethpb.SigningData{ObjectRoot: sr[:], Domain: domain}).HashTreeRoot()
	if err != nil {
		return nil, err
	}

	sig, err := signer(ctx, &pb.SignRequest{
		PublicKey:   depositPubkey.Marshal(),
		SigningRoot: root[:],
	})
	if err != nil {
		return nil, err
	}

	di := &ethpb.Deposit_Data{
		PublicKey:             depositMessage.PublicKey,
		WithdrawalCredentials: depositMessage.WithdrawalCredentials,
		Amount:                depositMessage.Amount,
		Signature:             sig.Marshal(),
	}

	dr, err := di.HashTreeRoot()
	if err != nil {
		return nil, err
	}
	dd := &pb.DepositData{
		Pubkey:                depositMessage.PublicKey,
		WithdrawalCredentials: depositMessage.WithdrawalCredentials,
		Signature:             sig.Marshal(),
		DepositDataRoot:       dr[:],
	}
	return dd, nil
}

// withdrawal_credentials[:1] == ETH1_ADDRESS_WITHDRAWAL_PREFIX
// withdrawal_credentials[1:12] == b`\x00' * 11`
// withdrawal_credentials[12:] == eth1_withdrawal_address
func withdrawalCredentialsHash(withdrawalAddress []byte) []byte {
	return append(append([]byte{params.BeaconConfig().ETH1AddressWithdrawalPrefixByte}, params.BeaconConfig().ZeroHash[1:12]...), withdrawalAddress[:20]...)[:32]
}

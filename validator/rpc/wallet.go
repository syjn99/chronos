package rpc

import (
	"context"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/pkg/errors"
	"github.com/prysmaticlabs/prysm/v4/crypto/aes"
	pb "github.com/prysmaticlabs/prysm/v4/proto/prysm/v1alpha1/validator-client"
	"github.com/prysmaticlabs/prysm/v4/validator/accounts/wallet"
	"github.com/prysmaticlabs/prysm/v4/validator/keymanager"
	"github.com/prysmaticlabs/prysm/v4/validator/keymanager/derived"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// InitializeDerivedWallet initialize wallet with Derived type keymanager
func (s *Server) InitializeDerivedWallet(
	ctx context.Context, req *pb.InitializeDerivedWalletRequest,
) (*pb.InitializeDerivedWalletResponse, error) {
	if !s.isOverNode {
		return nil, status.Error(codes.NotFound, "Only available when over node flag enabled")
	}
	// initialize derived wallet can only be called once
	if s.wallet != nil {
		return nil, status.Error(codes.AlreadyExists, "Wallet is Already Opened")
	}

	// check if wallet Initialized Event channel is Opened
	if !s.validatorService.IsWaitingKeyManagerInitialization() {
		return nil, status.Error(codes.Unavailable, "Client is not ready to listen wallet initialized event")
	}

	exists, err := wallet.Exists(req.WalletDir)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Could not check for existing wallet: %v", err)
	}
	password, err := hexutil.Decode(req.Password)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "Could not decode password")
	}
	decryptedPassword, err := aes.Decrypt(s.cipherKey, password)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "Could not decrypt password")
	}
	if exists {
		// Open wallet
		w, err := wallet.OpenWallet(ctx, &wallet.Config{
			WalletDir:      req.WalletDir,
			WalletPassword: string(decryptedPassword),
		})
		if err != nil {
			return nil, status.Error(codes.Internal, "Could not open wallet")
		}
		if w.KeymanagerKind() != keymanager.Derived {
			return nil, status.Error(codes.Internal, "Wallet is not a derived keymanager wallet")
		}
		s.wallet = w
	} else {
		// Create wallet and open it
		w, err := createDerivedKeymanagerWallet(ctx, req.WalletDir, string(decryptedPassword), req.MnemonicLang)
		if err != nil {
			return nil, err
		}
		s.wallet = w
	}
	s.walletInitialized = true
	s.walletInitializedFeed.Send(s.wallet)
	s.walletDir = req.WalletDir

	return &pb.InitializeDerivedWalletResponse{
		WalletDir: s.walletDir,
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

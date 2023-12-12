package rpc

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/pkg/errors"
	"github.com/prysmaticlabs/prysm/v4/crypto/aes"
	"github.com/prysmaticlabs/prysm/v4/io/file"
	pb "github.com/prysmaticlabs/prysm/v4/proto/prysm/v1alpha1/validator-client"
	"github.com/prysmaticlabs/prysm/v4/validator/accounts/wallet"
	"github.com/prysmaticlabs/prysm/v4/validator/keymanager"
	"github.com/prysmaticlabs/prysm/v4/validator/keymanager/derived"
	"github.com/prysmaticlabs/prysm/v4/validator/keymanager/local"
	keystorev4 "github.com/wealdtech/go-eth2-wallet-encryptor-keystorev4"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// InitializeDerivedWallet initialize wallet with Derived type keymanager
func (s *Server) InitializeDerivedWallet(
	ctx context.Context, req *pb.InitializeDerivedWalletRequest,
) (*pb.InitializeDerivedWalletResponse, error) {
	s.rpcMutex.Lock()
	defer s.rpcMutex.Unlock()
	if !s.isOverNode {
		log.Debug("InitializeDerivedWallet was called when over node flag disabled")
		return nil, status.Error(codes.NotFound, "Only available when over node flag enabled")
	}
	// initialize derived wallet can only be called once
	if s.wallet != nil {
		log.Debug("InitializeDerivedWallet was called when wallet is already opened")
		return nil, status.Error(codes.AlreadyExists, "Wallet is Already Opened")
	}

	// check if wallet Initialized Event channel is Opened
	if !s.validatorService.IsWaitingKeyManagerInitialization() {
		log.Debug("InitializeDerivedWallet was called when wallet initialized event channel is not opened")
		return nil, status.Error(codes.Unavailable, "Client is not ready to listen wallet initialized event")
	}

	exists, err := wallet.Exists(req.WalletDir)
	if err != nil {
		log.WithError(err).Error("Could not check for existing wallet")
		return nil, status.Errorf(codes.Internal, "Could not check for existing wallet: %v", err)
	}
	password, err := hexutil.Decode(req.Password)
	if err != nil {
		log.WithError(err).Error("Could not decode password")
		return nil, status.Error(codes.InvalidArgument, "Could not decode password")
	}
	decryptedPassword, err := aes.Decrypt(s.cipherKey, password)
	if err != nil {
		log.WithError(err).Error("Could not decrypt password")
		return nil, status.Error(codes.InvalidArgument, "Could not decrypt password")
	}
	if exists {
		// Open wallet
		w, err := wallet.OpenWallet(ctx, &wallet.Config{
			WalletDir:      req.WalletDir,
			WalletPassword: string(decryptedPassword),
		})
		if err != nil {
			log.WithError(err).Error("Could not open wallet")
			return nil, status.Error(codes.Internal, "Could not open wallet")
		}
		if w.KeymanagerKind() != keymanager.Derived {
			log.Error("Wallet is not a derived keymanager wallet")
			return nil, status.Error(codes.Internal, "Wallet is not a derived keymanager wallet")
		}
		// check mnemonic, second check keystore
		chk, err := checkPasswordValid(filepath.Join(w.AccountsDir(), derived.MnemonicStoreFileName), string(decryptedPassword))
		if err == nil {
			_, err = checkPasswordValid(filepath.Join(w.AccountsDir(), local.AccountsPath, local.AccountsKeystoreFileName), string(decryptedPassword))
			if err != nil {
				// keystore file exist and keystore file password incorrect
				log.Error("Password is not correct")
				return nil, status.Error(codes.InvalidArgument, "Password is not correct")
			}
			if !chk {
				// mnemonic file not exist so create it
				err = derived.GenerateAndSaveMnemonic(derived.DefaultMnemonicLanguage, string(decryptedPassword), w.AccountsDir())
				if err != nil {
					log.WithError(err).Error("Could not generate and save mnemonic")
					return nil, status.Error(codes.Internal, "Could not generate and save mnemonic")
				}
			}
		} else {
			//  password incorrect
			log.Error("Password is not correct")
			return nil, status.Error(codes.InvalidArgument, "Password is not correct")
		}
		s.wallet = w
	} else {
		// Create wallet and open it
		w, err := createDerivedKeymanagerWallet(ctx, req.WalletDir, string(decryptedPassword), req.MnemonicLang)
		if err != nil {
			log.WithError(err).Error("Could not create derived keymanager wallet")
			return nil, status.Error(codes.Internal, "Could not create derived keymanager wallet")
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

type KeyStoreRepresent struct {
	Crypto map[string]interface{} `json:"crypto"`
}

// checkPasswordValid check password valid
func checkPasswordValid(path string, password string) (bool, error) {
	if !file.FileExists(path) {
		return false, nil
	}
	rawData, err := os.ReadFile(filepath.Clean(path))
	if err != nil {
		return false, err
	}
	keystoreFile := &KeyStoreRepresent{}
	if err := json.Unmarshal(rawData, keystoreFile); err != nil {
		return false, err
	}

	decryptor := keystorev4.New()
	_, err = decryptor.Decrypt(keystoreFile.Crypto, password)
	if err != nil && strings.Contains(err.Error(), keymanager.IncorrectPasswordErrMsg) {
		return false, err
	} else if err != nil {
		return false, err
	}
	return true, nil
}

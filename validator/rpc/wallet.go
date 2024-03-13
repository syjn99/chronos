package rpc

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"

	"github.com/ethereum/go-ethereum/common/hexutil"
	emptypb "github.com/golang/protobuf/ptypes/empty"
	"github.com/pkg/errors"
	"github.com/prysmaticlabs/prysm/v4/crypto/aes"
	"github.com/prysmaticlabs/prysm/v4/io/file"
	pb "github.com/prysmaticlabs/prysm/v4/proto/prysm/v1alpha1/validator-client"
	"github.com/prysmaticlabs/prysm/v4/validator/accounts/wallet"
	"github.com/prysmaticlabs/prysm/v4/validator/keymanager"
	"github.com/prysmaticlabs/prysm/v4/validator/keymanager/local"
	keystorev4 "github.com/wealdtech/go-eth2-wallet-encryptor-keystorev4"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// InitializeWallet initialize wallet with local type keymanager
func (s *Server) InitializeWallet(
	ctx context.Context, req *pb.InitializeWalletRequest,
) (*pb.InitializeWalletResponse, error) {
	s.rpcMutex.Lock()
	defer s.rpcMutex.Unlock()
	if !s.isOverNode {
		log.Debug("InitializeWallet was called when over node flag disabled")
		return nil, status.Error(codes.NotFound, "Only available when over node flag enabled")
	}
	// initialize derived wallet can only be called once
	if s.wallet != nil {
		log.Debug("InitializeWallet was called when wallet is already opened")
		return nil, status.Error(codes.AlreadyExists, "Wallet is Already Opened")
	}

	// check if wallet Initialized Event channel is Opened
	if !s.validatorService.IsWaitingKeyManagerInitialization() {
		log.Debug("InitializeWallet was called when wallet initialized event channel is not opened")
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
		if w.KeymanagerKind() != keymanager.Local {
			log.Error("Wallet is not a local keymanager wallet")
			return nil, status.Error(codes.Internal, "Wallet is not a local keymanager wallet")
		}
		_, err = checkPasswordValid(filepath.Join(w.AccountsDir(), local.AccountsPath, local.AccountsKeystoreFileName), string(decryptedPassword))
		if err != nil {
			if strings.Contains(err.Error(), keymanager.IncorrectPasswordErrMsg) {
				log.Error("Password is not correct")
				return nil, status.Error(codes.InvalidArgument, "Password is not correct")
			} else {
				log.Error("Could not check password valid", err)
				return nil, status.Error(codes.Internal, "Could not check password valid")
			}
		}
		s.wallet = w
	} else {
		// Create wallet and open it
		w, err := createLocalKeymanagerWallet(ctx, req.WalletDir, string(decryptedPassword))
		if err != nil {
			log.WithError(err).Error("Could not create local keymanager wallet")
			return nil, status.Error(codes.Internal, "Could not create local keymanager wallet")
		}
		s.wallet = w
	}
	s.walletInitialized = true
	s.walletInitializedFeed.Send(s.wallet)
	s.walletDir = req.WalletDir

	return &pb.InitializeWalletResponse{
		WalletDir: s.walletDir,
	}, nil
}

func (s *Server) ChangeWalletPassword(
	ctx context.Context, req *pb.ChangePasswordRequest,
) (*emptypb.Empty, error) {
	s.rpcMutex.Lock()
	defer s.rpcMutex.Unlock()
	if !s.isOverNode {
		log.Debug("ChangeWalletPassword was called when over node flag disabled")
		return nil, status.Error(codes.NotFound, "Only available when over node flag enabled")
	}
	if s.validatorService == nil {
		log.Debug("ChangeWalletPassword was called when validator service is not opened")
		return nil, status.Error(codes.Unavailable, "Validator Service is Not Opened")
	}
	if s.wallet == nil {
		log.Debug("ChangeWalletPassword was called when wallet is not opened")
		return nil, status.Error(codes.Unavailable, "Wallet is Not Opened")
	}

	// Validate old password
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
	if s.wallet.Password() != string(decryptedPassword) {
		log.Error("password is not correct")
		return nil, status.Error(codes.InvalidArgument, "Old password is not correct")
	}
	newPassword, err := hexutil.Decode(req.NewPassword)
	if err != nil {
		log.WithError(err).Error("Could not decode new password")
		return nil, status.Error(codes.InvalidArgument, "Could not decode new password")
	}
	decryptedNewPassword, err := aes.Decrypt(s.cipherKey, newPassword)
	if err != nil {
		log.WithError(err).Error("Could not decrypt new password")
		return nil, status.Error(codes.InvalidArgument, "Could not decrypt new password")
	}
	// get keymanager
	km, err := s.validatorService.Keymanager()
	if err != nil {
		log.WithError(err).Error("Could not get keymanager")
		return nil, status.Error(codes.Internal, "Could not get keymanager")
	}

	// Change Password
	if err := s.wallet.ChangePassword(ctx, km, string(decryptedNewPassword)); err != nil {
		log.WithError(err).Error("Could not change password")
		return nil, status.Error(codes.Internal, "Could not change password")
	}

	return &emptypb.Empty{}, nil
}

// createLocalKeymanagerWallet creates a local keymanager wallet and saves it to disk.
func createLocalKeymanagerWallet(
	ctx context.Context,
	walletDir string,
	mnemonicPassphrase string,
) (*wallet.Wallet, error) {
	w := wallet.New(&wallet.Config{
		WalletDir:      walletDir,
		KeymanagerKind: keymanager.Local,
		WalletPassword: mnemonicPassphrase,
	})

	if err := w.SaveWallet(); err != nil {
		return nil, errors.Wrap(err, "could not save wallet to disk")
	}

	localKm, err := local.NewKeymanager(ctx, &local.SetupConfig{
		Wallet:           w,
		ListenForChanges: false,
	})
	if err != nil {
		return nil, errors.Wrap(err, "could not initialize HD keymanager")
	}
	// make empty accounts keystore json file
	accountsKeystore, err := localKm.CreateAccountsKeystore(ctx, make([][]byte, 0), make([][]byte, 0))
	if err != nil {
		return nil, err
	}
	encodedAccounts, err := json.MarshalIndent(accountsKeystore, "", "\t")
	if err != nil {
		return nil, err
	}
	if err = w.WriteFileAtPath(ctx, local.AccountsPath, local.AccountsKeystoreFileName, encodedAccounts); err != nil {
		return nil, err
	}

	return w, nil
}

type KeyStoreRepresent struct {
	Crypto map[string]interface{} `json:"crypto"`
}

// checkPasswordValid check password valid. return error if password is incorrect.
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

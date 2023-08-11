package rpc

import (
	"archive/zip"
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"
	"github.com/prysmaticlabs/prysm/v4/api/pagination"
	"github.com/prysmaticlabs/prysm/v4/beacon-chain/core/signing"
	"github.com/prysmaticlabs/prysm/v4/cmd"
	"github.com/prysmaticlabs/prysm/v4/config/params"
	"github.com/prysmaticlabs/prysm/v4/crypto/aes"
	"github.com/prysmaticlabs/prysm/v4/crypto/bls"
	ethpb "github.com/prysmaticlabs/prysm/v4/proto/prysm/v1alpha1"
	pb "github.com/prysmaticlabs/prysm/v4/proto/prysm/v1alpha1/validator-client"
	"github.com/prysmaticlabs/prysm/v4/validator/accounts"
	"github.com/prysmaticlabs/prysm/v4/validator/accounts/petnames"
	"github.com/prysmaticlabs/prysm/v4/validator/accounts/wallet"
	iface "github.com/prysmaticlabs/prysm/v4/validator/client/iface"
	"github.com/prysmaticlabs/prysm/v4/validator/keymanager"
	"github.com/prysmaticlabs/prysm/v4/validator/keymanager/derived"
	"github.com/prysmaticlabs/prysm/v4/validator/keymanager/local"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ListAccounts allows retrieval of validating keys and their petnames
// for a user's wallet via RPC.
// DEPRECATED: Prysm Web UI and associated endpoints will be fully removed in a future hard fork.
func (s *Server) ListAccounts(ctx context.Context, req *pb.ListAccountsRequest) (*pb.ListAccountsResponse, error) {
	if s.validatorService == nil {
		return nil, status.Error(codes.FailedPrecondition, "Validator service not yet initialized")
	}
	if !s.walletInitialized {
		return nil, status.Error(codes.FailedPrecondition, "Wallet not yet initialized")
	}
	if int(req.PageSize) > cmd.Get().MaxRPCPageSize {
		return nil, status.Errorf(codes.InvalidArgument, "Requested page size %d can not be greater than max size %d",
			req.PageSize, cmd.Get().MaxRPCPageSize)
	}
	km, err := s.validatorService.Keymanager()
	if err != nil {
		return nil, err
	}
	keys, err := km.FetchValidatingPublicKeys(ctx)
	if err != nil {
		return nil, err
	}
	accs := make([]*pb.Account, len(keys))
	for i := 0; i < len(keys); i++ {
		accs[i] = &pb.Account{
			ValidatingPublicKey: keys[i][:],
			AccountName:         petnames.DeterministicName(keys[i][:], "-"),
		}
		if s.wallet.KeymanagerKind() == keymanager.Derived {
			accs[i].DerivationPath = fmt.Sprintf(derived.ValidatingKeyDerivationPathTemplate, i)
		}
	}
	if req.All {
		return &pb.ListAccountsResponse{
			Accounts:      accs,
			TotalSize:     int32(len(keys)),
			NextPageToken: "",
		}, nil
	}
	start, end, nextPageToken, err := pagination.StartAndEndPage(req.PageToken, int(req.PageSize), len(keys))
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			"Could not paginate results: %v",
			err,
		)
	}
	return &pb.ListAccountsResponse{
		Accounts:      accs[start:end],
		TotalSize:     int32(len(keys)),
		NextPageToken: nextPageToken,
	}, nil
}

// BackupAccounts creates a zip file containing EIP-2335 keystores for the user's
// specified public keys by encrypting them with the specified password.
// DEPRECATED: Prysm Web UI and associated endpoints will be fully removed in a future hard fork.
func (s *Server) BackupAccounts(
	ctx context.Context, req *pb.BackupAccountsRequest,
) (*pb.BackupAccountsResponse, error) {
	if s.validatorService == nil {
		return nil, status.Error(codes.FailedPrecondition, "Validator service not yet initialized")
	}
	if req.PublicKeys == nil || len(req.PublicKeys) < 1 {
		return nil, status.Error(codes.InvalidArgument, "No public keys specified to backup")
	}
	if req.BackupPassword == "" {
		return nil, status.Error(codes.InvalidArgument, "Backup password cannot be empty")
	}

	if s.wallet == nil {
		return nil, status.Error(codes.FailedPrecondition, "No wallet found")
	}
	var err error
	km, err := s.validatorService.Keymanager()
	if err != nil {
		return nil, err
	}
	pubKeys := make([]bls.PublicKey, len(req.PublicKeys))
	for i, key := range req.PublicKeys {
		pubKey, err := bls.PublicKeyFromBytes(key)
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "%#x Not a valid BLS public key: %v", key, err)
		}
		pubKeys[i] = pubKey
	}

	var keystoresToBackup []*keymanager.Keystore
	switch km := km.(type) {
	case *local.Keymanager:
		keystoresToBackup, err = km.ExtractKeystores(ctx, pubKeys, req.BackupPassword)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "Could not backup accounts for local keymanager: %v", err)
		}
	case *derived.Keymanager:
		keystoresToBackup, err = km.ExtractKeystores(ctx, pubKeys, req.BackupPassword)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "Could not backup accounts for derived keymanager: %v", err)
		}
	default:
		return nil, status.Error(codes.FailedPrecondition, "Only HD or imported wallets can backup accounts")
	}
	if len(keystoresToBackup) == 0 {
		return nil, status.Error(codes.InvalidArgument, "No keystores to backup")
	}

	buf := new(bytes.Buffer)
	writer := zip.NewWriter(buf)
	for i, k := range keystoresToBackup {
		encodedFile, err := json.MarshalIndent(k, "", "\t")
		if err != nil {
			if err := writer.Close(); err != nil {
				log.WithError(err).Error("Could not close zip file after writing")
			}
			return nil, status.Errorf(codes.Internal, "could not marshal keystore to JSON file: %v", err)
		}
		f, err := writer.Create(fmt.Sprintf("keystore-%d.json", i))
		if err != nil {
			if err := writer.Close(); err != nil {
				log.WithError(err).Error("Could not close zip file after writing")
			}
			return nil, status.Errorf(codes.Internal, "Could not write keystore file to zip: %v", err)
		}
		if _, err = f.Write(encodedFile); err != nil {
			if err := writer.Close(); err != nil {
				log.WithError(err).Error("Could not close zip file after writing")
			}
			return nil, status.Errorf(codes.Internal, "Could not write keystore file contents")
		}
	}
	if err := writer.Close(); err != nil {
		log.WithError(err).Error("Could not close zip file after writing")
	}
	return &pb.BackupAccountsResponse{
		ZipFile: buf.Bytes(),
	}, nil
}

// VoluntaryExit performs a voluntary exit for the validator keys specified in a request.
// DEPRECATE: Prysm Web UI and associated endpoints will be fully removed in a future hard fork. There is a similar endpoint that is still used /eth/v1alpha1/validator/exit.
func (s *Server) VoluntaryExit(
	ctx context.Context, req *pb.VoluntaryExitRequest,
) (*pb.VoluntaryExitResponse, error) {
	if s.validatorService == nil {
		return nil, status.Error(codes.FailedPrecondition, "Validator service not yet initialized")
	}
	if len(req.PublicKeys) == 0 {
		return nil, status.Error(codes.InvalidArgument, "No public keys specified to delete")
	}
	if s.wallet == nil {
		return nil, status.Error(codes.FailedPrecondition, "No wallet found")
	}
	km, err := s.validatorService.Keymanager()
	if err != nil {
		return nil, err
	}
	formattedKeys := make([]string, len(req.PublicKeys))
	for i, key := range req.PublicKeys {
		formattedKeys[i] = fmt.Sprintf("%#x", key)
	}
	cfg := accounts.PerformExitCfg{
		ValidatorClient:  s.beaconNodeValidatorClient,
		NodeClient:       s.beaconNodeClient,
		Keymanager:       km,
		RawPubKeys:       req.PublicKeys,
		FormattedPubKeys: formattedKeys,
	}
	rawExitedKeys, _, err := accounts.PerformVoluntaryExit(ctx, cfg)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Could not perform voluntary exit: %v", err)
	}
	return &pb.VoluntaryExitResponse{
		ExitedKeys: rawExitedKeys,
	}, nil
}

/**
* PVER
 */

// CreateAccountsAndDepositData initialize validator accounts with deposit data
func (s *Server) CreateAccountsAndDepositData(
	ctx context.Context, req *pb.CreateAccountsRequest,
) (*pb.ListDepositDataResponse, error) {
	if s.validatorService == nil {
		return nil, status.Error(codes.FailedPrecondition, "Validator service not yet initialized")
	}
	if s.wallet == nil {
		return nil, status.Error(codes.NotFound, "Wallet is Not Opened")
	}
	km, err := s.validatorService.Keymanager()
	if err != nil {
		return nil, err
	}

	decryptedPassword, err := aes.Decrypt(s.cipherKey, []byte(req.Password))
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "Could not decrypt password")
	}
	latestIndex, err := createAccountsFromDerivedWallet(ctx, s.wallet, string(decryptedPassword), req.NumAccounts)
	if err != nil {
		return nil, status.Error(codes.Internal, "Could not recover accounts from wallet")
	}
	keys, err := km.FetchValidatingPublicKeys(ctx)
	if err != nil {
		return nil, err
	}
	depositDataList := make([]*pb.DepositDataResponse, int(req.NumAccounts))

	numAccounts := int(req.NumAccounts)

	for i := 0; i < numAccounts; i++ {
		keyIndex := i + latestIndex - numAccounts
		key, err := bls.PublicKeyFromBytes(keys[keyIndex][:])
		if err != nil {
			return nil, status.Error(codes.Internal, "Could not derive public key from bytes")
		}
		dd, err := createDepositData(ctx, key, req.WithdrawKey, req.AmountGwei, km.Sign)
		if err != nil {
			return nil, status.Error(codes.Internal, "Could not create deposit data")
		}
		depositDataList[i] = dd
	}

	return &pb.ListDepositDataResponse{
		DepositDataList: depositDataList,
	}, nil
}

// CreateDepositDataList creates DepositData list with given request.
// NOTE: Validator client does not store these values.
func (s *Server) CreateDepositDataList(ctx context.Context, req *pb.ListDepositDataRequest) (*pb.ListDepositDataResponse, error) {
	if s.validatorService == nil {
		return nil, status.Error(codes.NotFound, "Validator Service is Not Opened")
	}
	if s.wallet == nil {
		return nil, status.Error(codes.NotFound, "Wallet is Not Opened")
	}
	if len(req.DepositDataInputs) == 0 {
		return nil, status.Error(codes.InvalidArgument, "Deposit Data Keys is Empty")
	}
	km, err := s.validatorService.Keymanager()
	if err != nil {
		return nil, status.Error(codes.Internal, "Could not get keymanager")
	}

	datas := make([]*pb.DepositDataResponse, len(req.DepositDataInputs))

	for i, key := range req.DepositDataInputs {
		pubKey, err := bls.PublicKeyFromBytes(key.Pubkey)
		if err != nil {
			return nil, status.Error(codes.Internal, "Could not parse public key")
		}
		dd, err := createDepositData(ctx, pubKey, key.WithdrawKey, key.AmountGwei, km.Sign)
		if err != nil {
			return nil, status.Error(codes.Internal, "Could not create deposit data")
		}
		datas[i] = dd
	}

	return &pb.ListDepositDataResponse{
		DepositDataList: datas,
	}, nil
}

func createAccountsFromDerivedWallet(
	ctx context.Context,
	w *wallet.Wallet,
	mnemonicPassphrase string,
	numAccounts uint64,
) (int, error) {
	km, err := derived.NewKeymanager(ctx, &derived.SetupConfig{
		Wallet:           w,
		ListenForChanges: true,
	})
	if err != nil {
		return 0, errors.Wrap(err, "could not make keymanager for given phrase")
	}
	// TODO: Use encrypted mnemonic passphrase [@gazzua]
	mnemonicStore, err := derived.LoadMnemonic(w.AccountsDir(), mnemonicPassphrase)
	if err != nil {
		return 0, errors.Wrap(err, "could not load mnemonic")
	}

	mnemonicLanguage := "english"
	latestIndex := mnemonicStore.LatestIndex
	newIndex := int(latestIndex + numAccounts)

	err = km.RecoverAccountsFromMnemonic(ctx, mnemonicStore.Mnemonic, mnemonicLanguage, mnemonicPassphrase, newIndex)
	if err != nil {
		return 0, err
	}

	err = derived.SaveMnemonicStore(mnemonicStore.Mnemonic, mnemonicPassphrase, w.AccountsDir(), uint64(newIndex))
	if err != nil {
		return 0, err
	}

	return newIndex, nil
}

func createDepositData(
	ctx context.Context,
	depositPubkey bls.PublicKey,
	eth1WithdrawalAddress []byte,
	amountInGwei uint64,
	signer iface.SigningFunc,
) (*pb.DepositDataResponse, error) {
	depositMessage := &ethpb.DepositMessage{
		PublicKey:             depositPubkey.Marshal(),
		WithdrawalCredentials: withdrawalCredentialsHash(eth1WithdrawalAddress),
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
	dd := &pb.DepositDataResponse{
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

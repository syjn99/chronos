package rpc

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/pkg/errors"
	"github.com/prysmaticlabs/prysm/v4/api/pagination"
	"github.com/prysmaticlabs/prysm/v4/beacon-chain/core/signing"
	"github.com/prysmaticlabs/prysm/v4/cmd"
	"github.com/prysmaticlabs/prysm/v4/config/params"
	"github.com/prysmaticlabs/prysm/v4/crypto/aes"
	"github.com/prysmaticlabs/prysm/v4/crypto/bls"
	ethpb "github.com/prysmaticlabs/prysm/v4/proto/prysm/v1alpha1"
	pb "github.com/prysmaticlabs/prysm/v4/proto/prysm/v1alpha1/validator-client"
	"github.com/prysmaticlabs/prysm/v4/validator/accounts/petnames"
	"github.com/prysmaticlabs/prysm/v4/validator/accounts/wallet"
	iface "github.com/prysmaticlabs/prysm/v4/validator/client/iface"
	"github.com/prysmaticlabs/prysm/v4/validator/keymanager"
	"github.com/prysmaticlabs/prysm/v4/validator/keymanager/derived"
	keystorev4 "github.com/wealdtech/go-eth2-wallet-encryptor-keystorev4"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var tempPassword = "temp-password"

// ListAccounts allows retrieval of validating keys and their petnames
// for a user's wallet via RPC.
// DEPRECATED: Prysm Web UI and associated endpoints will be fully removed in a future hard fork.
// USED IN OVER-NODE
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

/**
* OVER-NODE
 */

// CreateAccountsAndDepositData initialize validator accounts with deposit data
func (s *Server) CreateAccountsAndDepositData(
	ctx context.Context, req *pb.CreateAccountsRequest,
) (*pb.ListDepositDataResponse, error) {
	if !s.isOverNode {
		log.Debug("CreateAccountsAndDepositData was called when over node flag disabled")
		return nil, status.Error(codes.NotFound, "Only available in over node flag enabled")
	}
	if s.validatorService == nil {
		log.Debug("CreateAccountsAndDepositData was called when validator service is not opened")
		return nil, status.Error(codes.Unavailable, "Validator service not yet initialized")
	}
	if s.wallet == nil {
		log.Debug("CreateAccountsAndDepositData was called when wallet is not opened")
		return nil, status.Error(codes.Unavailable, "Wallet is Not Opened")
	}

	km, err := s.validatorService.Keymanager()
	if err != nil {
		log.WithError(err).Error("Could not get keymanager")
		return nil, status.Error(codes.Internal, "Could not get keymanager")
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
	latestIndex, err := createAccountsFromDerivedWallet(ctx, s.wallet, string(decryptedPassword), req.NumAccounts)
	if err != nil {
		log.WithError(err).Error("Could not recover accounts from wallet")
		return nil, status.Error(codes.Internal, "Could not recover accounts from wallet")
	}
	keys, err := km.FetchValidatingPublicKeys(ctx)
	if err != nil {
		log.WithError(err).Error("Could not fetch validating public keys")
		return nil, status.Error(codes.Internal, "Could not fetch validating public keys")
	}
	depositDataList := make([]*pb.DepositDataResponse, int(req.NumAccounts))

	numAccounts := int(req.NumAccounts)

	for i := 0; i < numAccounts; i++ {
		keyIndex := i + latestIndex - numAccounts
		key, err := bls.PublicKeyFromBytes(keys[keyIndex][:])
		if err != nil {
			log.WithError(err).Error("Could not derive public key from bytes")
			return nil, status.Error(codes.Internal, "Could not derive public key from bytes")
		}
		dd, err := createDepositData(ctx, key, req.WithdrawKey, req.AmountGwei, km.Sign)
		if err != nil {
			log.WithError(err).Error("Could not create deposit data")
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
	if !s.isOverNode {
		log.Debug("CreateDepositDataList was called when over node flag disabled")
		return nil, status.Error(codes.NotFound, "Only available in over node flag enabled")
	}
	if s.validatorService == nil {
		log.Debug("CreateDepositDataList was called when validator service is not opened")
		return nil, status.Error(codes.Unavailable, "Validator Service is Not Opened")
	}
	if s.wallet == nil {
		log.Debug("CreateDepositDataList was called when wallet is not opened")
		return nil, status.Error(codes.Unavailable, "Wallet is Not Opened")
	}
	if len(req.DepositDataInputs) == 0 {
		log.Debug("CreateDepositDataList was called with empty Deposit Data")
		return nil, status.Error(codes.InvalidArgument, "Deposit Data Keys is Empty")
	}
	km, err := s.validatorService.Keymanager()
	if err != nil {
		log.WithError(err).Error("Could not get keymanager")
		return nil, status.Error(codes.Internal, "Could not get keymanager")
	}

	datas := make([]*pb.DepositDataResponse, len(req.DepositDataInputs))

	for i, key := range req.DepositDataInputs {
		pubKey, err := bls.PublicKeyFromBytes(key.Pubkey)
		if err != nil {
			log.WithError(err).Error("Could not parse public key")
			return nil, status.Error(codes.Internal, "Could not parse public key")
		}
		dd, err := createDepositData(ctx, pubKey, key.WithdrawKey, key.AmountGwei, km.Sign)
		if err != nil {
			log.WithError(err).Error("Could not create deposit data")
			return nil, status.Error(codes.Internal, "Could not create deposit data")
		}
		datas[i] = dd
	}

	return &pb.ListDepositDataResponse{
		DepositDataList: datas,
	}, nil
}

func (s *Server) ImportAccounts(ctx context.Context, req *pb.ImportAccountsRequest) (*pb.ImportAccountsResponse, error) {
	s.rpcMutex.Lock()
	defer s.rpcMutex.Unlock()
	if !s.isOverNode {
		log.Debug("ImportAccounts was called when over node flag disabled")
		return nil, status.Error(codes.NotFound, "Only available in over node flag enabled")
	}
	if s.validatorService == nil {
		log.Debug("ImportAccounts was called when validator service is not opened")
		return nil, status.Error(codes.Unavailable, "Validator Service is Not Opened")
	}
	if s.wallet == nil {
		log.Debug("ImportAccounts was called when wallet is not opened")
		return nil, status.Error(codes.Unavailable, "Wallet is Not Opened")
	}
	if !s.walletInitialized {
		return nil, status.Error(codes.FailedPrecondition, "Wallet not yet initialized")
	}
	if len(req.PrivateKeys) == 0 {
		log.Debug("ImportAccounts was called with empty Deposit Data")
		return nil, status.Error(codes.InvalidArgument, "Deposit Data Keys is Empty")
	}
	km, err := s.validatorService.Keymanager()
	if err != nil {
		log.WithError(err).Error("Could not get keymanager")
		return nil, status.Error(codes.Internal, "Could not get keymanager")
	}

	importer, ok := km.(keymanager.Importer)
	if !ok {
		log.WithError(err).Error("Keymanager cannot import local keys")
		return nil, status.Error(codes.Internal, "Keymanager cannot import local keys")
	}

	// Decrypt PrivateKeys with cipherKey
	secretKeys := make([]bls.SecretKey, 0)
	for _, encodedPrivateKey := range req.PrivateKeys {
		privateKey, err := hexutil.Decode(encodedPrivateKey)
		if err != nil {
			log.Debug("Could not decrypt private key", err)
			return nil, status.Error(codes.InvalidArgument, "Could not decrypt private key")
		}
		decryptedPrivateKey, err := aes.Decrypt(s.cipherKey, privateKey)
		if err != nil {
			log.Debug("Could not decrypt private key", err)
			return nil, status.Error(codes.InvalidArgument, "Could not decrypt private key")
		}

		secretKey, err := bls.SecretKeyFromBytes(decryptedPrivateKey)
		if err != nil {
			log.Debug("Could not decrypt private key", err)
			return nil, status.Error(codes.InvalidArgument, "Could not decrypt private key")
		}
		secretKeys = append(secretKeys, secretKey)
	}

	// Make keystore for each secret key
	keystores := make([]*keymanager.Keystore, len(secretKeys))
	passwords := make([]string, len(secretKeys))
	encryptor := keystorev4.New()
	for i := 0; i < len(secretKeys); i++ {
		pubkey := secretKeys[i].PublicKey()
		cryptoFields, err := encryptor.Encrypt(secretKeys[i].Marshal(), tempPassword)
		if err != nil {
			log.WithError(err).Error("Could not encrypt secret key when wrapping it in a keystore")
			return nil, status.Error(codes.Internal, "Could not encrypt secret key when wrapping it in a keystore")
		}
		k := &keymanager.Keystore{
			Crypto:      cryptoFields,
			ID:          fmt.Sprint(i), // temporary keystore ID
			Pubkey:      fmt.Sprintf("%x", pubkey.Marshal()),
			Version:     encryptor.Version(),
			Description: encryptor.Name(),
		}
		keystores[i] = k
		passwords[i] = tempPassword
	}

	statuses, err := importer.ImportKeystores(ctx, keystores, passwords)
	if err != nil {
		log.WithError(err).Error("Could not import keys", err)
		return nil, status.Error(codes.Internal, "Could not import keys")
	}

	// Map the statuses to the response
	importedStatuses := make([]*pb.ImportKeystoreStatus, len(statuses))
	for i, status := range statuses {
		importedStatuses[i] = &pb.ImportKeystoreStatus{
			PublicKey: "0x" + keystores[i].Pubkey,
			Status:    status,
		}
	}

	return &pb.ImportAccountsResponse{
		Data: importedStatuses,
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

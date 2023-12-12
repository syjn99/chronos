package derived

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/prysmaticlabs/prysm/v4/crypto/rand"
	"github.com/prysmaticlabs/prysm/v4/io/file"
	"github.com/prysmaticlabs/prysm/v4/io/prompt"
	"github.com/prysmaticlabs/prysm/v4/validator/keymanager"
	"github.com/tyler-smith/go-bip39"
	"github.com/tyler-smith/go-bip39/wordlists"
	keystorev4 "github.com/wealdtech/go-eth2-wallet-encryptor-keystorev4"
)

const (
	confirmationText      = "Confirm you have written down the recovery words somewhere safe (offline) [y|Y]"
	MnemonicStoreFileName = "mnemonic-store.json"
)

// MnemonicGenerator implements methods for creating
// mnemonic seed phrases in english using a given
// source of entropy such as a private key.
type MnemonicGenerator struct {
	skipMnemonicConfirm bool
}

type MnemonicStore struct {
	Mnemonic    string
	LatestIndex uint64
}

type MnemonicStoreRepresentation struct {
	Crypto  map[string]interface{} `json:"crypto"`
	ID      string                 `json:"uuid"`
	Version uint                   `json:"version"`
	Name    string                 `json:"name"`
}

// ErrUnsupportedMnemonicLanguage is returned when trying to use an unsupported mnemonic language.
var (
	mutex                          sync.Mutex
	DefaultMnemonicLanguage        = "english"
	ErrUnsupportedMnemonicLanguage = errors.New("unsupported mnemonic language")
)

// GenerateAndConfirmMnemonic requires confirming the generated mnemonics.
func GenerateAndConfirmMnemonic(mnemonicLanguage string, skipMnemonicConfirm bool) (string, error) {
	mnemonicRandomness := make([]byte, 32)
	if _, err := rand.NewGenerator().Read(mnemonicRandomness); err != nil {
		return "", errors.Wrap(err, "could not initialize mnemonic source of randomness")
	}
	err := setBip39Lang(mnemonicLanguage)
	if err != nil {
		return "", err
	}
	m := &MnemonicGenerator{
		skipMnemonicConfirm: skipMnemonicConfirm,
	}
	phrase, err := m.Generate(mnemonicRandomness)
	if err != nil {
		return "", errors.Wrap(err, "could not generate wallet seed")
	}
	if err := m.ConfirmAcknowledgement(phrase); err != nil {
		return "", errors.Wrap(err, "could not confirm mnemonic acknowledgement")
	}
	return phrase, nil
}

// GenerateAndSaveMnemonic generates a mnemonic and saves it to a file.
func GenerateAndSaveMnemonic(mnemonicLanguage string, password string, path string) error {
	mnemonicRandomness := make([]byte, 32)
	if _, err := rand.NewGenerator().Read(mnemonicRandomness); err != nil {
		return errors.Wrap(err, "could not initialize mnemonic source of randomness")
	}
	err := setBip39Lang(mnemonicLanguage)
	if err != nil {
		return err
	}
	m := &MnemonicGenerator{
		skipMnemonicConfirm: true,
	}
	phrase, err := m.Generate(mnemonicRandomness)
	if err != nil {
		return errors.Wrap(err, "could not generate wallet seed")
	}
	return SaveMnemonicStore(phrase, password, path, 0)
}

// Generate a mnemonic seed phrase in english using a source of
// entropy given as raw bytes.
func (_ *MnemonicGenerator) Generate(data []byte) (string, error) {
	return bip39.NewMnemonic(data)
}

// ConfirmAcknowledgement displays the mnemonic phrase to the user
// and confirms the user has written down the phrase securely offline.
func (m *MnemonicGenerator) ConfirmAcknowledgement(phrase string) error {
	log.Info(
		"Write down the sentence below, as it is your only " +
			"means of recovering your wallet",
	)
	fmt.Printf(
		`=================Wallet Seed Recovery Phrase====================

%s

===================================================================`,
		phrase)
	fmt.Println("")
	if m.skipMnemonicConfirm {
		return nil
	}
	// Confirm the user has written down the mnemonic phrase offline.
	_, err := prompt.ValidatePrompt(os.Stdin, confirmationText, prompt.ValidateConfirmation)
	if err != nil {
		log.Errorf("Could not confirm acknowledgement of userprompt, please enter y")
	}
	return nil
}

// Uses the provided mnemonic seed phrase to generate the
// appropriate seed file for recovering a derived wallets.
func seedFromMnemonic(mnemonic, mnemonicLanguage, mnemonicPassphrase string) ([]byte, error) {
	err := setBip39Lang(mnemonicLanguage)
	if err != nil {
		return nil, err
	}
	if ok := bip39.IsMnemonicValid(mnemonic); !ok {
		return nil, bip39.ErrInvalidMnemonic
	}
	return bip39.NewSeed(mnemonic, mnemonicPassphrase), nil
}

func setBip39Lang(lang string) error {
	// mutex is used to prevent concurrent access to bip39.SetWordList
	mutex.Lock()
	defer mutex.Unlock()
	var wordlist []string
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

	if wl, ok := allowedLanguages[lang]; ok {
		wordlist = wl
	} else {
		return errors.Wrapf(ErrUnsupportedMnemonicLanguage, "%s", lang)
	}
	bip39.SetWordList(wordlist)
	return nil
}

func SaveMnemonicStore(mnemonic string, password string, path string, lastIndex uint64) error {
	encryptor := keystorev4.New()
	id, err := uuid.NewRandom()
	if err != nil {
		return err
	}

	mnemonicStore := &MnemonicStore{
		Mnemonic:    mnemonic,
		LatestIndex: lastIndex,
	}

	encodedStore, err := json.MarshalIndent(mnemonicStore, "", "\t")
	if err != nil {
		return err
	}

	cryptoFields, err := encryptor.Encrypt([]byte(encodedStore), password)
	if err != nil {
		return err
	}

	mnemonicStoreRepresentation := &MnemonicStoreRepresentation{
		ID:      id.String(),
		Crypto:  cryptoFields,
		Version: encryptor.Version(),
		Name:    encryptor.Name(),
	}
	encoded, err := json.MarshalIndent(mnemonicStoreRepresentation, "", "\t")
	if err != nil {
		return err
	}
	err = file.WriteFile(filepath.Join(filepath.Clean(path), MnemonicStoreFileName), encoded)
	if err != nil {
		return err
	}
	return nil
}

func LoadMnemonic(path string, password string) (*MnemonicStore, error) {
	mnemonicFile := filepath.Join(path, MnemonicStoreFileName)
	if file.FileExists(mnemonicFile) {
		rawData, err := os.ReadFile(filepath.Clean(mnemonicFile))
		if err != nil {
			return nil, errors.Wrap(err, "could not read mnemonic store file")
		}
		keystoreFile := &MnemonicStoreRepresentation{}
		if err := json.Unmarshal(rawData, keystoreFile); err != nil {
			return nil, errors.Wrap(err, "could not unmarshal mnemonic store file")
		}

		decryptor := keystorev4.New()
		enc, err := decryptor.Decrypt(keystoreFile.Crypto, password)
		if err != nil && strings.Contains(err.Error(), keymanager.IncorrectPasswordErrMsg) {
			return nil, errors.Wrap(err, "wrong password for wallet entered")
		} else if err != nil {
			return nil, errors.Wrap(err, "could not decrypt keystore")
		}

		mnemonicStore := &MnemonicStore{}
		if err := json.Unmarshal(enc, mnemonicStore); err != nil {
			return nil, errors.Wrap(err, "could not unmarshal mnemonic store")
		}

		return mnemonicStore, nil
	} else {
		return nil, errors.New("mnemonic file does not exist")
	}
}

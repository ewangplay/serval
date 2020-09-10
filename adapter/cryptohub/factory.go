package cryptohub

import (
	"errors"
	"strings"
	"sync"
)

const (
	// ED25519 signatures are elliptic-curve signatures,
	// carefully engineered at several levels of design
	// and implementation to achieve very high speeds
	// without compromising security.
	ED25519 = "ED25519"
)

// Global instance definition
var (
	initOnce sync.Once
	gSigner  Signer
)

// SignerConfig represents the signer configure
type SignerConfig struct {
	Algorithm string
}

// Config represetns the crypto hub configure
type Config struct {
	Signer *SignerConfig
}

// GetDefaultConfig returns a default configure for crypto hub
func GetDefaultConfig() *Config {
	return &Config{
		Signer: &SignerConfig{
			Algorithm: ED25519,
		},
	}
}

// InitCryptoHub initializes the cryptohub instance with singleton mode
// This method MUST be called before any other method can be called.
func InitCryptoHub(cfg *Config) error {
	var err error

	initOnce.Do(func() {
		err = initCryptHub(cfg)
	})

	return err
}

func initCryptHub(cfg *Config) error {
	if cfg == nil {
		cfg = GetDefaultConfig()
	}
	if cfg.Signer == nil {
		cfg.Signer = GetDefaultConfig().Signer
	}
	if cfg.Signer.Algorithm == "" {
		cfg.Signer.Algorithm = ED25519
	}

	var err error
	err = initSigner(cfg.Signer)
	if err != nil {
		return err
	}

	return nil
}

func initSigner(cfg *SignerConfig) error {
	var err error

	switch strings.ToUpper(cfg.Algorithm) {
	case ED25519:
		gSigner, err = CreateEd25519Signer()
	default:
		err = errors.New("unsupported signature algorithm")
	}

	return err
}

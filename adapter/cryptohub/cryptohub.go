package cryptohub

import (
	"sync"

	ch "github.com/ewangplay/cryptohub"
)

// Global instance definition
var (
	initOnce sync.Once
	gCSP     ch.CSP
)

// InitCryptoHub initializes the cryptohub instance with singleton mode
func InitCryptoHub() error {
	var err error

	initOnce.Do(func() {
		err = initCryptHub()
	})

	return err
}

func initCryptHub() (err error) {
	cfg := &ch.Config{
		ProviderName: "SW",
	}
	gCSP, err = ch.GetCSP(cfg)
	if err != nil {
		return err
	}
	return nil
}

// GenEd25519Key generates an Ed25519 key.
func GenEd25519Key() (k ch.Key, err error) {
	assertCSPValid()
	return gCSP.KeyGen(&ch.ED25519KeyGenOpts{})
}

// Sign signs digest using key k.
func Sign(k ch.Key, digest []byte) (signature []byte, err error) {
	assertCSPValid()
	return gCSP.Sign(k, digest)
}

// Verify verifies signature against key k and digest
func Verify(k ch.Key, digest, signature []byte) (valid bool, err error) {
	assertCSPValid()
	return gCSP.Verify(k, digest, signature)
}

func assertCSPValid() {
	if gCSP == nil {
		panic("CSP not be initialized")
	}
}

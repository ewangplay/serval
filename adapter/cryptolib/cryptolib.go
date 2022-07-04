package cryptolib

import (
	"sync"

	cl "github.com/ewangplay/cryptolib"
)

// Global instance definition
var (
	initOnce sync.Once
	gCSP     cl.CSP
)

// Initcryptolib initializes the cryptolib instance with singleton mode
func Initcryptolib() error {
	var err error

	initOnce.Do(func() {
		err = initCryptHub()
	})

	return err
}

func initCryptHub() (err error) {
	cfg := &cl.Config{
		ProviderName: "SW",
	}
	gCSP, err = cl.GetCSP(cfg)
	if err != nil {
		return err
	}
	return nil
}

// GenEd25519Key generates an Ed25519 key.
func GenEd25519Key() (k cl.Key, err error) {
	assertCSPValid()
	return gCSP.KeyGen(&cl.ED25519KeyGenOpts{})
}

// Sign signs digest using key k.
func Sign(k cl.Key, digest []byte) (signature []byte, err error) {
	assertCSPValid()
	return gCSP.Sign(k, digest, nil)
}

// Verify verifies signature against key k and digest
func Verify(k cl.Key, digest, signature []byte) (valid bool, err error) {
	assertCSPValid()
	return gCSP.Verify(k, digest, signature, nil)
}

func assertCSPValid() {
	if gCSP == nil {
		panic("CSP not be initialized")
	}
}

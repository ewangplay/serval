package cryptolib

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"hash"
	"sync"

	cl "github.com/ewangplay/cryptolib"
	"github.com/ewangplay/serval/io"
	"github.com/jerray/qsign"
)

// Global instance definition
var (
	initOnce sync.Once
	gCSP     cl.CSP
	gQsign   *qsign.Qsign
)

// InitCryptolib initializes the cryptolib instance with singleton mode
func InitCryptolib() error {
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

	// New Qsign instance
	gQsign = qsign.NewQsign(qsign.Options{
		// To use a hash.Hash other than md5
		Hasher: func() hash.Hash {
			return sha256.New()
		},
		// To use a encoding other than hex
		Encoder: func() qsign.Encoding {
			return base64.StdEncoding
		},
	})

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

func VerifyDDO(ddo *io.DDO) (valid bool, err error) {
	if ddo == nil {
		return false, fmt.Errorf("did document is nil")
	}
	if ddo.Proof.Type == "" || ddo.Proof.Creator == "" || ddo.Proof.SignatureValue == "" {
		return false, fmt.Errorf("did document proof invalid")
	}
	if len(ddo.PublicKey) == 0 {
		return false, fmt.Errorf("did document public key list invalid")
	}
	var pk io.PublicKey
	hasPubKey := false
	for _, pk = range ddo.PublicKey {
		if pk.ID == ddo.Proof.Creator && pk.Type == ddo.Proof.Type {
			hasPubKey = true
			break
		}
	}
	if !hasPubKey {
		return false, fmt.Errorf("did document proof public key missing")
	}

	pubKeyBytes, err := hex.DecodeString(pk.PublicKeyHex)
	if err != nil {
		return false, fmt.Errorf("Decode public key for %s failed: %v", pk.ID, err)
	}

	var k cl.Key
	switch pk.Type {
	case cl.ED25519:
		k = &cl.Ed25519PublicKey{
			PubKey: pubKeyBytes,
		}
	default:
		return false, fmt.Errorf("unsupported key type for did: %v", pk.Type)
	}

	signature, err := base64.StdEncoding.DecodeString(ddo.Proof.SignatureValue)
	if err != nil {
		return false, fmt.Errorf("did document proof signature invalid")
	}

	data, err := gQsign.Digest(ddo)
	if err != nil {
		return
	}
	digest, err := gCSP.Hash(data, &cl.SHA256Opts{})
	if err != nil {
		return
	}
	return gCSP.Verify(k, digest, signature, nil)
}

func assertCSPValid() {
	if gCSP == nil {
		panic("CSP not be initialized")
	}
}

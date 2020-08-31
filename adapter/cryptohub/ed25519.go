package cryptohub

import (
	"crypto/ed25519"
	"fmt"
)

// Ed25519PublicKey represents the ed25519 public key
type Ed25519PublicKey []byte

// GetPublicKey returns the public key bytes
func (edpubk Ed25519PublicKey) GetPublicKey() []byte {
	return edpubk
}

// Ed25519PrivateKey represents the ed25519 private key
type Ed25519PrivateKey []byte

// GetPrivateKey returns the private key bytes
func (edprik Ed25519PrivateKey) GetPrivateKey() []byte {
	return edprik
}

// Ed25519CryptoHub represents the ed25519 crypto hub
type Ed25519CryptoHub struct {
}

// CreateEd25519CryptoHub creates an instance of ed25519 crypto hub
func CreateEd25519CryptoHub() *Ed25519CryptoHub {
	return &Ed25519CryptoHub{}
}

// GenKey generates one ed25519 public key and private key pair
func (ed *Ed25519CryptoHub) GenKey() (PublicKey, PrivateKey, error) {
	pubKey, priKey, err := ed25519.GenerateKey(nil)
	if err != nil {
		return nil, nil, err
	}
	return Ed25519PublicKey(pubKey), Ed25519PrivateKey(priKey), nil
}

// Sign signs the message with privateKey and returns signature
func (ed *Ed25519CryptoHub) Sign(privateKey PrivateKey, message []byte) (sigature []byte, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = fmt.Errorf("sign error: %v", e)
		}
	}()
	return ed25519.Sign(ed25519.PrivateKey(privateKey.GetPrivateKey()), message), nil
}

// Verify verifies signature against publicKey and message
func (ed *Ed25519CryptoHub) Verify(publicKey PublicKey, message []byte, signature []byte) (valid bool, err error) {
	valid = ed25519.Verify(ed25519.PublicKey(publicKey.GetPublicKey()), message, signature)
	return valid, nil
}

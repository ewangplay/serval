package cryptohub

import (
	"crypto/ed25519"
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

// KeyPair generates one ed25519 public key and private key pair
func (ed *Ed25519CryptoHub) KeyPair() (PublicKey, PrivateKey, error) {
	pubKey, priKey, err := ed25519.GenerateKey(nil)
	if err != nil {
		return nil, nil, err
	}
	return Ed25519PublicKey(pubKey), Ed25519PrivateKey(priKey), nil
}

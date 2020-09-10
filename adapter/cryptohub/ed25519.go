package cryptohub

import (
	"crypto/ed25519"
	"fmt"
)

// Ed25519PublicKey represents the ed25519 public key
type Ed25519PublicKey []byte

// Bytes converts this key to its byte representation.
func (edpubk Ed25519PublicKey) Bytes() ([]byte, error) {
	return edpubk, nil
}

// Symmetric returns true if this key is a symmetric key,
// false is this key is asymmetric
func (edpubk Ed25519PublicKey) Symmetric() bool {
	return false
}

// Private returns true if this key is a private key,
// false otherwise.
func (edpubk Ed25519PublicKey) Private() bool {
	return false
}

// PublicKey returns the corresponding public key part of an asymmetric public/private key pair.
// This method returns an error in symmetric key schemes.
func (edpubk Ed25519PublicKey) PublicKey() (Key, error) {
	return nil, fmt.Errorf("it is already a public key")
}

// Ed25519PrivateKey represents the ed25519 private key
type Ed25519PrivateKey struct {
	PrivateKeyBytes []byte
	PublicKeyBytes  []byte
}

// Bytes converts this key to its byte representation.
func (edprik Ed25519PrivateKey) Bytes() ([]byte, error) {
	return edprik.PrivateKeyBytes, nil
}

// Symmetric returns true if this key is a symmetric key,
// false is this key is asymmetric
func (edprik Ed25519PrivateKey) Symmetric() bool {
	return false
}

// Private returns true if this key is a private key,
// false otherwise.
func (edprik Ed25519PrivateKey) Private() bool {
	return true
}

// PublicKey returns the corresponding public key part of an asymmetric public/private key pair.
// This method returns an error in symmetric key schemes.
func (edprik Ed25519PrivateKey) PublicKey() (Key, error) {
	return Ed25519PublicKey(edprik.PublicKeyBytes), nil
}

// Ed25519Signer represents the ed25519 crypto hub
type Ed25519Signer struct {
}

// CreateEd25519Signer creates an instance of ed25519 crypto hub
func CreateEd25519Signer() (*Ed25519Signer, error) {
	return &Ed25519Signer{}, nil
}

// GenKey generates a private key of ed25519 algorithm
func (ed *Ed25519Signer) GenKey() (k Key, err error) {
	pubKey, priKey, err := ed25519.GenerateKey(nil)
	if err != nil {
		return nil, err
	}
	return &Ed25519PrivateKey{
		PrivateKeyBytes: priKey,
		PublicKeyBytes:  pubKey,
	}, nil
}

// Sign signs digest using key k
func (ed *Ed25519Signer) Sign(k Key, digest []byte) (signature []byte, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = fmt.Errorf("sign error: %v", e)
		}
	}()

	priKeyBytes, err := k.Bytes()
	if err != nil {
		return nil, err
	}
	return ed25519.Sign(ed25519.PrivateKey(priKeyBytes), digest), nil
}

// Verify verifies signature against key k and digest
func (ed *Ed25519Signer) Verify(k Key, digest, signature []byte) (valid bool, err error) {
	pubKeyBytes, err := k.Bytes()
	if err != nil {
		return false, err
	}
	valid = ed25519.Verify(ed25519.PublicKey(pubKeyBytes), digest, signature)
	return valid, nil
}

package cryptohub

// PublicKey represents the public key interface
type PublicKey interface {
	GetPublicKey() []byte
}

// PrivateKey represent the private key interface
type PrivateKey interface {
	GetPrivateKey() []byte
}

// CryptoHub represents the crypto hub interface
type CryptoHub interface {
	GenKey() (publicKey PublicKey, privateKey PrivateKey, err error)
	Sign(privateKey PrivateKey, message []byte) (signature []byte, err error)
	Verify(publicKey PublicKey, message []byte, signature []byte) (valid bool, err error)
}

// Global crypto hub instance
var gCryptoHub CryptoHub

// GetCryptoHub returns the crypto hub instance in singleton mode
func GetCryptoHub() CryptoHub {
	if gCryptoHub == nil {
		gCryptoHub = CreateEd25519CryptoHub()
	}
	return gCryptoHub
}

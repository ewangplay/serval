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
var gCH CryptoHub

// InitCryptoHub initializes the cryptohub instance with singleton mode
// This method MUST be called before any other method can be called.
func InitCryptoHub() error {
	if gCH == nil {
		gCH = CreateEd25519CryptoHub()
	}

	return nil
}

// GenKey returns a public and private key pair
func GenKey() (publicKey PublicKey, privateKey PrivateKey, err error) {
	checkInitState()
	return gCH.GenKey()
}

// Sign signs the message with privateKey and returns signature
func Sign(privateKey PrivateKey, message []byte) (signature []byte, err error) {
	checkInitState()
	return gCH.Sign(privateKey, message)
}

// Verify verifies signature against publicKey and message
func Verify(publicKey PublicKey, message []byte, signature []byte) (valid bool, err error) {
	checkInitState()
	return gCH.Verify(publicKey, message, signature)
}

func checkInitState() {
	if gCH == nil {
		panic("cryptohub not initialized")
	}
}

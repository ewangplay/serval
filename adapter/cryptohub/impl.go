package cryptohub

// GenKey generates a key.
func GenKey() (k Key, err error) {
	checkGlobalSigner()
	return gSigner.GenKey()
}

// Sign signs digest using key k.
//
// Note that when a signature of a hash of a larger message is needed,
// the caller is responsible for hashing the larger message and passing
// the hash (as digest).
func Sign(k Key, digest []byte) (signature []byte, err error) {
	checkGlobalSigner()
	return gSigner.Sign(k, digest)
}

// Verify verifies signature against key k and digest
func Verify(k Key, digest, signature []byte) (valid bool, err error) {
	checkGlobalSigner()
	return gSigner.Verify(k, digest, signature)
}

func checkGlobalSigner() {
	if gSigner == nil {
		panic("global signer not be initialized")
	}
}

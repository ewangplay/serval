package utils

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io"

	cl "github.com/ewangplay/cryptolib"
	didio "github.com/ewangplay/serval/io"
	"github.com/jerray/qsign"
)

// GenerateUUID returns a UUID as a string based on RFC 4122
func GenerateUUID() string {
	uuid := GenerateBytesUUID()
	return uuidBytesToStr(uuid)
}

// GenerateBytesUUID returns a UUID as []byte based on RFC 4122
func GenerateBytesUUID() []byte {
	uuid := make([]byte, 16)
	_, err := io.ReadFull(rand.Reader, uuid)
	if err != nil {
		panic(err)
	}

	// variant bits; see section 4.1.1
	uuid[8] = uuid[8]&^0xc0 | 0x80

	// version 4 (pseudo-random); see section 4.1.3
	uuid[6] = uuid[6]&^0xf0 | 0x40

	return uuid
}

func uuidBytesToStr(uuid []byte) string {
	return fmt.Sprintf("%x-%x-%x-%x-%x", uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:])
}

// SHA256 returns the SHA256 checksum of the data.
func SHA256(data []byte) string {
	h := sha256.New()
	h.Write(data)
	cs := h.Sum(nil)
	return hex.EncodeToString(cs)
}

func SignDDO(csp cl.CSP, qs *qsign.Qsign, keyID string, key cl.Key, ddo *didio.DDO) (err error) {
	data, err := qs.Digest(ddo)
	if err != nil {
		return
	}
	digest, err := csp.Hash(data, &cl.SHA256Opts{})
	if err != nil {
		return
	}
	signature, err := csp.Sign(key, digest, nil)
	if err != nil {
		return
	}
	// Set the ddo proof
	ddo.Proof = didio.Proof{
		Type:           key.Type(),
		Creator:        keyID,
		SignatureValue: base64.StdEncoding.EncodeToString(signature),
	}

	return nil
}

func VerifyDDO(csp cl.CSP, qs *qsign.Qsign, ddo *didio.DDO) (err error) {
	if csp == nil {
		return fmt.Errorf("CSP provider is nil")
	}
	if qs == nil {
		return fmt.Errorf("Qsign instance is nil")
	}
	if ddo == nil {
		return fmt.Errorf("DID document is nil")
	}
	if ddo.Proof.Type == "" || ddo.Proof.Creator == "" || ddo.Proof.SignatureValue == "" {
		return fmt.Errorf("The proof of the DID document is missing")
	}
	if len(ddo.PublicKey) == 0 {
		return fmt.Errorf("The public key list of the DID document is missing")
	}

	// Retrieve the public key corresponding to the signature
	var pk didio.PublicKey
	hasPubKey := false
	for _, pk = range ddo.PublicKey {
		if pk.ID == ddo.Proof.Creator && pk.Type == ddo.Proof.Type {
			hasPubKey = true
			break
		}
	}
	if !hasPubKey {
		return fmt.Errorf("The public key corresponding to the signature is missing")
	}

	pubKeyBytes, err := hex.DecodeString(pk.PublicKeyHex)
	if err != nil {
		return fmt.Errorf("Decode the public key (%s) failed: %v", pk.ID, err)
	}

	var k cl.Key
	switch pk.Type {
	case cl.ED25519:
		k = &cl.Ed25519PublicKey{
			PubKey: pubKeyBytes,
		}
	default:
		return fmt.Errorf("unsupported key type: %v", pk.Type)
	}

	// Decode the signature of the DID document
	signature, err := base64.StdEncoding.DecodeString(ddo.Proof.SignatureValue)
	if err != nil {
		return fmt.Errorf("The signature of the DID document is invalid")
	}

	// Verifying the signature of the DID document
	data, err := qs.Digest(ddo)
	if err != nil {
		return err
	}
	digest, err := csp.Hash(data, &cl.SHA256Opts{})
	if err != nil {
		return err
	}
	valid, err := csp.Verify(k, digest, signature, nil)
	if err != nil {
		return err
	}
	if !valid {
		return fmt.Errorf("Verifying the signature of the DID document failed")
	}

	return nil
}

func SignProof(csp cl.CSP, did string, k cl.Key) (signature []byte, err error) {
	digest, err := csp.Hash([]byte(did), &cl.SHA256Opts{})
	if err != nil {
		return
	}

	return csp.Sign(k, digest, nil)
}

func VerifyProof(csp cl.CSP, did string, proof *didio.Proof, ddo *didio.DDO) (valid bool, err error) {

	if len(ddo.PublicKey) == 0 {
		return false, fmt.Errorf("did document public key list invalid")
	}

	var pk didio.PublicKey
	hasPubKey := false
	for _, pk = range ddo.PublicKey {
		if pk.ID == proof.Creator && pk.Type == proof.Type {
			hasPubKey = true
			break
		}
	}
	if !hasPubKey {
		return false, fmt.Errorf("did document recovery public key missing")
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

	signature, err := base64.StdEncoding.DecodeString(proof.SignatureValue)
	if err != nil {
		return false, fmt.Errorf("proof signature is invalid")
	}

	digest, err := csp.Hash([]byte(did), &cl.SHA256Opts{})
	if err != nil {
		return
	}
	return csp.Verify(k, digest, signature, nil)
}

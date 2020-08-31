package v1

import (
	"time"
)

// Key represents the key pair
type Key struct {
	ID            string  `json:"id"`
	Type          KeyType `json:"type"`
	PrivateKeyHex string  `json:"privateKeyHex"`
	PublicKeyHex  string  `json:"publicKeyHex"`
}

// CreateDidResponse represents the CreateDid response
type CreateDidResponse struct {
	Did     string    `json:"did"`
	Key     []*Key    `json:"key"`
	Created time.Time `json:"created"`
}

// ResolveDidResponse represents the ResolveDid response
type ResolveDidResponse struct {
	Did      string `json:"did"`
	Document *DDO   `json:"document"`
}

// KeyType represents key type
type KeyType string

// KeyType Definition
const (
	RSAKey        KeyType = "RSA"
	Ed25519Key    KeyType = "Ed25519"
	Secp256k1Key  KeyType = "Secp256k1"
	secp256r1Key  KeyType = "secp256r1"
	Curve25519Key KeyType = "Curve25519"
)

// PublicKey represents publicKey in DID Document,
// used for digital signatures, encryption and other cryptographic operations
type PublicKey struct {
	ID           string  `json:"id"`
	Type         KeyType `json:"type"`
	PublicKeyHex string  `json:"publicKeyHex"`
}

// ServiceType represents service type
type ServiceType string

// Service represent the services provided by DID Subject
type Service struct {
	ID              string      `json:"id"`
	Type            ServiceType `json:"type"`
	ServiceEndpoint string      `json:"serviceEndpoint"`
}

// Proof represents the signature signed by DID Controller
type Proof struct {
	Type           KeyType `json:"type"`
	Creator        string  `json:"creator"`
	SignatureValue string  `json:"signatureValue"`
}

// DDO represents DID Document
type DDO struct {
	Context        string      `json:"@context"`
	ID             string      `json:"id"`
	Version        int8        `json:"version"`
	PublicKey      []PublicKey `json:"publicKey"`
	Controller     string      `json:"controller"`
	Authentication []string    `json:"authentication"`
	Recovery       []string    `json:"recovery"`
	Service        []Service   `json:"service"`
	Proof          *Proof      `json:"proof"`
	Created        time.Time   `json:"created"`
	Updated        time.Time   `json:"updated"`
}

// DIDPackage represents DID context which contains all infomations for DID
type DIDPackage struct {
	Did           string `json:"did"`
	Document      []byte `json:"document"`
	Hash          string `json:"hash"`
	ProviderProof *Proof `json:"providerProof"`
}

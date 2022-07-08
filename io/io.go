package io

import (
	"time"
)

// Key represents the key pair
type Key struct {
	ID            string `json:"id"`
	Type          string `json:"type"`
	PrivateKeyHex string `json:"privateKeyHex"`
	PublicKeyHex  string `json:"publicKeyHex"`
}

// PublicKey represents publicKey in DID Document,
// used for digital signatures, encryption and other cryptographic operations
type PublicKey struct {
	ID           string `json:"id"`
	Type         string `json:"type"`
	PublicKeyHex string `json:"publicKeyHex"`
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
	Type           string `json:"type"`
	Creator        string `json:"creator"`
	SignatureValue string `json:"signatureValue"`
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
	Proof          Proof       `json:"proof"`
	Created        time.Time   `json:"created"`
	Updated        time.Time   `json:"updated"`
}

// CreateDidReq represents the CreateDid request body
type CreateDidReq struct {
	Did      string `json:"did"`
	Document DDO    `json:"document"`
}

// ResolveDidResp represents the ResolveDid response
type ResolveDidResp struct {
	Did      string `json:"did"`
	Document DDO    `json:"document"`
}

// RevokeDidReq represents the ResolveDid request body
type RevokeDidReq struct {
	Did string `json:"did"`
	// Signature to did using recovery prviate key, Base64 encoded
	Signature string `json:"signature"`
}

// Response represents the response body
type Response struct {
	Code int    `json:"code"`
	Data any    `json:"data"`
	Msg  string `json:"msg"`
}

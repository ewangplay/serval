package io

import (
	"fmt"
	"sort"
	"strings"
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

type PublicKeyList []PublicKey

func (kl PublicKeyList) Len() int           { return len(kl) }
func (kl PublicKeyList) Swap(i, j int)      { kl[i], kl[j] = kl[j], kl[i] }
func (kl PublicKeyList) Less(i, j int) bool { return strings.Compare(kl[i].ID, kl[j].ID) == -1 }

func (kl PublicKeyList) MarshalQsign() string {
	sort.Sort(kl)

	s := "["
	for i, k := range kl {
		if i != 0 {
			s += "|"
		}
		s += fmt.Sprintf("id=%s&publicKeyHex=%s&type=%s", k.ID, k.PublicKeyHex, k.Type)
	}
	s += "]"
	return s
}

type StringList []string

func (sl StringList) MarshalQsign() string {
	sort.Strings(sl)

	s := "["
	for i, a := range sl {
		if i != 0 {
			s += "|"
		}
		s += a
	}
	s += "]"
	return s
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
	Context        string        `json:"@context" qsign:"-"`
	ID             string        `json:"id"`
	Version        int8          `json:"version"`
	PublicKey      PublicKeyList `json:"publicKey"`
	Controller     string        `json:"controller"`
	Authentication StringList    `json:"authentication"`
	Recovery       StringList    `json:"recovery"`
	Service        []Service     `json:"service" qsign:"-"`
	Proof          Proof         `json:"proof" qsign:"-"`
	Created        time.Time     `json:"created" qsign:"-"`
	Updated        time.Time     `json:"updated" qsign:"-"`
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

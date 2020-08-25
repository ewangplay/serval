package v1

import (
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	bc "github.com/ewangplay/serval/adapter/blockchain"
	ch "github.com/ewangplay/serval/adapter/cryptohub"
	"github.com/ewangplay/serval/utils"
	"github.com/gin-gonic/gin"
)

// CreateDid handles the /api/v1/did/create request to create a DID
func CreateDid(c *gin.Context) {

	// Get CryptoHub instance
	cryptoHub, err := getCryptoHub(c)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	// Get BlockChain instance
	blockChain, err := getBlockChain(c)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	// Generate DID
	methodName := "gfa"
	methodSpecificID := strings.ReplaceAll(utils.GenerateUUID(), "-", "")
	did := fmt.Sprintf("did:%s:%s", methodName, methodSpecificID)

	// Created time
	now := time.Now()

	// Generate master public / private key pair
	key1 := fmt.Sprintf("%s#keys-1", did)
	pubKey1, priKey1, err := cryptoHub.GenKey()
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	// Use master private key to sign self public key
	signature, err := cryptoHub.Sign(priKey1, pubKey1.GetPublicKey())
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	// Generate standby public / private key pair
	key2 := fmt.Sprintf("%s#keys-2", did)
	pubKey2, priKey2, err := cryptoHub.GenKey()
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	// DID Document
	ddo := &DDO{
		Context: "https://www.w3.org/ns/did/v1",
		ID:      did,
		Version: 1,
		PublicKey: []PublicKey{
			PublicKey{
				ID:           key1,
				Type:         Ed25519Key,
				PublicKeyHex: hex.EncodeToString(pubKey1.GetPublicKey()),
			},
			PublicKey{
				ID:           key2,
				Type:         Ed25519Key,
				PublicKeyHex: hex.EncodeToString(pubKey2.GetPublicKey()),
			},
		},
		Controller:     did,
		Authentication: []string{key1},
		Recovery:       []string{key2},
		Proof: Proof{
			Type:           Ed25519Key,
			Creator:        key1,
			SignatureValue: base64.StdEncoding.EncodeToString(signature),
		},
		Created: now,
		Updated: now,
	}
	fmt.Println("DDO: ", ddo)

	// Submit did / ddo to block chain
	ddoBytes, err := json.Marshal(ddo)
	result, err := blockChain.Submit("CreateDID", did, string(ddoBytes))
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	fmt.Printf("Blockchain submit result: %s\n", string(result))

	// Response body
	respBody := CreateDidResponse{
		Did:     did,
		Created: now,
		Key: []*Key{
			&Key{
				ID:            key1,
				Type:          Ed25519Key,
				PrivateKeyHex: hex.EncodeToString(priKey1.GetPrivateKey()),
				PublicKeyHex:  hex.EncodeToString(pubKey1.GetPublicKey()),
			},
			&Key{
				ID:            key2,
				Type:          Ed25519Key,
				PrivateKeyHex: hex.EncodeToString(priKey2.GetPrivateKey()),
				PublicKeyHex:  hex.EncodeToString(pubKey2.GetPublicKey()),
			},
		},
	}

	c.JSON(http.StatusOK, respBody)
}

func getCryptoHub(c *gin.Context) (ch.CryptoHub, error) {
	obj, exists := c.Get(ch.CryptoHubKey)
	if !exists {
		return nil, fmt.Errorf("cyprto hub instance does not exist in context")
	}
	if obj == nil {
		return nil, fmt.Errorf("cyprto hub instance in context is nil")
	}
	cryptoHub, ok := obj.(ch.CryptoHub)
	if !ok {
		return nil, fmt.Errorf("cyprto hub instance type is invalid")
	}
	return cryptoHub, nil
}

func getBlockChain(c *gin.Context) (bc.BlockChain, error) {
	obj, exists := c.Get(bc.BlockChainKey)
	if !exists {
		return nil, fmt.Errorf("block chain instance does not exist in context")
	}
	if obj == nil {
		return nil, fmt.Errorf("block chain instance in context is nil")
	}
	blockChain, ok := obj.(bc.BlockChain)
	if !ok {
		return nil, fmt.Errorf("block chain instance type is invalid")
	}
	return blockChain, nil
}

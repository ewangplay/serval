package v1

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"
	"time"

	ch "github.com/ewangplay/serval/adapter/cryptohub"
	"github.com/ewangplay/serval/utils"
	"github.com/gin-gonic/gin"
)

// KeyPair represents the public / private key pair structure
type KeyPair struct {
	PrivateKey string `json:"private_key"`
	PublicKey  string `json:"public_key"`
}

// CreateDidResponse represents the response structure
// that requests the creation of did.
type CreateDidResponse struct {
	Did     string    `json:"did"`
	KeyPair *KeyPair  `json:"key_pair"`
	Created time.Time `json:"created"`
}

func getCryptoHub(c *gin.Context) (ch.CryptoHub, error) {
	obj, exists := c.Get(ch.CryptoHubKey)
	if !exists {
		return nil, fmt.Errorf("cyprto hub does not exist in context")
	}
	if obj == nil {
		return nil, fmt.Errorf("cyprto hub in context is nil")
	}
	cryptoHub, ok := obj.(ch.CryptoHub)
	if !ok {
		return nil, fmt.Errorf("cyprto hub type invalid")
	}
	return cryptoHub, nil
}

// CreateDid handles the /api/v1/did/create request to create a DID
func CreateDid(c *gin.Context) {
	// Generate DID
	methodName := "gfa"
	methodSpecificID := strings.ReplaceAll(utils.GenerateUUID(), "-", "")
	did := fmt.Sprintf("did:%s:%s", methodName, methodSpecificID)

	// Created time
	created := time.Now()

	// Generate public / private key pair
	cryptoHub, err := getCryptoHub(c)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	pubKey, priKey, err := cryptoHub.KeyPair()
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	keyPair := &KeyPair{
		PrivateKey: base64.StdEncoding.EncodeToString(priKey.GetPrivateKey()),
		PublicKey:  base64.StdEncoding.EncodeToString(pubKey.GetPublicKey()),
	}

	// Response body
	respBody := CreateDidResponse{
		Did:     did,
		Created: created,
		KeyPair: keyPair,
	}

	c.JSON(http.StatusOK, respBody)
}

package v1

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"
	"time"

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

// CreateDid handles the /api/v1/did/create request to create a DID
func CreateDid(c *gin.Context) {
	// Generate DID
	methodName := "gfa"
	methodSpecificID := strings.ReplaceAll(utils.GenerateUUID(), "-", "")
	did := fmt.Sprintf("did:%s:%s", methodName, methodSpecificID)

	// Created time
	created := time.Now()

	// Generate public / private key pair
	bytePrivateKey := []byte("private key")
	bytePublicKey := []byte("public key")
	keyPair := &KeyPair{
		PrivateKey: base64.StdEncoding.EncodeToString(bytePrivateKey),
		PublicKey:  base64.StdEncoding.EncodeToString(bytePublicKey),
	}

	// Response body
	respBody := CreateDidResponse{
		Did:     did,
		Created: created,
		KeyPair: keyPair,
	}

	c.JSON(http.StatusOK, respBody)
}

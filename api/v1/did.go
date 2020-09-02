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
	"github.com/ewangplay/serval/log"
	"github.com/ewangplay/serval/utils"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

var appKey *Key

func init() {
	appKey = &Key{
		ID:            viper.GetString("appKey.id"),
		Type:          KeyType(viper.GetString("appKey.type")),
		PrivateKeyHex: viper.GetString("appKey.privateKeyHex"),
		PublicKeyHex:  viper.GetString("appKey.publicKeyHex"),
	}
}

// CreateDid handles the /api/v1/did/create request to create a DID
func CreateDid(c *gin.Context) {

	// Generate DID
	methodName := "example"
	methodSpecificID := strings.ReplaceAll(utils.GenerateUUID(), "-", "")
	did := fmt.Sprintf("did:%s:%s", methodName, methodSpecificID)

	// Created time
	now := time.Now()

	// Generate master public / private key pair
	key1 := fmt.Sprintf("%s#keys-1", did)
	pubKey1, priKey1, err := ch.GetCryptoHub().GenKey()
	if err != nil {
		log.Error("Gen master key pair failed: %v", err)
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	// Use master private key to sign did
	// Once an entity's DID is generated,
	// it does not change, so signing did is appropriate.
	signature, err := ch.GetCryptoHub().Sign(priKey1, []byte(did))
	if err != nil {
		log.Error("Self signing failed: %v", err)
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	// Generate standby public / private key pair
	key2 := fmt.Sprintf("%s#keys-2", did)
	pubKey2, priKey2, err := ch.GetCryptoHub().GenKey()
	if err != nil {
		log.Error("Gen standby key pair failed: %v", err)
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
		Proof: &Proof{
			Type:           Ed25519Key,
			Creator:        key1,
			SignatureValue: base64.StdEncoding.EncodeToString(signature),
		},
		Created: now,
		Updated: now,
	}

	log.Debug("DDO: %v", ddo)

	// Hash DID Document
	ddoBytes, _ := json.Marshal(ddo)
	hash := utils.SHA256(ddoBytes)

	// Use application private key to sign DID Document content
	appPriKey, err := hex.DecodeString(appKey.PrivateKeyHex)
	if err != nil {
		log.Error("Load app private key failed: %v", err)
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	signature, err = ch.GetCryptoHub().Sign(ch.Ed25519PrivateKey(appPriKey), []byte(hash))
	if err != nil {
		log.Error("Provider signing failed: %v", err)
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	// Build DID Context
	didPkg := &DIDPackage{
		Did:      did,
		Document: ddoBytes,
		Hash:     hash,
		ProviderProof: &Proof{
			Type:           appKey.Type,
			Creator:        appKey.ID,
			SignatureValue: base64.StdEncoding.EncodeToString(signature),
		},
	}

	// Submit did context to block chain
	didPkgBytes, _ := json.Marshal(didPkg)
	result, err := bc.GetBlockChain().Submit("CreateDID", did, string(didPkgBytes))
	if err != nil {
		log.Error("Submit did document to blockchain failed: %v", err)
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	log.Info("Blockchain submit result: %s", string(result))

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

	log.Debug("CreateDID response: %v", respBody)

	c.JSON(http.StatusOK, respBody)
}

// ResolveDid handles the /api/v1/did/resolve request to resolve a DID
// Request URL: http://IP:Port/api/v1/did/resolve/:did
func ResolveDid(c *gin.Context) {
	// Retrieve did from path param
	did := c.Param("did")

	// Query DDO from blockchain
	result, err := bc.GetBlockChain().Evaluate("QueryDID", did)
	if err != nil {
		log.Error("Query DID Document from blockchain failed: %v", err)
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	log.Debug("Blockchain query result: %v", string(result))

	var didPkg DIDPackage
	err = json.Unmarshal(result, &didPkg)
	if err != nil {
		log.Error("Parse DID Package failed: %v", err)
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	// Verify DID Document
	err = verifyDIDPackage(did, &didPkg)
	if err != nil {
		log.Error("Verify DID Package failed: %v", err)
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	var ddo DDO
	err = json.Unmarshal(didPkg.Document, &ddo)
	if err != nil {
		log.Error("Parse DID Document failed: %v", err)
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	// Response body
	respBody := ResolveDidResponse{
		Did:      did,
		Document: &ddo,
	}

	log.Debug("ResolveDid response: %v", respBody)

	c.JSON(http.StatusOK, respBody)
}

func verifyDIDPackage(did string, didPkg *DIDPackage) error {
	var err error

	// Verify DID Identifier
	if did != didPkg.Did {
		err = fmt.Errorf("got did %v, want %v", didPkg.Did, did)
		return err
	}

	// Verify DID Document hash value
	hash := utils.SHA256(didPkg.Document)
	if hash != didPkg.Hash {
		err = fmt.Errorf("DID Document hash does not match to actual")
		return err
	}

	// Verify provider signature
	if didPkg.ProviderProof == nil {
		err = fmt.Errorf("no provider proof")
		return err
	}
	if didPkg.ProviderProof.Creator != appKey.ID {
		err = fmt.Errorf("provider did does not match")
		return err
	}
	if didPkg.ProviderProof.Type != appKey.Type {
		err = fmt.Errorf("signature algorithm not match")
		return err
	}
	appPubKey, err := hex.DecodeString(appKey.PublicKeyHex)
	if err != nil {
		return err
	}
	signature, err := base64.StdEncoding.DecodeString(didPkg.ProviderProof.SignatureValue)
	if err != nil {
		return err
	}
	valid, err := ch.GetCryptoHub().Verify(ch.Ed25519PublicKey(appPubKey), []byte(hash), signature)
	if err != nil {
		return err
	}
	if !valid {
		err = fmt.Errorf("verify signature fail")
		return err
	}

	return nil
}

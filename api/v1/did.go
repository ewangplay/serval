package v1

import (
	"fmt"
	"net/http"

	"github.com/ewangplay/serval/adapter/store"
	"github.com/ewangplay/serval/io"
	"github.com/ewangplay/serval/log"
	"github.com/gin-gonic/gin"
)

// CreateDid handles the /api/v1/did/create request to create a DID
/*
func CreateDid(c *gin.Context) {

	log.Debug("app key id: %v", viper.GetString("appKey.id"))
	log.Debug("app key type: %v", viper.GetString("appKey.type"))
	log.Debug("app private key: %v", viper.GetString("appKey.privateKeyHex"))
	log.Debug("app public key: %v", viper.GetString("appKey.publicKeyHex"))

	// Generate DID
	methodName := "example"
	methodSpecificID := strings.ReplaceAll(utils.GenerateUUID(), "-", "")
	did := fmt.Sprintf("did:%s:%s", methodName, methodSpecificID)

	// Created time
	now := time.Now()

	// Generate master public / private key pair
	key1 := fmt.Sprintf("%s#keys-1", did)
	priKey1, err := cl.GenEd25519Key()
	if err != nil {
		log.Error("Gen master key failed: %v", err)
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if priKey1.Symmetric() || !priKey1.Private() {
		err = fmt.Errorf("the generated key's type incorrect")
		log.Error("Gen master key failed: %v", err)
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	priKey1Bytes, err := priKey1.Bytes()
	if err != nil {
		log.Error("Gen master key failed: %v", err)
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	pubKey1, err := priKey1.PublicKey()
	if err != nil {
		log.Error("Gen master key failed: %v", err)
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	pubKey1Bytes, err := pubKey1.Bytes()
	if err != nil {
		log.Error("Gen master key failed: %v", err)
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	// Use master private key to sign did
	// Once an entity's DID is generated,
	// it does not change, so signing did is appropriate.
	signature, err := cl.Sign(priKey1, []byte(did))
	if err != nil {
		log.Error("Self signing failed: %v", err)
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	// Generate standby public / private key pair
	key2 := fmt.Sprintf("%s#keys-2", did)
	priKey2, err := cl.GenEd25519Key()
	if err != nil {
		log.Error("Gen standby key failed: %v", err)
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if priKey2.Symmetric() || !priKey2.Private() {
		err = fmt.Errorf("the generated key's type incorrect")
		log.Error("Gen standby key failed: %v", err)
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	priKey2Bytes, err := priKey2.Bytes()
	if err != nil {
		log.Error("Gen standby key failed: %v", err)
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	pubKey2, err := priKey2.PublicKey()
	if err != nil {
		log.Error("Gen standby key failed: %v", err)
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	pubKey2Bytes, err := pubKey2.Bytes()
	if err != nil {
		log.Error("Gen standby key failed: %v", err)
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	// Build DID Document
	ddo := &io.DDO{
		Context: "https://www.w3.org/ns/did/v1",
		ID:      did,
		Version: 1,
		PublicKey: []io.PublicKey{
			{
				ID:           key1,
				Type:         ecl.ED25519,
				PublicKeyHex: hex.EncodeToString(pubKey1Bytes),
			},
			{
				ID:           key2,
				Type:         ecl.ED25519,
				PublicKeyHex: hex.EncodeToString(pubKey2Bytes),
			},
		},
		Controller:     did,
		Authentication: []string{key1},
		Recovery:       []string{key2},
		Proof: &io.Proof{
			Type:           ecl.ED25519,
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
	appPriKeyBytes, err := hex.DecodeString(viper.GetString("appKey.privateKeyHex"))
	if err != nil {
		log.Error("Load app private key failed: %v", err)
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	appPriKey := &ecl.Ed25519PrivateKey{
		PrivKey: appPriKeyBytes,
	}
	signature, err = cl.Sign(appPriKey, []byte(hash))
	if err != nil {
		log.Error("Provider signing failed: %v", err)
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	// Build DID Context
	didPkg := &io.DIDPackage{
		Did:      did,
		Document: ddoBytes,
		Hash:     hash,
		ProviderProof: &io.Proof{
			Type:           viper.GetString("appKey.type"),
			Creator:        viper.GetString("appKey.id"),
			SignatureValue: base64.StdEncoding.EncodeToString(signature),
		},
	}

	// Set to store
	err = store.Set(did, didPkg)
	if err != nil {
		log.Error("Set did document to Store failed: %v", err)
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	// Response body
	respBody := io.CreateDidResp{
		Did:     did,
		Created: now,
		Key: []*io.Key{
			{
				ID:            key1,
				Type:          ecl.ED25519,
				PrivateKeyHex: hex.EncodeToString(priKey1Bytes),
				PublicKeyHex:  hex.EncodeToString(pubKey1Bytes),
			},
			{
				ID:            key2,
				Type:          ecl.ED25519,
				PrivateKeyHex: hex.EncodeToString(priKey2Bytes),
				PublicKeyHex:  hex.EncodeToString(pubKey2Bytes),
			},
		},
	}

	log.Debug("CreateDID response: %v", respBody)

	c.JSON(http.StatusOK, respBody)
}
*/

func CreateDid(c *gin.Context) {
	var err error

	var req io.CreateDidReq
	err = c.BindJSON(&req)
	if err != nil {
		errMsg := fmt.Sprintf("Parse request body failed: %v", err)
		log.Error(errMsg)
		FailWithMessage(http.StatusBadRequest, errMsg, c)
		return
	}

	log.Debug("Did: %v\nDocument: %v", req.Did, req.Document)

	// TODO: Verify the signature in DDO

	// Set to store
	err = store.Set(req.Did, req.Document)
	if err != nil {
		errMsg := fmt.Sprintf("Set did document to Store failed: %v", err)
		log.Error(errMsg)
		FailWithMessage(http.StatusInternalServerError, errMsg, c)
		return
	}

	Ok(c)
}

// ResolveDid handles the /api/v1/did/resolve request to resolve a DID
// Request URL: http://IP:Port/api/v1/did/resolve/:did
func ResolveDid(c *gin.Context) {
	// Retrieve did from path param
	did := c.Param("did")

	// Query DDO from Store
	var ddo io.DDO
	found, err := store.Get(did, &ddo)
	if err != nil {
		errMsg := fmt.Sprintf("Get DID(%v) Document from Store failed: %v", did, err)
		log.Error(errMsg)
		FailWithMessage(http.StatusInternalServerError, errMsg, c)
		return
	}
	if !found {
		errMsg := fmt.Sprintf("Did %v not found", did)
		log.Error(errMsg)
		FailWithMessage(http.StatusInternalServerError, errMsg, c)
		return
	}

	// Response body
	resp := io.ResolveDidResp{
		Did:      did,
		Document: ddo,
	}

	log.Debug("ResolveDid response: %v", resp)

	OkWithData(resp, c)
}

// RevokeDid handles the /api/v1/did/revoke request to revoke a DID
func RevokeDid(c *gin.Context) {
	var err error

	// Parse the request body
	var req io.RevokeDidReq
	err = c.BindJSON(&req)
	if err != nil {
		errMsg := fmt.Sprintf("Parse request body failed: %v", err)
		log.Error(errMsg)
		FailWithMessage(http.StatusBadRequest, errMsg, c)
		return
	}

	// Check the params
	if req.Did == "" {
		errMsg := fmt.Sprintf("The request parameter did is empty")
		log.Error(errMsg)
		FailWithMessage(http.StatusBadRequest, errMsg, c)
		return
	}
	if req.Signature == "" {
		errMsg := fmt.Sprintf("The request parameter signature is empty")
		log.Error(errMsg)
		FailWithMessage(http.StatusBadRequest, errMsg, c)
		return
	}

	// TODO: Verify the signature
	// Client uses the recovery key to make the signature, so we must use the
	// corresponding public key to verify the signature.

	// Delete the Did/DDO record from Store
	err = store.Delete(req.Did)
	if err != nil {
		errMsg := fmt.Sprintf("Delete %s item from Store failed: %v", req.Did, err)
		log.Error(errMsg)
		FailWithMessage(http.StatusInternalServerError, errMsg, c)
		return
	}

	Ok(c)
}

package v1

import (
	"fmt"
	"net/http"

	acl "github.com/ewangplay/serval/adapter/cryptolib"
	"github.com/ewangplay/serval/adapter/store"
	"github.com/ewangplay/serval/io"
	"github.com/ewangplay/serval/log"
	"github.com/gin-gonic/gin"
)

// CreateDid handles the /api/v1/did/create request to create a DID
func CreateDid(c *gin.Context) {
	var err error

	// Parse the request body
	req, err := parseCreateDidReq(c)
	if err != nil {
		errMsg := fmt.Sprintf("Parse the request body failed: %v", err)
		log.Error(errMsg)
		FailWithMessage(http.StatusBadRequest, errMsg, c)
		return
	}

	// Set the DID/DDO record to store
	err = store.Set(req.Did, req.Document)
	if err != nil {
		errMsg := fmt.Sprintf("Set the DID/DDO record to store failed: %v", err)
		log.Error(errMsg)
		FailWithMessage(http.StatusInternalServerError, errMsg, c)
		return
	}

	Ok(c)
}

func parseCreateDidReq(c *gin.Context) (*io.CreateDidReq, error) {
	var err error
	var req io.CreateDidReq

	err = c.BindJSON(&req)
	if err != nil {
		return nil, err
	}

	log.Debug("Did: %v\nDocument: %v", req.Did, req.Document)

	// Verify the DID document
	valid, err := acl.VerifyDDO(&req.Document)
	if err != nil {
		return nil, err
	}
	if !valid {
		return nil, fmt.Errorf("signature verifying failed")
	}

	return &req, nil
}

// ResolveDid handles the /api/v1/did/resolve request to resolve a DID
// Request URL: http://IP:Port/api/v1/did/resolve/:did
func ResolveDid(c *gin.Context) {
	// Retrieve did from path param
	did := c.Param("did")

	// Get the DID/DDO record from store
	var ddo io.DDO
	found, err := store.Get(did, &ddo)
	if err != nil {
		errMsg := fmt.Sprintf("Get the DID/DDO (%v) record from store failed: %v", did, err)
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
	req, err := parseRevokeDidReq(c)
	if err != nil {
		errMsg := fmt.Sprintf("Parse the request body failed: %v", err)
		log.Error(errMsg)
		FailWithMessage(http.StatusBadRequest, errMsg, c)
		return
	}

	// Delete the DID/DDO record from store
	err = store.Delete(req.Did)
	if err != nil {
		errMsg := fmt.Sprintf("Delete the DID/DDO (%s) record from store failed: %v", req.Did, err)
		log.Error(errMsg)
		FailWithMessage(http.StatusInternalServerError, errMsg, c)
		return
	}

	Ok(c)
}

func parseRevokeDidReq(c *gin.Context) (*io.RevokeDidReq, error) {
	var err error
	var req io.RevokeDidReq

	err = c.BindJSON(&req)
	if err != nil {
		return nil, err
	}

	// Check the params
	if req.Did == "" {
		err = fmt.Errorf("The request parameter did cannot be empty")
		return nil, err
	}
	if req.Signature == "" {
		err = fmt.Errorf("The request parameter signature cannot be empty")
		return nil, err
	}

	// TODO: Verify the signature
	// Client uses the recovery key to make the signature, so we must use the
	// corresponding public key to verify the signature.

	return &req, nil
}

package v1

import (
	"encoding/json"
	"fmt"
	"net/http"

	ctx "github.com/ewangplay/serval/context"
	"github.com/ewangplay/serval/io"
	"github.com/ewangplay/serval/log"
	"github.com/ewangplay/serval/utils"
)

// CreateDid handles the /api/v1/did/create request to create a DID
func CreateDid(c *ctx.Context) {
	var err error

	// Parse the request body
	req, err := parseCreateDidReq(c)
	if err != nil {
		errMsg := fmt.Sprintf("Parse the request body failed: %v", err)
		log.Error(errMsg)
		FailWithMessage(http.StatusBadRequest, errMsg, c.Context)
		return
	}

	// debug
	data, _ := json.Marshal(req)
	log.Debug("CreateDid request: %s", string(data))

	// Set the DID/DDO record to store
	err = c.Store.Set(req.Did, req.Document)
	if err != nil {
		errMsg := fmt.Sprintf("Set the DID/DDO record to store failed: %v", err)
		log.Error(errMsg)
		FailWithMessage(http.StatusInternalServerError, errMsg, c.Context)
		return
	}

	Ok(c.Context)
}

func parseCreateDidReq(c *ctx.Context) (*io.CreateDidReq, error) {
	var err error
	var req io.CreateDidReq

	err = c.BindJSON(&req)
	if err != nil {
		return nil, err
	}

	// Verify the DID document
	err = utils.VerifyDDO(c.CSP, c.Qsign, &req.Document)
	if err != nil {
		return nil, err
	}

	return &req, nil
}

// ResolveDid handles the /api/v1/did/resolve request to resolve a DID
// Request URL: http://IP:Port/api/v1/did/resolve/:did
func ResolveDid(c *ctx.Context) {
	// Retrieve did from path param
	did := c.Param("did")

	// Get the DID/DDO record from store
	var ddo io.DDO
	found, err := c.Store.Get(did, &ddo)
	if err != nil {
		errMsg := fmt.Sprintf("Failed to retrieve the DID document (%v) from store: %v", did, err)
		log.Error(errMsg)
		FailWithMessage(http.StatusInternalServerError, errMsg, c.Context)
		return
	}
	if !found {
		errMsg := fmt.Sprintf("DID document (%v) not found", did)
		log.Error(errMsg)
		FailWithMessage(http.StatusInternalServerError, errMsg, c.Context)
		return
	}

	// Verify the DID document
	err = utils.VerifyDDO(c.CSP, c.Qsign, &ddo)
	if err != nil {
		errMsg := fmt.Sprintf("Failed to verify the DID document (%v): %v", did, err)
		log.Error(errMsg)
		FailWithMessage(http.StatusInternalServerError, errMsg, c.Context)
		return
	}

	// Response body
	resp := io.ResolveDidResp{
		Did:      did,
		Document: ddo,
	}

	log.Debug("ResolveDid response: %v", resp)

	OkWithData(resp, c.Context)
}

// RevokeDid handles the /api/v1/did/revoke request to revoke a DID
func RevokeDid(c *ctx.Context) {
	var err error

	// Parse the request body
	req, err := parseRevokeDidReq(c)
	if err != nil {
		errMsg := fmt.Sprintf("Parse the request body failed: %v", err)
		log.Error(errMsg)
		FailWithMessage(http.StatusBadRequest, errMsg, c.Context)
		return
	}

	// debug
	data, _ := json.Marshal(req)
	log.Debug("RevokeDid request: %s", string(data))

	// Delete the DID/DDO record from store
	err = c.Store.Delete(req.Did)
	if err != nil {
		errMsg := fmt.Sprintf("Delete the DID/DDO (%s) record from store failed: %v", req.Did, err)
		log.Error(errMsg)
		FailWithMessage(http.StatusInternalServerError, errMsg, c.Context)
		return
	}

	Ok(c.Context)
}

func parseRevokeDidReq(c *ctx.Context) (*io.RevokeDidReq, error) {
	var err error
	var req io.RevokeDidReq

	err = c.BindJSON(&req)
	if err != nil {
		return nil, err
	}

	// Check the params
	if req.Did == "" {
		err = fmt.Errorf("The DID parameter cannot be empty")
		return nil, err
	}
	if req.Proof.Type == "" || req.Proof.Creator == "" || req.Proof.SignatureValue == "" {
		err = fmt.Errorf("The Proof parameter cannot be empty")
		return nil, err
	}

	// Verify the proof
	var ddo io.DDO
	found, err := c.Store.Get(req.Did, &ddo)
	if err != nil {
		return nil, err
	}
	if !found {
		err = fmt.Errorf("DID document (%v) not found", req.Did)
		return nil, err
	}

	// Verify the DID document
	err = utils.VerifyDDO(c.CSP, c.Qsign, &ddo)
	if err != nil {
		return nil, err
	}

	valid, err := utils.VerifyProof(c.CSP, req.Did, &req.Proof, &ddo)
	if err != nil {
		return nil, err
	}
	if !valid {
		return nil, fmt.Errorf("Failed to verify the signature of the Proof")
	}

	return &req, nil
}

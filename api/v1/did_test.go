package v1

import (
	"crypto/ed25519"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	bc "github.com/ewangplay/serval/adapter/blockchain"
	ch "github.com/ewangplay/serval/adapter/cryptohub"
	"github.com/gin-gonic/gin"
)

// mockBlockChain represents the mock block chain
type mockBlockChain struct {
}

// Submit will submit a transaction to the ledger
func (m *mockBlockChain) Submit(fn string, args ...string) ([]byte, error) {
	return []byte("TxID-0001"), nil
}

// Evaluate will evaluate a transaction function and return its results
func (m *mockBlockChain) Evaluate(fn string, args ...string) ([]byte, error) {
	did := "did:gfa:54f2c1506bc54747a5d36b55b240a2ac"
	key1 := fmt.Sprintf("%s#keys-1", did)
	key2 := fmt.Sprintf("%s#keys-2", did)
	now := time.Now()
	ddo := &DDO{
		Context: "https://www.w3.org/ns/did/v1",
		ID:      did,
		Version: 1,
		PublicKey: []PublicKey{
			PublicKey{
				ID:           key1,
				Type:         Ed25519Key,
				PublicKeyHex: hex.EncodeToString([]byte("public key 1")),
			},
			PublicKey{
				ID:           key2,
				Type:         Ed25519Key,
				PublicKeyHex: hex.EncodeToString([]byte("public key 2")),
			},
		},
		Controller:     did,
		Authentication: []string{key1},
		Recovery:       []string{key2},
		Proof: Proof{
			Type:           Ed25519Key,
			Creator:        key1,
			SignatureValue: base64.StdEncoding.EncodeToString([]byte("signature")),
		},
		Created: now,
		Updated: now,
	}
	ddoBytes, _ := json.Marshal(ddo)
	return ddoBytes, nil
}

func TestCreateDidSucc(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Prepare request params
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set(ch.CryptoHubKey, ch.GetCryptoHub())
	c.Set(bc.BlockChainKey, &mockBlockChain{})
	c.Request = httptest.NewRequest("POST", "http://localhost:8080/api/v1/did/create", nil)

	// Call target handler
	CreateDid(c)

	// Verify the results returned
	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("response status code got %v, want %v", resp.StatusCode, http.StatusOK)
	}
	contentType := resp.Header.Get("Content-Type")
	if !strings.Contains(contentType, "application/json; charset=utf-8") {
		t.Fatalf("response content type shoud contain %v", "application/json; charset=utf-8")
	}
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))
	var respBody CreateDidResponse
	err := json.Unmarshal(body, &respBody)
	if err != nil {
		t.Fatalf("response body can't be resolved, %v", err)
	}
}

func TestCreateDidWithoutCryptoHub(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Prepare request params
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "http://localhost:8080/api/v1/did/create", nil)

	// Call target handler
	CreateDid(c)

	// Verify the results returned
	resp := w.Result()
	if resp.StatusCode != http.StatusInternalServerError {
		t.Fatalf("response status code got %v, want %v", resp.StatusCode, http.StatusInternalServerError)
	}
}

func TestCreateDidWithCryptoHubIsNil(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Prepare request params
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set(ch.CryptoHubKey, nil)
	c.Request = httptest.NewRequest("POST", "http://localhost:8080/api/v1/did/create", nil)

	// Call target handler
	CreateDid(c)

	// Verify the results returned
	resp := w.Result()
	if resp.StatusCode != http.StatusInternalServerError {
		t.Fatalf("response status code got %v, want %v", resp.StatusCode, http.StatusInternalServerError)
	}
}

func TestCreateDidWithCryptoHubTypeIncorret(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Prepare request params
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set(ch.CryptoHubKey, "I am not a crypto hub")
	c.Request = httptest.NewRequest("POST", "http://localhost:8080/api/v1/did/create", nil)

	// Call target handler
	CreateDid(c)

	// Verify the results returned
	resp := w.Result()
	if resp.StatusCode != http.StatusInternalServerError {
		t.Fatalf("response status code got %v, want %v", resp.StatusCode, http.StatusInternalServerError)
	}
}

func TestCreateDidWithoutBlockChain(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Prepare request params
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set(ch.CryptoHubKey, ch.GetCryptoHub())
	c.Request = httptest.NewRequest("POST", "http://localhost:8080/api/v1/did/create", nil)

	// Call target handler
	CreateDid(c)

	// Verify the results returned
	resp := w.Result()
	if resp.StatusCode != http.StatusInternalServerError {
		t.Fatalf("response status code got %v, want %v", resp.StatusCode, http.StatusInternalServerError)
	}
}

func TestCreateDidWithBlockChainIsNil(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Prepare request params
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set(ch.CryptoHubKey, ch.GetCryptoHub())
	c.Set(bc.BlockChainKey, nil)
	c.Request = httptest.NewRequest("POST", "http://localhost:8080/api/v1/did/create", nil)

	// Call target handler
	CreateDid(c)

	// Verify the results returned
	resp := w.Result()
	if resp.StatusCode != http.StatusInternalServerError {
		t.Fatalf("response status code got %v, want %v", resp.StatusCode, http.StatusInternalServerError)
	}
}

func TestCreateDidWithBlockChainTypeIncorret(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Prepare request params
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set(ch.CryptoHubKey, ch.GetCryptoHub())
	c.Set(bc.BlockChainKey, "I am not a blockchain")
	c.Request = httptest.NewRequest("POST", "http://localhost:8080/api/v1/did/create", nil)

	// Call target handler
	CreateDid(c)

	// Verify the results returned
	resp := w.Result()
	if resp.StatusCode != http.StatusInternalServerError {
		t.Fatalf("response status code got %v, want %v", resp.StatusCode, http.StatusInternalServerError)
	}
}

// mockCryptoHubGenKeyError represents the mock crypto hub with GenKey failed
type mockCryptoHubGenKeyError struct {
}

// GenKey returns error
func (ed *mockCryptoHubGenKeyError) GenKey() (ch.PublicKey, ch.PrivateKey, error) {
	return nil, nil, fmt.Errorf("generate key error")
}

// Sign signs the message with privateKey and returns signature
func (ed *mockCryptoHubGenKeyError) Sign(privateKey ch.PrivateKey, message []byte) (sigature []byte, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = fmt.Errorf("sign error: %v", e)
		}
	}()
	return ed25519.Sign(ed25519.PrivateKey(privateKey.GetPrivateKey()), message), nil
}

func TestCreateDidWithCryptoHubGenKeyError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Prepare request params
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set(ch.CryptoHubKey, &mockCryptoHubGenKeyError{})
	c.Set(bc.BlockChainKey, &mockBlockChain{})
	c.Request = httptest.NewRequest("POST", "http://localhost:8080/api/v1/did/create", nil)

	// Call target handler
	CreateDid(c)

	// Verify the results returned
	resp := w.Result()
	if resp.StatusCode != http.StatusInternalServerError {
		t.Fatalf("response status code got %v, want %v", resp.StatusCode, http.StatusInternalServerError)
	}
}

// mockBlockChainSubmitError represents the mock block chain with submit error
type mockBlockChainSubmitError struct {
}

// Submit will submit a transaction to the ledger
func (m *mockBlockChainSubmitError) Submit(fn string, args ...string) ([]byte, error) {
	return nil, fmt.Errorf("block chain submit error")
}

// Evaluate will evaluate a transaction function and return its results
func (m *mockBlockChainSubmitError) Evaluate(fn string, args ...string) ([]byte, error) {
	did := "did:gfa:54f2c1506bc54747a5d36b55b240a2ac"
	key1 := fmt.Sprintf("%s#keys-1", did)
	key2 := fmt.Sprintf("%s#keys-2", did)
	now := time.Now()
	ddo := &DDO{
		Context: "https://www.w3.org/ns/did/v1",
		ID:      did,
		Version: 1,
		PublicKey: []PublicKey{
			PublicKey{
				ID:           key1,
				Type:         Ed25519Key,
				PublicKeyHex: hex.EncodeToString([]byte("public key 1")),
			},
			PublicKey{
				ID:           key2,
				Type:         Ed25519Key,
				PublicKeyHex: hex.EncodeToString([]byte("public key 2")),
			},
		},
		Controller:     did,
		Authentication: []string{key1},
		Recovery:       []string{key2},
		Proof: Proof{
			Type:           Ed25519Key,
			Creator:        key1,
			SignatureValue: base64.StdEncoding.EncodeToString([]byte("signature")),
		},
		Created: now,
		Updated: now,
	}
	ddoBytes, _ := json.Marshal(ddo)
	return ddoBytes, nil
}

func TestCreateDidWithBlockChainSubmitError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Prepare request params
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set(ch.CryptoHubKey, ch.GetCryptoHub())
	c.Set(bc.BlockChainKey, &mockBlockChainSubmitError{})
	c.Request = httptest.NewRequest("POST", "http://localhost:8080/api/v1/did/create", nil)

	// Call target handler
	CreateDid(c)

	// Verify the results returned
	resp := w.Result()
	if resp.StatusCode != http.StatusInternalServerError {
		t.Fatalf("response status code got %v, want %v", resp.StatusCode, http.StatusInternalServerError)
	}
}

func TestResolveDidSucc(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Prepare request params
	did := "did:gfa:54f2c1506bc54747a5d36b55b240a2ac"
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	// Set blockchain mock instance
	c.Set(bc.BlockChainKey, &mockBlockChain{})
	// Set path param: did
	c.Params = append(c.Params, gin.Param{Key: "did", Value: did})
	c.Request = httptest.NewRequest("GET", "http://localhost:8080/api/v1/did/resolve/did:gfa:54f2c1506bc54747a5d36b55b240a2ac", nil)

	// Call target handler
	ResolveDid(c)

	// Verify the results returned
	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("response status code got %v, want %v", resp.StatusCode, http.StatusOK)
	}
	contentType := resp.Header.Get("Content-Type")
	if !strings.Contains(contentType, "application/json; charset=utf-8") {
		t.Fatalf("response content type shoud contain %v", "application/json; charset=utf-8")
	}
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))
	var respBody ResolveDidResponse
	err := json.Unmarshal(body, &respBody)
	if err != nil {
		t.Fatalf("response body can't be resolved, %v", err)
	}
	if respBody.Did != did {
		t.Fatalf("response did got %v, want %v", respBody.Did, did)
	}
}

func TestResolveDidWithoutBlockchain(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Prepare request params
	did := "did:gfa:54f2c1506bc54747a5d36b55b240a2ac"
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	// Set path param: did
	c.Params = append(c.Params, gin.Param{Key: "did", Value: did})
	c.Request = httptest.NewRequest("GET", "http://localhost:8080/api/v1/did/resolve/did:gfa:54f2c1506bc54747a5d36b55b240a2ac", nil)

	// Call target handler
	ResolveDid(c)

	// Verify the results returned
	resp := w.Result()
	if resp.StatusCode != http.StatusInternalServerError {
		t.Fatalf("response status code got %v, want %v", resp.StatusCode, http.StatusInternalServerError)
	}
}

// mockBlockChain represents the mock block chain with evaluate error
type mockBlockChainEvaluateError struct {
}

// Submit will submit a transaction to the ledger
func (m *mockBlockChainEvaluateError) Submit(fn string, args ...string) ([]byte, error) {
	return []byte("TxID-0001"), nil
}

// Evaluate will evaluate a transaction function and return its results
func (m *mockBlockChainEvaluateError) Evaluate(fn string, args ...string) ([]byte, error) {
	return nil, fmt.Errorf("blockchain evaluate error")
}

func TestResolveDidWithBlockChainEvaluateError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Prepare request params
	did := "did:gfa:54f2c1506bc54747a5d36b55b240a2ac"
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	// Set blockchain mock instance with evalute error
	c.Set(bc.BlockChainKey, &mockBlockChainEvaluateError{})
	// Set path param: did
	c.Params = append(c.Params, gin.Param{Key: "did", Value: did})
	c.Request = httptest.NewRequest("GET", "http://localhost:8080/api/v1/did/resolve/did:gfa:54f2c1506bc54747a5d36b55b240a2ac", nil)

	// Call target handler
	ResolveDid(c)

	// Verify the results returned
	resp := w.Result()
	if resp.StatusCode != http.StatusInternalServerError {
		t.Fatalf("response status code got %v, want %v", resp.StatusCode, http.StatusInternalServerError)
	}
}

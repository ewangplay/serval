package v1

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	ch "github.com/ewangplay/serval/adapter/cryptohub"
	"github.com/gin-gonic/gin"
)

func TestCreateDidSucc(t *testing.T) {
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

// MockCryptoHubKeyPairError represents the mock crypto hub
type MockCryptoHubKeyPairError struct {
}

// KeyPair mock generating key pair failure
func (ed *MockCryptoHubKeyPairError) KeyPair() (ch.PublicKey, ch.PrivateKey, error) {
	return nil, nil, errors.New("generating key pair fail")
}

func TestCreateDidWithKeyPairError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Prepare request params
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set(ch.CryptoHubKey, &MockCryptoHubKeyPairError{})
	c.Request = httptest.NewRequest("POST", "http://localhost:8080/api/v1/did/create", nil)

	// Call target handler
	CreateDid(c)

	// Verify the results returned
	resp := w.Result()
	if resp.StatusCode != http.StatusInternalServerError {
		t.Fatalf("response status code got %v, want %v", resp.StatusCode, http.StatusInternalServerError)
	}
}

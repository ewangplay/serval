package v1

import (
	"fmt"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestCreateDid(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Prepare request params
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
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

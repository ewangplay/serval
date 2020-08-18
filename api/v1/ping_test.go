package v1

import (
	"fmt"
	"io/ioutil"
	"net/http/httptest"

	"github.com/gin-gonic/gin"
)

func ExamplePong() {
	gin.SetMode(gin.TestMode)

	// Prepare request params
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "http://localhost:8080/api/v1/ping", nil)

	// Call target handler
	Pong(c)

	// Verify the results returned
	resp := w.Result()
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(resp.StatusCode)
	fmt.Println(resp.Header.Get("Content-Type"))
	fmt.Println(string(body))
	// Output:
	// 200
	// application/json; charset=utf-8
	// {"message":"pong"}
}

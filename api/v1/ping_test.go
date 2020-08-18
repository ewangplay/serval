package v1

import (
	"fmt"
	"io/ioutil"
	"net/http/httptest"

	"github.com/gin-gonic/gin"
)

func ExamplePong() {
	w := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(w)

	c.Request = httptest.NewRequest("GET", "http://localhost:8080/ping", nil)
	Pong(c)

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

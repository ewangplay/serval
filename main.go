package main

import "github.com/ewangplay/serval/router"

func main() {
	r := router.InitRouter()

	// listen and serve on 0.0.0.0:8080
	r.Run()
}

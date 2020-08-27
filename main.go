package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/ewangplay/serval/config"
	"github.com/ewangplay/serval/router"
	"github.com/spf13/viper"
)

func main() {
	var filename = flag.String("config", "serval.yaml", "path to config file")
	flag.Parse()

	// Init configure
	err := config.InitConfig(*filename)
	if err != nil {
		fmt.Printf("Init config failed: %v\n", err)
		os.Exit(1)
	}

	// Init router
	r := router.InitRouter()

	// listen and serve on 0.0.0.0:<port>
	r.Run(fmt.Sprintf(":%s", viper.GetString("server.port")))
}

package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/ewangplay/serval/config"
	"github.com/ewangplay/serval/log"
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

	// New Rotate Writer
	var w io.Writer
	rwCfg := &log.RotateWriterConfig{
		Module:      "serval",
		Path:        viper.GetString("log.path"),
		MaxSize:     viper.GetInt64("log.maxSize"),
		RotateDaily: viper.GetBool("log.rotateDaily"),
	}
	w, err = log.NewRotateWriter(rwCfg)
	if err != nil {
		fmt.Printf("Create rotate writer failed: %v\n", err)
		os.Exit(1)
	}

	// Init Logger
	logCfg := &log.LoggerConfig{
		Module:   "serval",
		LogLevel: viper.GetString("log.level"),
		Color:    0,
		Writer:   w,
	}
	err = log.InitLogger(logCfg)
	if err != nil {
		fmt.Printf("Init logger failed: %v\n", err)
		os.Exit(1)
	}

	// Init router
	r := router.InitRouter(w)

	// listen and serve on 0.0.0.0:<port>
	r.Run(fmt.Sprintf(":%s", viper.GetString("server.port")))
}

package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/ewangplay/rwriter"
	"github.com/ewangplay/serval/adapter"
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

	// Debug prints all configuration registries for debugging
	viper.Debug()

	// New Rotate Writer
	var w io.Writer
	rwCfg := &rwriter.Config{
		Module:      "serval",
		Path:        viper.GetString("log.path"),
		MaxSize:     viper.GetInt64("log.maxSize"),
		RotateDaily: viper.GetBool("log.rotateDaily"),
	}
	w, err = rwriter.NewRotateWriter(rwCfg)
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

	// Init Store
	var opts adapter.StoreOptions
	err = viper.UnmarshalKey("store", &opts)
	if err != nil {
		fmt.Printf("Read store options failed: %v\n", err)
		os.Exit(1)
	}
	store, err := adapter.InitStore(&opts)
	if err != nil {
		fmt.Printf("Init store failed: %v\n", err)
		os.Exit(1)
	}
	defer store.Close()

	// Init cryptolib
	csp, err := adapter.InitCryptolib()
	if err != nil {
		fmt.Printf("Init cryptolib failed: %v\n", err)
		os.Exit(1)
	}

	// Init Qsign
	qsign, err := adapter.InitQsign()
	if err != nil {
		fmt.Printf("Init Qsign failed: %v\n", err)
		os.Exit(1)
	}

	// Init router
	r := router.InitRouter(w, store, csp, qsign)

	// listen and serve on 0.0.0.0:<port>
	r.Run(fmt.Sprintf(":%s", viper.GetString("server.port")))
}

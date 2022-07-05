package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/ewangplay/rwriter"
	cl "github.com/ewangplay/serval/adapter/cryptolib"
	"github.com/ewangplay/serval/adapter/store"
	st "github.com/ewangplay/serval/adapter/store"
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
	var opts store.Options
	err = viper.UnmarshalKey("store", &opts)
	if err != nil {
		fmt.Printf("Read store options failed: %v\n", err)
		os.Exit(1)
	}
	err = st.InitStore(&opts)
	if err != nil {
		fmt.Printf("Init store failed: %v\n", err)
		os.Exit(1)
	}
	defer st.ReleaseStore()

	// Init cryptolib
	err = cl.InitCryptolib()
	if err != nil {
		fmt.Printf("Init cryptolib failed: %v\n", err)
		os.Exit(1)
	}

	// Init router
	r := router.InitRouter(w)

	// listen and serve on 0.0.0.0:<port>
	r.Run(fmt.Sprintf(":%s", viper.GetString("server.port")))
}

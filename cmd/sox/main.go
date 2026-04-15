package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/orvice/sox/internal/biz"
	"github.com/orvice/sox/internal/config"
	"github.com/weeon/log"
	"github.com/weeon/utils"
	"github.com/weeon/utils/process"
)

var (
	configFile string
)

func main() {
	var err error

	flag.StringVar(&configFile, "c", "config", "config file path")
	flag.Parse()

	log.SetupStdoutLogger()
	// 初始化
	log.Infof("Run")
	err = utils.ExecFuncs([]func() error{
		func() error {
			return config.Init(configFile)
		},
		biz.InitCache,
		biz.InitMastodon,
		biz.InitTwitter,
	})
	if err != nil {
		fmt.Println("Init Error: ", err)
		os.Exit(1)
	}

	go biz.Daemon(context.Background())
	process.WaitSignal()
}

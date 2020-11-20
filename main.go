package main

import (
	"github.com/alecthomas/kingpin"
	"github.com/ozonru/file.d/cfg"
	"github.com/ozonru/file.d/fd"
	"github.com/ozonru/file.d/logger"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/daniel-orlov/output-plugin-file.d"
	_ "github.com/daniel-orlov/parser-plugin-file.d"
	//_ "github.com/ozonru/file.d/plugin/action/discard"
	_ "github.com/ozonru/file.d/plugin/input/file"
	//_ "github.com/ozonru/file.d/plugin/input/http"
	//_ "github.com/ozonru/file.d/plugin/input/kafka"
	//_ "github.com/ozonru/file.d/plugin/output/devnull"
	//_ "github.com/ozonru/file.d/plugin/output/elasticsearch"
	//_ "github.com/ozonru/file.d/plugin/output/gelf"
	//_ "github.com/ozonru/file.d/plugin/output/kafka"
	//_ "github.com/ozonru/file.d/plugin/output/stdout"
)

var (
	fileD   *fd.FileD
	exit    = make(chan bool)
	version = "v0.0.1"

	config = kingpin.Flag("config", `config file name`).Required().ExistingFile()
	http   = kingpin.Flag("http", `http listen addr eg. ":9000", "off" to disable`).Default(":9000").String()
)

func main() {
	kingpin.Version(version)
	kingpin.Parse()

	go listenSignals()
	go start()

	<-exit
	logger.Infof("Exiting...")
}

func start() {
	fileD = fd.New(cfg.NewConfigFromFile(*config), *http)
	fileD.Start()
}

func listenSignals() {
	signalChan := make(chan os.Signal)
	signal.Notify(signalChan, syscall.SIGHUP, syscall.SIGTERM)

	for {
		s := <-signalChan

		switch s {
		case syscall.SIGHUP:
			logger.Infof("SIGHUP received")
			fileD.Stop()
			start()
		case syscall.SIGINT:
			fallthrough
		case syscall.SIGTERM:
			logger.Infof("SIGTERM received")
			fileD.Stop()
			exit <- true
		}
	}
}

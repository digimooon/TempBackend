package main

import (
	"TempBackend/model"
	"TempBackend/server"
	"flag"
	"io/ioutil"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

var (
	configFilePath = flag.String("config", "etc/config.yaml", "temperature backend config file path")
)

func main() {
	flag.Parse()
	configData, err := ioutil.ReadFile(*configFilePath)
	if err != nil {
		println("read config file failed, err: " + err.Error())
		os.Exit(1)
	}

	cfg, err := model.UnmarshalConfig(configData)
	if err != nil {
		println("parse config file failed, err: " + err.Error())
		os.Exit(1)
	}

	srv, err := server.Init(cfg)
	if err != nil {
		println("init failed, err: " + err.Error())
		os.Exit(1)
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGPIPE, syscall.SIGUSR1)
	wg := new(sync.WaitGroup)
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			sig := <-c
			if sig == syscall.SIGINT || sig == syscall.SIGTERM || sig == syscall.SIGQUIT {
				println("got signal, quit, signal: " + sig.String())
				srv.Close()
				break
			}
			println("ignore signal: " + sig.String())
		}
	}()

	srv.Run()
	wg.Wait()
}

package main

import (
	"flag"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"udns/config"
	"udns/handler"
	"udns/logger"
	"udns/storage"
	"udns/watcher"

	"github.com/miekg/dns"
)

func main() {
	var cfgFile string

	flag.StringVar(&cfgFile, "cfg", "", "config file path, eg. /etc/udns.cfg.yml")
	flag.Parse()

	config.Init(cfgFile)
	storage.Init(config.Cfg.Source)

	addr := ":" + strconv.Itoa(config.Cfg.Port)
	logger.Infof("main", "Listen to %s(%s), parent DNS '%s'", addr, config.Cfg.Proto, config.Cfg.ParentDNS)
	server := &dns.Server{Addr: addr, Net: config.Cfg.Proto}
	for _, addr := range config.Cfg.Zones {
		dns.HandleFunc(addr, handler.OnSmartDNSRequest)
	}
	dns.HandleFunc(config.Cfg.MyIP, handler.OnMyIPRequest)
	go func() {
		logger.Info("main", "Server start to listen...")
		err := server.ListenAndServe()
		if err != nil {
			logger.Fatal("main", err)
		}
	}()
	sig := make(chan os.Signal)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	select {
	case s := <-sig:
		logger.Infof("main", "Signal (%s) received, stopping...", s)
		if err := server.Shutdown(); err != nil {
			logger.Error("main", err)
		}
		storage.Shutdown()
		watcher.Stop()
		logger.Info("main", "Bye.")
	}
}

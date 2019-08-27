package main

import (
	"flag"
	"github.com/miekg/dns"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"udns/config"
	"udns/handler"
	"udns/logger"
	"udns/storage"
)

func main() {
	var cfgFile string

	flag.StringVar(&cfgFile, "cfg", "", "config file path, eg. /etc/udns.cfg.yml")
	flag.Parse()

	config.Init(cfgFile)
	logger.SetLogLevel(config.Cfg.LogLevel)
	storage.Init(config.Cfg.Source)

	addr := ":" + strconv.Itoa(config.Cfg.Port)
	logger.Infof("main", "Listen to %s(%s), parent DNS '%s'", addr, config.Cfg.Proto, config.Cfg.ParentDNS)
	server := &dns.Server{Addr: addr, Net: config.Cfg.Proto}
	for _, addr := range config.Cfg.Zones {
		dns.HandleFunc(addr, handler.OnSmartDNSRequest)
	}
	dns.HandleFunc(config.Cfg.MyIP, handler.OnMyIPRequest)
	dns.HandleFunc(".", handler.OnAnyDNSRequest)
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
			logger.Fatal("main", err)
		}
	}
}

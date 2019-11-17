package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"strings"
	"udns/logger"
	"udns/utils"
)

const (
	defaultCfgFile = "udns.cfg.yml"
)

var (
	Cfg = &config{}
)

type config struct {
	Proto     string   `yaml:"proto"`
	Port      int      `yaml:"port"`
	LogLevel  string   `yaml:"logLevel"`
	TTL       uint32   `yaml:"ttl"`
	Zones     []string `yaml:"zones"`
	ParentDNS string   `yaml:"parent"`
	MyIP      string   `yaml:"myip"`
	Source    string   `yaml:"source"`
}

func Init(cfgFile string) {
	if cfgFile == "" {
		logger.Warnf("config", "config file is not set, use %s", defaultCfgFile)
		cfgFile = defaultCfgFile
	}
	data, err := ioutil.ReadFile(cfgFile)
	if err != nil {
		logger.Fatal("config", err)
	}
	if err := yaml.Unmarshal(data, Cfg); err != nil {
		logger.Fatal("config", err)
	}

	logger.SetLogLevel(Cfg.LogLevel)

	if Cfg.ParentDNS == "auto" {
		logger.Info("config", "try to detect parent DNS")
		gw := utils.GetDefaultGateway()
		if gw == "" {
			logger.Fatal("config", "cannot detect default gateway address")
		}
		Cfg.ParentDNS = gw + ":53"
	} else if !strings.ContainsRune(Cfg.ParentDNS, ':') {
		logger.Info("config", "parent DNS has no port set, use :53")
		Cfg.ParentDNS += ":53"
	}
}

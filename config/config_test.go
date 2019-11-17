package config

import (
	"io/ioutil"
	"os"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

const (
	cfgFile = "test.cfg.yml"
)

func TestConfig(t *testing.T) {
	Convey("init config from file", t, func() {
		ioutil.WriteFile(cfgFile, []byte(`
proto: udp4
port: 53
logLevel: info
ttl: 5
zones:
  - my.internal.
  - you.external.
parent: auto
myip: my.ip
source: file://local.dns`), 0644)
		defer os.Remove(cfgFile)
		Init(cfgFile)
		So(Cfg.Proto, ShouldEqual, "udp4")
		So(Cfg.Port, ShouldEqual, 53)
		So(Cfg.LogLevel, ShouldEqual, "info")
		So(Cfg.Zones, ShouldResemble, []string{"my.internal.", "you.external."})
		So(Cfg.ParentDNS, ShouldNotEqual, "auto")
		So(Cfg.ParentDNS[len(Cfg.ParentDNS)-3:], ShouldEqual, ":53")
		So(Cfg.MyIP, ShouldEqual, "my.ip")
		So(Cfg.Source, ShouldEqual, "file://local.dns")
	})
}

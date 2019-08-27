package storage

import (
	"io/ioutil"
	"os"
	"testing"
	"udns/config"

	"github.com/miekg/dns"
	. "github.com/smartystreets/goconvey/convey"
)

func init() {
	config.Cfg.TTL = 5
}

func TestFileStorage(t *testing.T) {
	Convey("load dns records with empty file", t, func() {
		ioutil.WriteFile("test_dns.data", []byte(""), 0644)
		defer os.Remove("test_dns.data")
		defer Cache.Reset()
		newFileStorage("test_dns.data")
		So(len(Cache), ShouldEqual, 0)
	})

	Convey("load dns records with comment line", t, func() {
		ioutil.WriteFile("test_dns.data", []byte("# fqdn    ip     nettype"), 0644)
		defer os.Remove("test_dns.data")
		defer Cache.Reset()
		newFileStorage("test_dns.data")
		So(len(Cache), ShouldEqual, 0)
	})

	Convey("load dns records with empty line", t, func() {
		ioutil.WriteFile("test_dns.data", []byte("\n\n"), 0644)
		defer os.Remove("test_dns.data")
		defer Cache.Reset()
		newFileStorage("test_dns.data")
		So(len(Cache), ShouldEqual, 0)
	})

	Convey("load dns records with invalid lines", t, func() {
		ioutil.WriteFile("test_dns.data", []byte(`
# fqdn              ip        type
  yum.dx.corp       10.1.2.3  A
#yum.dx.corp        10.1.2.3  A
falcon.dx.corp      10.1.2.4  A
invalid.dx.corp    

invalid2.dx.corp    10.1.2.10  XXX
invalid3.dx.corp    10.1.2.10  A    invalid

pypi.dx.internal    10.1.2.5,10.1.2.6  A
        `), 0644)
		defer os.Remove("test_dns.data")
		defer Cache.Reset()
		ds := newFileStorage("test_dns.data")
		So(len(Cache), ShouldEqual, 1)
		So(len(Cache[dns.TypeA]), ShouldEqual, 3)
		rrs, err := ds.Query(dns.TypeA, "yum.dx.corp.")
		So(err, ShouldBeNil)
		So(len(rrs), ShouldEqual, 1)
		So(rrs[0].String(), ShouldEqual, "yum.dx.corp.\t5\tIN\tA\t10.1.2.3")
		rrs, err = ds.Query(dns.TypeA, "falcon.dx.corp.")
		So(err, ShouldBeNil)
		So(len(rrs), ShouldEqual, 1)
		So(rrs[0].String(), ShouldEqual, "falcon.dx.corp.\t5\tIN\tA\t10.1.2.4")
		rrs, err = ds.Query(dns.TypeA, "pypi.dx.internal.")
		So(err, ShouldBeNil)
		So(len(rrs), ShouldEqual, 2)
		So(rrs[0].String(), ShouldEqual, "pypi.dx.internal.\t5\tIN\tA\t10.1.2.5")
		So(rrs[1].String(), ShouldEqual, "pypi.dx.internal.\t5\tIN\tA\t10.1.2.6")
	})
}

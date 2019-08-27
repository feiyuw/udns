package storage

import (
	"bufio"
	"errors"
	"github.com/miekg/dns"
	"net"
	"os"
	"strings"
	"udns/config"
	"udns/logger"
)

func newFileStorage(filePath string) dnsStorage {
	stor := &fileStorage{filePath: filePath}
	stor.load()
	return stor
}

type fileStorage struct {
	filePath string
}

func (stor *fileStorage) Query(dnsType uint16, name string) ([]dns.RR, error) {
	return Cache.Get(dnsType, name)
}

func (ds *fileStorage) load() {
	file, err := os.Open(ds.filePath)
	if err != nil {
		logger.Fatalf("storage/file", "read data source error: %v", err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || line[0] == '#' {
			continue
		}
		// # fqdn       ip           type(*_)
		// a.dx.corp    10.1.2.17    A
		// b.dx.corp    10.1.2.18
		var rrs []dns.RR
		var err error
		var rrType uint16

		words := strings.Fields(line)
		fqdn := words[0]
		if fqdn[len(fqdn)-1] != '.' {
			fqdn += "."
		}

		switch len(words) {
		case 2: // type is A by default
			rrType = dns.TypeA
			rrs, err = recordToA(fqdn, words[1])
		case 3:
			switch words[2] {
			case "A":
				rrType = dns.TypeA
				rrs, err = recordToA(fqdn, words[1])
			default:
				logger.Errorf("storage/file", "unsupported type %s, in line %s", words[2], line)
				continue
			}
		default:
			logger.Errorf("storage/file", "invalid line %s", line)
			continue
		}

		if err != nil {
			logger.Errorf("storage/file", "parse record error: %v", err)
			continue
		}
		Cache.Set(rrType, fqdn, rrs)
	}

	if err := scanner.Err(); err != nil {
		return
	}
}

func recordToA(fqdn string, ipStr string) ([]dns.RR, error) {
	ips := strings.Split(ipStr, ",")
	rrs := make([]dns.RR, len(ips))
	for idx, ip := range ips {
		realIp := net.ParseIP(ip)
		if realIp == nil {
			return nil, errors.New("invalid ips format " + ipStr)
		}
		rrs[idx] = &dns.A{
			Hdr: dns.RR_Header{Name: fqdn, Rrtype: dns.TypeA, Class: dns.ClassINET, Ttl: config.Cfg.TTL},
			A:   realIp,
		}
	}

	return rrs, nil
}

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
	"udns/watcher"
)

func newFileStorage(filePath string) dnsStorage {
	stor := &fileStorage{filePath: filePath}
	if err := stor.load(); err != nil {
		logger.Fatal("storage/file", err)
	}
	watcher.On(filePath, stor.load)
	return stor
}

type fileStorage struct {
	filePath string
}

// Query used to fetch DNS record
func (stor *fileStorage) Query(dnsType uint16, name string) ([]dns.RR, error) {
	return Cache.Get(dnsType, name)
}

// Shutdown is used in server stop step
func (stor *fileStorage) Shutdown() {
	if err := watcher.Remove(stor.filePath); err != nil {
		logger.Error("storage/file", err)
	}
}

func (ds *fileStorage) load() error {
	logger.Info("storage/file", "load DNS records from file storage")
	file, err := os.Open(ds.filePath)
	if err != nil {
		return errors.New("read data source error: " + err.Error())
	}
	defer file.Close()

	newCache := NewCache()

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
		newCache.Set(rrType, fqdn, rrs)
	}

	UpdateCache(newCache)

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
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

package storage

import (
	"fmt"
	"github.com/miekg/dns"
)

var (
	Cache = cache{}
)

type cache map[uint16]map[string][]dns.RR

func (c cache) Get(dnsType uint16, record string) ([]dns.RR, error) {
	typeCache, exists := c[dnsType]
	if !exists {
		return nil, fmt.Errorf("%s 404 not found", record)
	}
	result, exists := typeCache[record]
	if !exists {
		return nil, fmt.Errorf("%s 404 not found", record)
	}
	return result, nil
}

func (c cache) Set(dnsType uint16, key string, value []dns.RR) {
	typeCache, exists := c[dnsType]
	if !exists {
		c[dnsType] = map[string][]dns.RR{}
		typeCache = c[dnsType]
	}

	typeCache[key] = value
}

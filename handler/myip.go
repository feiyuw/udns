package handler

import (
	"github.com/miekg/dns"
	"net"
	"udns/config"
	"udns/logger"
)

// OnMyIPRequest is used to handle smart dns request
func OnMyIPRequest(w dns.ResponseWriter, r *dns.Msg) {
	ip, ok := w.RemoteAddr().(*net.UDPAddr)
	if !ok {
		logger.Error("handler/myip", "failed to get remote addr as UDP!")
		return
	}
	fqdn := r.Question[0].Name
	logger.Debugf("handler/myip", "FQDN: %s, Remote IP: %s", fqdn, ip.IP)

	rr := &dns.A{
		Hdr: dns.RR_Header{
			Name:   fqdn,
			Rrtype: dns.TypeA,
			Class:  dns.ClassINET,
			Ttl:    config.Cfg.TTL,
		},
		A: ip.IP,
	}
	m := new(dns.Msg)
	m.Authoritative = true
	m.SetReply(r)
	m.Answer = []dns.RR{rr}
	w.WriteMsg(m)
}

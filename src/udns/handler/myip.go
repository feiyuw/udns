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
		logger.Error("dns", "failed to get remote addr as UDP!")
		return
	}
	name := r.Question[0].Name
	logger.Debugf("dns", "FQDN: %s, Remote IP: %s", name, ip.IP)

	rr := &dns.A{
		Hdr: dns.RR_Header{Name: name, Rrtype: dns.TypeA, Class: dns.ClassINET, Ttl: config.Cfg.TTL},
		A:   ip.IP,
	}
	m := &dns.Msg{Answer: []dns.RR{rr}}
	m.SetReply(r)
	w.WriteMsg(m)
}

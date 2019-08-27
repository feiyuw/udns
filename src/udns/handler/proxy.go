package handler

import (
	"github.com/miekg/dns"
	"udns/config"
	"udns/logger"
)

var (
	dnsClient = &dns.Client{Net: "udp"}
)

// OnAnyDNSRequest forward dns request to upstream
func OnAnyDNSRequest(w dns.ResponseWriter, r *dns.Msg) {
	logger.Debugf("handler/proxy", "Forward request %s to parent DNS", r.Question[0].Name)
	resp, _, err := dnsClient.Exchange(r, config.Cfg.ParentDNS)
	if err != nil {
		logger.Error("handler/proxy", "forward parent DNS", w.RemoteAddr(), r.Question[0].Name, err)
		dns.HandleFailed(w, r)
		return
	}
	if err = w.WriteMsg(resp); err != nil {
		logger.Error("handler/proxy", err)
	}
}

package handler

import (
	"github.com/miekg/dns"
	"udns/logger"
	"udns/storage"
)

// OnSmartDNSRequest is used to handle smart dns request
func OnSmartDNSRequest(w dns.ResponseWriter, r *dns.Msg) {
	// TODO: smart handler for incoming IP
	rrs, err := storage.Cache.Get(r.Question[0].Qtype, r.Question[0].Name)

	if err != nil {
		logger.Error("handler/smart", err)
		OnAnyDNSRequest(w, r) // forward to upstream request
		return
	}

	m := new(dns.Msg)
	m.SetReply(r)
	m.Authoritative = true
	m.Answer = rrs
	w.WriteMsg(m)
}
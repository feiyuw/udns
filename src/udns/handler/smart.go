package handler

import (
	"github.com/miekg/dns"
	"udns/logger"
	"udns/storage"
)

// OnSmartDNSRequest is used to handle smart dns request
func OnSmartDNSRequest(w dns.ResponseWriter, r *dns.Msg) {
	// TODO: smart handler for incoming IP
	rrs, err := storage.Query(r.Question[0].Qtype, r.Question[0].Name)

	if err != nil {
		logger.Error("handler/smart", err)
		OnAnyDNSRequest(w, r) // forward to upstream request
		return
	}

	m := new(dns.Msg)
	m.Authoritative = true
	m.SetReply(r)
	m.Answer = rrs
	w.WriteMsg(m)
}

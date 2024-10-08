package worker

import (
	"ihosts/global"
	"strings"
)

type Result struct {
	Domain string
	IP     []string
}

var (
	hostChanList = make(chan *global.Hosts, global.EnvConfig.GoroutineCount)
	chanCount    = 0
	controlChan  = make(chan struct{})
)

func (r *Result) GetDomain() string {
	return strings.Replace(r.Domain, global.IPAddressDomain, "", -1)
}

func (r *Result) GetIP() (ip string) {
	for _, ips := range r.IP {
		ip = ips
	}
	return
}

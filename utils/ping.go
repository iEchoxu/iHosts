package utils

import (
	"errors"
	"github.com/go-ping/ping"
	"time"
)

type PingOpts struct {
	Count         int
	Interval      time.Duration
	Timeout       time.Duration
	SetPrivileged bool
}

func NewPing(ip string, pingOpts *PingOpts) (*PingInfo, error) {
	var pingInfos *PingInfo

	pinger, err := ping.NewPinger(ip)
	if err != nil {
		return nil, errors.New("发起 ping 请求失败")
	}

	pinger.Count = pingOpts.Count
	pinger.Interval = pingOpts.Interval
	pinger.Timeout = pingOpts.Timeout            // 总耗时 3s
	pinger.SetPrivileged(pingOpts.SetPrivileged) // windows 上运行必须添加此行代码,Linux/Unix 上必须注释此行代码
	pinger.OnFinish = func(statistics *ping.Statistics) {
		pingInfos = &PingInfo{
			RttIP: &RttIP{
				IP:  statistics.Addr,
				Rtt: statistics.AvgRtt.Seconds(),
			},
			PacketLoss: statistics.PacketLoss,
		}
	}

	err = pinger.Run()
	if err != nil {
		return nil, errors.New("ping 运行出错")
	}

	return pingInfos, nil
}

func IsEnablePrivileged() bool {
	var isEnablePrivileged = false
	sysTypeRes := GetPlatform()
	if sysTypeRes == "windows" {
		isEnablePrivileged = true
	}
	return isEnablePrivileged
}

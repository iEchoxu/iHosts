package utils

import (
	"ihosts/global"
	"log"
	"math"
	"strings"
)

type RttIP struct {
	IP  string
	Rtt float64
}

type PingInfo struct {
	*RttIP
	PacketLoss float64
}

type DomainPingInfo struct {
	Domain string
	Ping   []*PingInfo
}

// GetDomain 去除 ipaddress.com 前缀获取真实域名
func (d *DomainPingInfo) GetDomain() string {
	//return strings.Split(d.domain, "/")[3]
	return strings.Replace(d.Domain, "https://sites.ipaddress.com/", "", -1)
}

// GetPingInfo 获取 ping 结果
func (d *DomainPingInfo) GetPingInfo() []*PingInfo {
	var pingInfoList []*PingInfo

	for _, pingRes := range d.Ping {
		pingInfoList = append(pingInfoList, &PingInfo{
			RttIP: &RttIP{
				IP:  pingRes.IP,
				Rtt: pingRes.Rtt,
			},
			PacketLoss: pingRes.PacketLoss,
		})
	}

	return pingInfoList
}

// PrintPingDetails 打印详细的 ping 过程
func (d *DomainPingInfo) PrintPingDetails() {
	if global.EnvConfig.PingLog {
		for _, pingInfos := range d.GetPingInfo() {
			log.Printf("%s 的 ping 结果：ip:%s rtt:%f 丢包率:%f", d.Domain, pingInfos.IP, pingInfos.Rtt, pingInfos.PacketLoss)
		}
	}
}

// GetPacketLossOfIPList 检测丢包率
func (d *DomainPingInfo) GetPacketLossOfIPList() []*RttIP {
	var NoPacketLossOfIPList, PacketLossOfIPList []*RttIP

	for _, pingRes := range d.Ping {
		if pingRes.PacketLoss == 0 {
			NoPacketLossOfIPList = append(NoPacketLossOfIPList, &RttIP{
				IP:  pingRes.IP,
				Rtt: pingRes.Rtt,
			})
		}

		// 取丢包率为 0-100 之间的数据，bug: 在丢包率全部为 100 时可能导致NoPacketLossOfIPList没有值
		// TODO 当丢包率全部为 100 时表示全部 ip 无法 ping 通,无 ip 可用
		// TODO 从丢包率最小的 ip 里选取 ip
		if pingRes.PacketLoss != 100 && pingRes.PacketLoss > 0 {
			PacketLossOfIPList = append(PacketLossOfIPList, &RttIP{
				IP:  pingRes.IP,
				Rtt: pingRes.Rtt,
			})
		}

	}

	// 当解析到的 ip 都出现丢包的时候，只能从 PacketLossOfIPList 中选择值
	if len(NoPacketLossOfIPList) == 0 {
		NoPacketLossOfIPList = PacketLossOfIPList
	}

	return NoPacketLossOfIPList
}

// GetMinimumRttValue 获取最小 rtt 值
func (d *DomainPingInfo) GetMinimumRttValue() float64 {
	var compareNumber float64 = 10
	NoPacketLossOfIPList := d.GetPacketLossOfIPList()
	for _, pingRes := range NoPacketLossOfIPList {
		if math.Min(pingRes.Rtt, compareNumber) == pingRes.Rtt {
			compareNumber = pingRes.Rtt
		}
	}

	return compareNumber
}

// GetMinimumRttOfIP 获取最小 rtt 值对应的 ip
func (d *DomainPingInfo) GetMinimumRttOfIP() string {
	minimumRtt := d.GetMinimumRttValue()
	var Accuracy = 0.000000000001 // 精准度
	var minimumRttOfIP string
	for _, pingRes := range d.Ping {
		// 查找列表中的值是否与上面获取到的最小 rtt 值相等
		if math.Abs(pingRes.Rtt-minimumRtt) < Accuracy {
			minimumRttOfIP = pingRes.IP
		}
	}

	return minimumRttOfIP
}

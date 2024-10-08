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
	return strings.Replace(d.Domain, global.IPAddressDomain, "", -1)
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

// GetMinimumPacketLossHandler 获取最小丢包率的 handler 函数,返回 丢包率为0的列表以及丢包率不为0 的列表
func (d *DomainPingInfo) GetMinimumPacketLossHandler() (NoPacketLossList, PacketLossList []*PingInfo) {
	PacketLossOfIPList := d.GetPingInfo()
	var FullPacketLossList []*PingInfo
	for _, pingRes := range PacketLossOfIPList {
		if pingRes.PacketLoss == 0 {
			NoPacketLossList = append(NoPacketLossList, &PingInfo{
				RttIP: &RttIP{
					IP:  pingRes.IP,
					Rtt: pingRes.Rtt,
				},
				PacketLoss: pingRes.PacketLoss,
			})
			continue
		}

		if pingRes.PacketLoss == 100 {
			FullPacketLossList = append(FullPacketLossList, &PingInfo{
				RttIP: &RttIP{
					IP:  pingRes.IP,
					Rtt: pingRes.Rtt,
				},
				PacketLoss: pingRes.PacketLoss,
			})
			continue
		}

		// 丢包率为 0-100 时
		if pingRes.PacketLoss > 0 && pingRes.PacketLoss != 100 {
			PacketLossList = append(PacketLossList, &PingInfo{
				RttIP: &RttIP{
					IP:  pingRes.IP,
					Rtt: pingRes.Rtt,
				},
				PacketLoss: pingRes.PacketLoss,
			})
		}
	}

	if len(FullPacketLossList) == len(d.Ping) {
		return nil, nil
	}

	return NoPacketLossList, PacketLossList
}

// GetMinimumPacketLossOfIPList 获取最小丢包率的 pingInfo,当最小丢包率对应的 ip 有多个时都要添加
func (d *DomainPingInfo) GetMinimumPacketLossOfIPList() (MinimumPacketLossList []*RttIP) {
	NoPacketLossList, PacketLossList := d.GetMinimumPacketLossHandler()

	// 当解析到的所有 ip 的丢包率都为 100% 时
	if NoPacketLossList == nil && PacketLossList == nil {
		log.Printf("%s 获取到的 ip 丢包率全部为百分之百,请稍后再试...", d.Domain)
		return nil
	}

	// 当丢包率为 0 的列表有数据时直接返回,最终 ip 从这个列表里选
	if len(NoPacketLossList) != 0 {
		for _, i := range NoPacketLossList {
			MinimumPacketLossList = append(MinimumPacketLossList, i.RttIP)
		}
		return
	}

	// 当丢包率为 0 的列表里没有数据, 最终 ip 从丢包率为 1-99.9999% 里选
	if len(NoPacketLossList) == 0 || len(PacketLossList) != 0 {
		var (
			compareNumber     float64 = 100
			firstValueOfRange PingInfo
			RttIPListTemp     []*PingInfo
		)

		for k, v := range PacketLossList {
			v := v
			//fmt.Printf("总列表：序号: %d,IP 为；%s,Rtt 为: %f,丢包率为: %f\n", k, v.IP, v.Rtt, v.PacketLoss)

			// 将序号为 1 的 pingInfo 存储下来提供给后面的 pingInfo 进行对比
			if k == 0 {
				firstValueOfRange = PingInfo{
					RttIP:      v.RttIP,
					PacketLoss: v.PacketLoss,
				}
				continue
			}

			// 与第一个里面的丢包率进行对比,如果小于它就添加进临时列表中
			if math.Min(v.PacketLoss, firstValueOfRange.PacketLoss) == v.PacketLoss {
				RttIPListTemp = append(RttIPListTemp, &PingInfo{
					RttIP:      v.RttIP,
					PacketLoss: v.PacketLoss,
				})
			}

			// 获取最小的丢包率的值
			if math.Min(v.PacketLoss, compareNumber) == v.PacketLoss {
				compareNumber = v.PacketLoss
			}

		}

		for _, i := range RttIPListTemp {
			// 当临时列表里的丢包率与最小丢包率的值相同时，添加进最终列表中
			if i.PacketLoss == compareNumber {
				MinimumPacketLossList = append(MinimumPacketLossList, &RttIP{
					IP:  i.IP,
					Rtt: i.Rtt,
				})
			}

			// 当第一个的丢包率与最小丢包率相同时添加进最终列表中
			if compareNumber == firstValueOfRange.PacketLoss {
				MinimumPacketLossList = append(MinimumPacketLossList, &RttIP{
					IP:  firstValueOfRange.IP,
					Rtt: firstValueOfRange.Rtt,
				})
			}
		}
	}

	return
}

// GetMinimumRttValue 从最小丢包率列表里获取最小 rtt 值
func (d *DomainPingInfo) GetMinimumRttValue() float64 {
	var compareNumber float64 = 10
	NoPacketLossOfIPList := d.GetMinimumPacketLossOfIPList()
	for _, pingRes := range NoPacketLossOfIPList {
		//fmt.Printf("最终获取到的: IP 为 %s,其 rtt 为: %f\n", pingRes.IP, pingRes.Rtt)
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

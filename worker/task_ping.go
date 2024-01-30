package worker

import (
	"ihosts/global"
	"ihosts/utils"
	"log"
	"sync"
	"time"
)

func Parse(result chan *Result) {
	var wg sync.WaitGroup

	defer close(hostChanList)

	for item := range result {
		item := item
		if len(item.IP) == 0 {
			log.Printf("%s 未获取到数据,请稍后再试...", item.Domain)
			continue
		}

		if len(item.IP) == 1 {
			domain := item.GetDomain()
			ip := item.GetIP()
			log.Printf("%s 解析到 1 个 ip, 其值为: %s", item.Domain, ip)
			hostChanList <- &global.Hosts{
				Domain: domain,
				IP:     ip,
			}
			chanCount++ //计数器加 1,用于统计成功获取到 ip 的 domain 数量
			continue
		}

		if len(item.IP) > 1 {
			wg.Add(1)
			go func() {
				defer wg.Done()
				doPing(item)
			}()
		}
	}

	controlChan <- struct{}{} // 用于控制刷新 DNS 缓存的协程在 Parse 协程后执行

	wg.Wait()
}

func doPing(item *Result) {
	var domainPingRes utils.DomainPingInfo
	var pingInfos []*utils.PingInfo

	pingOpts := &utils.PingOpts{
		Count:         10,
		Interval:      time.Millisecond * 100,
		Timeout:       time.Second * 3,
		SetPrivileged: utils.IsEnablePrivileged(),
	}

	for _, ip := range item.IP {
		pingInfoOnce, err := utils.NewPing(ip, pingOpts)
		if err != nil {
			log.Printf("ping 出现错误: %s", err)
		}

		pingInfos = append(pingInfos, pingInfoOnce)
	}

	domainPingRes = utils.DomainPingInfo{
		Domain: item.Domain,
		Ping:   pingInfos,
	}

	// 打印 ping 详细
	domainPingRes.PrintPingDetails()

	log.Printf("%s 解析到 %d 个 ip,最终 ip(取最小Avgrtt值) 为： %s", domainPingRes.Domain, len(domainPingRes.Ping), domainPingRes.GetMinimumRttOfIP())

	hostChanList <- &global.Hosts{
		Domain: domainPingRes.GetDomain(),
		IP:     domainPingRes.GetMinimumRttOfIP(),
	}
	chanCount++
}

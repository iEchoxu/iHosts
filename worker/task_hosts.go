package worker

import (
	"ihosts/utils"
	"log"
	"sync"
)

func hosts() {
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		platform, err := utils.GoosMap(utils.GetPlatform())
		if err != nil {
			log.Printf("获取系统平台出现错误:%s", err)
		}
		platform.UpdateHosts(hostChanList)
	}()

	wg.Wait()
}

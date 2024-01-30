package worker

import (
	"fmt"
	"ihosts/global"
	"log"
	"sync"
)

func Run() {
	var waitgroup sync.WaitGroup
	results := make(chan *Result, global.EnvConfig.GoroutineCount)
	log.Printf("正在解析网页获取最新的 IP，请稍后...")

	waitgroup.Add(3)

	go func() {
		defer waitgroup.Done()
		Fetch(results)
	}()

	go func() {
		defer waitgroup.Done()
		Parse(results)
	}()

	go func() {
		<-controlChan // 当该 channel 没有数据时会阻塞，可以让 hosts 协程运行在 Parse 协程后面

		// 频繁的爬取会触发 ipaddress.com 反爬机制，会导致有可能获取的数据全部为空
		// 当 hostChanList 为空值结束 hosts 协程，这样就不会运行刷新 DNS 的代码
		if len(hostChanList) == 0 {
			close(controlChan)
			waitgroup.Done()
		}

		hosts()
		waitgroup.Done()
	}()

	waitgroup.Wait()

	fmt.Printf("\nFetched %d/%d(total) sites \n", chanCount, global.NumberOfUrlsAfterDeDuplication)

}

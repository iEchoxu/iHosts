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
		defer waitgroup.Done()
		hosts()
	}()

	waitgroup.Wait()

	fmt.Printf("\nFetched %d/%d(total) sites \n", chanCount, global.NumberOfUrlsAfterDeDuplication)

}

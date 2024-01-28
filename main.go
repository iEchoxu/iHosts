package main

import (
	"fmt"
	"ihosts/utils"
	"ihosts/worker"
	"log"
	"time"
)

func main() {
	start := time.Now()
	if !utils.IsRoot() {
		log.Fatal("请用 root 账号执行此程序")
	}
	worker.Run()
	duration := time.Since(start)
	fmt.Printf("总耗时：%s\n\n", duration)
}

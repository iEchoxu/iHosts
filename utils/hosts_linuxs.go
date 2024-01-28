package utils

import (
	"fmt"
	"ihosts/global"
	"log"
	"os/exec"
)

type Linux struct {
	*BaseType
}

func NewLinux() *Linux {
	return &Linux{
		BaseType: &BaseType{
			HostsPath: "/etc/hosts",
			FlushCMD:  exec.Command("systemd-resolve", "--flush-caches"),
		},
	}
}

func (l *Linux) UpdateHosts(hostsChanList chan *global.Hosts) {
	writeBytes, err := l.HostsWrite(hostsChanList)
	if err != nil {
		log.Printf("写入 hosts 出现错误,报错为: %s", err)
	}

	log.Printf("操作系统为 %s，其 hosts 文件路径为： %s", GetPlatform(), l.HostsPath)
	log.Printf("正在更新系统 hosts 文件,请不要关闭此窗口...")

	if writeBytes != 0 {
		log.Println("已更新 hosts 文件，正在刷新 DNS 缓存...")
		output, err := l.FlushDNSCache()
		if err != nil {
			log.Printf("执行 DNS 缓存刷新命令失败:  %s", err)
		}
		fmt.Printf(output)
	}
}

package utils

import (
	"bufio"
	"errors"
	"fmt"
	"golang.org/x/text/encoding/simplifiedchinese"
	"ihosts/global"
	"io"
	"os"
	"os/exec"
	"runtime"
	"time"
)

var Goos = map[string]Hosts{
	"windows": NewWindows(),
	"linux":   NewLinux(),
	"darwin":  NewDarwin(),
}

type Hosts interface {
	UpdateHosts(hostsChanList chan *global.Hosts)
}

type BaseType struct {
	HostsPath string
	FlushCMD  *exec.Cmd
}

func (b *BaseType) HostsWrite(hostsChanList chan *global.Hosts) (int64, error) {
	// 备份文件
	hostsBack := b.HostsPath + ".bak"
	_, err := CopyFile(b.HostsPath, hostsBack)
	if err != nil {
		return 0, errors.New("hosts 文件备份失败")
	}

	fileHost, err := os.OpenFile(b.HostsPath, os.O_RDWR, os.ModePerm)
	if err != nil {
		return 0, errors.New("无法打开 hosts 文件")
	}

	tmpFile := b.HostsPath + "tmp"
	tmpContent, err := os.OpenFile(tmpFile, os.O_RDWR|os.O_CREATE|os.O_TRUNC, os.ModePerm)
	if err != nil {
		return 0, errors.New("无法打开 hosts 临时文件")
	}

	defer func(fileHost *os.File) {
		err := fileHost.Close()
		if err != nil {

		}
	}(fileHost)

	defer func(tmpContent *os.File) {
		err := tmpContent.Close()
		if err != nil {

		}
	}(tmpContent)

	// 取 hosts 中非配置文件中 start_urls 的内容并写入到临时文件中
	reader := bufio.NewReader(fileHost)
	for {
		line, err := reader.ReadString('\n')
		isDuplication := FileDeduplication(line)
		if err == io.EOF {
			// 因为是以换行符为分隔符，如果最后一行没有换行符
			// 那么返回 io.EOF 错误时，也是可能读到数据的，因此判断一下是否读到了数据
			if len(line) > 0 {
				_, err = tmpContent.WriteString(line)
			}
			break
		}
		if err != nil {
			return 0, errors.New("循环读 hosts 文件时出现错误")
		}

		if !isDuplication {
			_, err = tmpContent.WriteString(line)
		}

	}

	today := time.Now().Format("2006-01-02 15:04:05")
	_, err = tmpContent.WriteString("# Domain-IP 信息已于  " + string(today) + "  完成更新 \n")

	for hostsInfo := range hostsChanList {
		fmt.Printf("IP: %s,Domain: %s\n", hostsInfo.IP, hostsInfo.Domain)
		_, err = tmpContent.WriteString(hostsInfo.IP + "\t" + hostsInfo.Domain + "\n")
	}

	// 拷贝临时文件里的内容到真实 hosts 文件
	writeBytes, err := CopyFile(tmpFile, b.HostsPath)
	if err != nil {
		return 0, errors.New("无法复制内容到 hosts 文件中")
	}

	return writeBytes, nil
}

func (b *BaseType) FlushDNSCache() (string, error) {
	bytes, err := b.FlushCMD.Output()
	if err != nil {
		return "", errors.New("cmd.Run() failed with")
	}

	// 解决命令行输出中文乱码问题
	output, err := simplifiedchinese.GBK.NewDecoder().Bytes(bytes)
	if err != nil {
		return "", errors.New("utf-8 转 gbk 失败")
	}

	return string(output), nil
}

func GetPlatform() string {
	sysType := runtime.GOOS
	return sysType
}

func IsRoot() (isRoot bool) {
	isRoot = true
	goos := GetPlatform()
	if goos != "windows" && os.Geteuid() != 0 {
		isRoot = false
	}
	return isRoot
}

func GoosMap(goosName string) (Hosts, error) {
	if _, ok := Goos[goosName]; !ok {
		return nil, errors.New("暂不支持该平台")
	}
	return Goos[goosName], nil
}

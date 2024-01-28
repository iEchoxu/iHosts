package worker

import (
	"ihosts/global"
	"ihosts/utils"
	"io"
	"log"
	"regexp"
	"sync"
	"time"
)

func Fetch(result chan *Result) {
	var wg sync.WaitGroup
	defer close(result)

	for _, url := range genUrls() {
		wg.Add(1)
		go func(url string) {
			defer wg.Done()
			res := response(url)
			result <- res
		}(url)
	}

	wg.Wait()
}

// 获取去重后的 url 列表
func genUrls() []string {
	var domainList []string
	domainMap := make(map[string]struct{})
	for _, domain := range global.EnvConfig.StartUrls {
		domain := domain
		if _, ok := domainMap[domain]; !ok {
			url := "https://sites.ipaddress.com/" + domain
			domainList = append(domainList, url)
			domainMap[domain] = struct{}{}
			global.NumberOfUrlsAfterDeDuplication++
		}
	}

	return domainList
}

func response(url string) *Result {
	HTMLContent := request(url)
	re := regexp.MustCompile(`href="https://www.ipaddress.com/ipv4/(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)(\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)){3}">`)
	matchContent := re.FindAllString(HTMLContent, -1)
	var ipaddr []string
	for _, content := range matchContent {
		getIp := content[37 : len(content)-2]
		ipaddr = append(ipaddr, getIp)
	}

	results := new(Result)
	results.IP = ipaddr
	results.Domain = url

	return results
}

func request(url string) string {
	requestIns := &utils.RequestOpts{
		UserAgent:      utils.RandomUserAgent(),
		Referer:        utils.RandomReferer(),
		AcceptLanguage: "zh-CN,zh;q=0.9",
		Timeout:        60 * time.Second,
	}

	resp, err := utils.HttpDo(requestIns, url)
	if err != nil || resp.StatusCode != 200 {
		log.Printf("连接 ipaddress.com 异常：%s", err)
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("读取 ipaddress.com 返回的 HTML 数据失败：%s", err)
	}

	//utils.PrintRequestInfo(resp)  //验证随机 user-agent 是否生效

	return string(body)
}

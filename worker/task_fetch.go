package worker

import (
	"github.com/antchfx/htmlquery"
	"ihosts/global"
	"ihosts/utils"
	"io"
	"log"
	"regexp"
	"strings"
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
			url := global.IPAddressDomain + domain
			domainList = append(domainList, url)
			domainMap[domain] = struct{}{}
			global.NumberOfUrlsAfterDeDuplication++
		}
	}

	return domainList
}

func response(url string) *Result {
	HTMLContent := request(url)
	IPContent, err := htmlquery.Parse(strings.NewReader(HTMLContent))
	if err != nil {
		log.Println("htmlQuery 解析 HTML 内容出错:", err)
	}

	IPNode := htmlquery.Find(IPContent, `//*[@id="tabpanel-dns-a"]/pre/a/text()`)
	re := regexp.MustCompile(`(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)(\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)){3}`)
	var IPList []string
	for i := range IPNode {
		ip := re.FindString(IPNode[i].Data)
		IPList = append(IPList, ip)
	}

	results := new(Result)
	results.IP = IPList
	results.Domain = url

	return results
}

func responseBack(url string) *Result {
	HTMLContent := request(url)
	IPContent, err := htmlquery.Parse(strings.NewReader(HTMLContent))
	if err != nil {
		log.Println("htmlQuery 解析 HTML 内容出错:", err)
	}

	IPNode := htmlquery.FindOne(IPContent, `//*[@id="main-wrapper"]//*[@class='summary']/p/text()`)
	re := regexp.MustCompile(`(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)(\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)){3}`)
	IPS := re.FindAllString(IPNode.Data, -1)

	results := new(Result)
	results.IP = IPS
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

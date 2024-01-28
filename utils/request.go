package utils

import (
	"errors"
	"log"
	"math/rand"
	"net/http"
	"time"
)

type RequestOpts struct {
	UserAgent      string
	Referer        string
	AcceptLanguage string
	Timeout        time.Duration
}

func HttpDo(opts *RequestOpts, url string) (*http.Response, error) {
	client := &http.Client{
		Timeout: opts.Timeout,
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, errors.New("无法发送网络请求，请检查你的网络")
	}

	req.Header.Set("User-Agent", opts.UserAgent)
	req.Header.Add("Referer", opts.Referer)                // 必须添加此行，不然会报错 403
	req.Header.Add("accept-language", opts.AcceptLanguage) // 必须添加此行，不然会报错 403

	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.New("发送请求异常")
	}

	return resp, nil
}

func RandomReferer() string {
	var referer = []string{
		"https://github.com",
		"https://www.baidu.com",
		"https://www.google.com",
	}
	rand.NewSource(time.Now().UnixNano())
	refererRand := referer[rand.Intn(len(referer))]

	return refererRand
}

func RandomUserAgent() string {
	var userAgent = []string{
		"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4896.127 Safari/537.36",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36 Edg/120.0.0.0",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
		"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
		"Mozilla/5.0 (X11; Linux x86_64; rv:121.0) Gecko/20100101 Firefox/121.0",
	}

	rand.NewSource(time.Now().UnixNano())
	useAgentRand := userAgent[rand.Intn(len(userAgent))]

	return useAgentRand
}

// PrintRequestInfo 检测随机 userAgent 是否生效
func PrintRequestInfo(resp *http.Response) {
	log.Printf("referer is: %s", resp.Request.Referer())
	log.Printf("user-agent is: %s", resp.Request.UserAgent())
}

package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

func main() {
	start := time.Now()
	for i := 0; i < 10; i++ {
		go parseUrls("https://movie.douban.com/top250?start=" + strconv.Itoa(25*i))
	}
	elapsed := time.Since(start)
	fmt.Printf("took %v", elapsed)
}

func fetch(cusURL string) string {
	// 设置代理
	proxyUrl := "http://127.0.0.1:10080"
	proxy, _ := url.Parse(proxyUrl)
	transport := &http.Transport{Proxy: http.ProxyURL(proxy)}
	client := &http.Client{Transport: transport}

	// 发起请求
	req, _ := http.NewRequest("GET", cusURL, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("http get error", err)
		return ""
	}
	defer resp.Body.Close() // 关闭连接

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("read error")
		return ""
	}
	return string(body)
}

func parseUrls(url string) {
	fmt.Println(url)
	fetch(url)
	//body = strings.Replace(body, "\n", "", -1)
	//rp := regexp.MustCompile(`<div class="hd">(.*?)</div>`)
	//titleRe := regexp.MustCompile(`<span class="title">(.*?)</span>`)
	//idRe := regexp.MustCompile(`<a href="https://movie.douban.com/subject/(\d+)/"`)
	//items := rp.FindAllStringSubmatch(body, -1)
	//for _, item := range items {
	//	fmt.Println(idRe.FindStringSubmatch(item[1])[1],
	//		titleRe.FindStringSubmatch(item[1])[1])
	//}
}

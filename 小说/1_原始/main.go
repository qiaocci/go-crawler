package main

import (
	"fmt"
	"github.com/antchfx/htmlquery"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

const BaseURL = "https://top.chinaz.com"

func main() {
	allUrls := prepareUrls()
	start := time.Now()
	for _, url := range allUrls {
		parse(url)
		time.Sleep(1 * time.Second)
	}
	elapsed := time.Since(start)
	fmt.Printf("took %v", elapsed)
}

func prepareUrls() []string {
	allUrls := []string{
		"https://top.chinaz.com/hangye/index_yule_xiaoshuo.html",
	}
	for i := 2; i < 40; i++ {
		allUrls = append(allUrls, fmt.Sprintf("https://top.chinaz.com/hangye/index_yule_xiaoshuo_%v.html", i))
	}
	return allUrls
}

func fetch(cusURL string) string {
	client := &http.Client{}

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

func parse(url string) {
	body := fetch(url)
	root, _ := htmlquery.Parse(strings.NewReader(body))
	tr := htmlquery.Find(root, "//*[@id=\"content\"]/div[3]/div[3]/div[2]/ul/li/div/h3/a")

	for _, row := range tr {
		title := htmlquery.InnerText(row)
		titleURL := BaseURL + htmlquery.SelectAttr(row, "href")
		fmt.Println(title, titleURL)
	}
}

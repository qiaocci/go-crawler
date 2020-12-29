package main

import (
	"fmt"
	"github.com/gocolly/colly"
)

func main() {
	c := colly.NewCollector()

	// 如果收到的响应内容是HTML调用它
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		fmt.Printf("解析colly官网, 找到网址%v\n", e.Attr("href"))
		//e.Request.Visit(e.Attr("href"))
	})

	// 请求前
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	c.Visit("http://go-colly.org/")
}

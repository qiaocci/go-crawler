package main

import (
	"fmt"
	"github.com/antchfx/htmlquery"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

type NovelSite struct {
	gorm.Model
	Title string `gorm:"type:varchar(100);"`
	Url   string `gorm:"type:varchar(100);"`
}

const BaseURL = "https://top.chinaz.com"

func main() {
	start := time.Now()

	ch := make(chan bool)

	// 建立连接
	dsn := "root:123@tcp(127.0.0.1:3306)/test?charset=utf8mb4&parseTime=true&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	checkError(err)

	db.Migrator().DropTable(&NovelSite{})
	db.Set("grom:table_options", "ENGINE=InnoDB").AutoMigrate(&NovelSite{})

	allUrls := prepareUrls()
	for _, url := range allUrls {
		go parse(url, ch, db)
	}
	for i := 0; i < 39; i++ {
		<-ch
	}
	elapsed := time.Since(start)
	fmt.Printf("took %v", elapsed)
}
func parse(url string, ch chan bool, db *gorm.DB) {
	body := fetch(url)
	//fmt.Println(body)
	root, _ := htmlquery.Parse(strings.NewReader(body))
	tr := htmlquery.Find(root, "//*[@id=\"content\"]/div[3]/div[3]/div[2]/ul/li/div/h3/a")

	for _, row := range tr {
		title := htmlquery.InnerText(row)
		titleURL := BaseURL + htmlquery.SelectAttr(row, "href")
		fmt.Println(title, titleURL)

		movie := NovelSite{
			Title: title,
			Url:   titleURL,
		}
		db.Create(&movie)

	}
	ch <- true
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

func checkError(err error) {
	if err != nil {
		fmt.Println("异常啦", err)
		panic(err)
	}
}

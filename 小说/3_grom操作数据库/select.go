package main

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
)

type NovelSite struct {
	gorm.Model
	Title string `gorm:"type:varchar(100);"`
	Url   string `gorm:"type:varchar(100);"`
}

func main() {
	dsn := "root:123@tcp(127.0.0.1:3306)/test?charset=utf8mb4&parseTime=true&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println("出错啦", err)
	}

	var novelSite NovelSite
	//var novelSites []NovelSite
	db.First(&novelSite, 1) // 查询id=1的记录，只支持主键
	log.Println(novelSite)
	//log.Println(novelSite.Title, novelSite.Url, novelSite.CreatedAt)
	//
	//db.Order("id desc").Limit(3).Find(&novelSites)
	//log.Println(novelSites)
	//log.Println(novelSites[0].ID, novelSites[0].Title)
	//
	//db.Order("id desc").Limit(3).Offset(1).Find(&novelSites)
	//log.Println(novelSites[0].ID, novelSites[0].Title)

	db.Select("title").Where("id=?", 3).Find(&novelSite)
	log.Printf("id=%v, title=%v, url=%v", novelSite.ID, novelSite.Title, novelSite.Url)
}

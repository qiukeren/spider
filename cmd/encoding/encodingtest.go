package main

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/qiukeren/spider/model"
	. "github.com/qiukeren/spider/utils"

	"log"
)

var db *gorm.DB

func main() {
	var err error

	db, err = gorm.Open("sqlite3", "spider.db")
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()
	log.SetFlags(log.Lshortfile | log.Ltime | log.Ldate)
	var contentStruct model.Content
	db.Where("url <> ?", "").Where("status = ?", 200).First(&contentStruct)
	content := contentStruct.Content
	log.Println(EncodingTest(&content))
}

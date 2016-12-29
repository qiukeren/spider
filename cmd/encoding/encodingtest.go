package main

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/qiukeren/spider/model"
	. "github.com/qiukeren/spider/utils"
	"gopkg.in/iconv.v1"

	"log"
)

var db *gorm.DB

func main() {
	var err error

	db, err = gorm.Open("sqlite3", "spider.db")
	if err != nil {
		panic("failed to connect database")
	}
	P("a", "b")
	defer db.Close()
	log.SetFlags(log.Lshortfile | log.Ltime | log.Ldate)
	var contentStruct model.Content
	db.Where("url <> ?", "").Where("status = ?", 200).First(&contentStruct)
	content := contentStruct.Content
	encoding, err := EncodingTest(&content)
	if err != nil {
		return
	}
	cd, err := iconv.Open("utf-8", encoding) // convert utf-8 to gbk
	if err != nil {
		log.Println("iconv.Open failed!")
		return
	}
	defer cd.Close()

	out := make([]byte, len(content)*2)
	_, _, _ = cd.Conv(content, out)

	//gbk := cd.ConvString(string(content))
	log.Println(string(content))
	log.Println("================================")
	log.Println(string(out))

}

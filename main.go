package main

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/sqlite"

	"github.com/qiukeren/spider/model"
	. "github.com/qiukeren/spider/utils"

	"bytes"
	"log"
	"sync"
	"time"
)

var db *gorm.DB
var wg sync.WaitGroup
var lock sync.Mutex

func init() {
	var err error

	// db, err = gorm.Open("sqlite3", "spider.db")
	db, err = gorm.Open("mysql", "spider:spider_pass@/spider?charset=utf8&parseTime=True&loc=Local")

	if err != nil {
		panic("failed to connect database")
	}
	//defer db.Close()

	db.AutoMigrate(&model.Site{})
	db.AutoMigrate(&model.Url{})
	db.AutoMigrate(&model.Content{})
}

func main() {
	log.SetFlags(log.Lshortfile | log.Ltime | log.Ldate)
	var err error
	//db, err = gorm.Open("sqlite3", "spider.db")
	db.LogMode(true)
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()

	array := []string{
		"http://www.duwenzhang.com/wenzhang/shenghuosuibi/20140520/291739.html",
		"http://www.oschina.net/news/80475/bfs-0-5-0",
		"https://my.oschina.net/lujianing/blog/787745",
		"http://coolshell.cn/articles/17583.html",
		"http://www.mike.org.cn/articles/some-classic-quotations-1-2/",
		"https://segmentfault.com/",
	}
	for _, v := range array {
		wg.Add(1)
		go GoSpide(v)
	}
	wg.Wait()

}

func GoSpide(url1 string) {
	a, err := StoreGetSite(url1)
	P("title", a)
	P("title", err)
	if err != nil {
		log.Println(err)
		return
	}
	SpidePage(a, url1)
	wg.Done()
}

func StoreGetSite(randomUrl string) (*model.Site, error) {
	urlStruct, err := ParseUrl(randomUrl)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	siteName := urlStruct.Host

	count := 0

	db.Model(&model.Site{}).Where("name = ?", siteName).Count(&count)

	if count == 0 {
		newSite := model.Site{Name: siteName, Url: siteName, Protocol: urlStruct.Scheme}
		db.Create(&newSite)
	}

	var siteStruct model.Site
	db.Where("name = ?", siteName).First(&siteStruct)
	return &siteStruct, nil
}

var map2 map[string]bool

func SpidePage(siteStruct *model.Site, url1 string) {

	if map2 == nil {
		map2 = make(map[string]bool)
	}

	lock.Lock()
	if _, b := map2[url1]; b {
		lock.Unlock()
		return
	}
	map2[url1] = true
	lock.Unlock()

	content, err := Get(url1)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println("spidering " + url1)
	StoreContent(siteStruct, url1, content)

	reader := bytes.NewReader(content)
	doc, err := goquery.NewDocumentFromReader(reader)

	if err != nil {
		log.Println(err)
	}
	site := siteStruct.Url

	doc.Find("a").Each(func(i int, contentSelection *goquery.Selection) {
		urlTemp, boolean := contentSelection.Attr("href")
		a, err := FormatUrl(urlTemp, site)

		if err != nil {
			log.Println(err)
			return
		}

		if boolean && IsCurrentSite(a, site, siteStruct.Protocol) {

			// _, err := Get(a)
			// if err != nil {
			// 	log.Println(err)
			// 	return
			// }
			//_ = content
			//log.Println("spidering " + a)
			//StorePage(siteStruct, urlStructTemp, content)

			//StoreContentUrl(siteStruct, a)

			time.Sleep(time.Millisecond * 10)
			SpidePage(siteStruct, a)
		} else {
			//log.Printf("none current url " + site + " " + a)
		}
	})

}

func StoreContent(siteStruct *model.Site, url1 string, content []byte) {

	count := 0

	db.Model(&model.Content{}).Where("url = ?", url1).Count(&count)
	encoding, _ := EncodingTest(&content)
	// p("title", count)
	if count == 0 {
		newContent := model.Content{
			Url:      url1,
			SiteId:   siteStruct.ID,
			Status:   200,
			Code:     200,
			Content:  content,
			Encoding: encoding,
		}
		db.Create(&newContent)
	} else {
		db.Model(model.Content{}).Where("url = ?", url1).Update(
			map[string]interface{}{
				"content":  content,
				"encoding": encoding,
			},
		)
	}
}

func StoreContentUrl(siteStruct *model.Site, url1 string) {

	count := 0

	db.Model(&model.Content{}).Where("url = ?", url1).Count(&count)
	// p("title", count)
	if count == 0 {
		newContent := model.Content{Url: url1, SiteId: siteStruct.ID, Status: 100, Code: 100}
		db.Create(&newContent)
	}
}

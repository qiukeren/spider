package main

import (
	"github.com/MouseSun/goprint"
	"github.com/PuerkitoBio/goquery"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/qiukeren/spider/model"
	. "github.com/qiukeren/spider/utils"

	"bytes"
	"log"
	"net/url"
	"time"
)

func p(title string, c interface{}) {
	goprint.P(title, c)
}

var db *gorm.DB

func init() {
	var err error

	db, err = gorm.Open("sqlite3", "spider.db")
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()

	db.AutoMigrate(&model.Site{})
	db.AutoMigrate(&model.Url{})
	db.AutoMigrate(&model.Content{})
}

func main() {
	log.SetFlags(log.Lshortfile | log.Ltime | log.Ldate)
	var err error
	db, err = gorm.Open("sqlite3", "spider.db")
	// db.LogMode(true)
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()

	a, err := StoreGetSite("http://www.duwenzhang.com/wenzhang/shenghuosuibi/20140520/291739.html")
	p("title", a)
	p("title", err)
	if err != nil {
		log.Println(err)
		return
	}
	SpidePage(a, "http://www.duwenzhang.com/wenzhang/shenghuosuibi/20140520/291739.html")
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

	if _, b := map2[url1]; b {
		return
	}

	content, err := Get(url1)
	if err != nil {
		log.Println(err)
		return
	}

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

		if boolean && IsCurrentSite(a, site) {

			// _, err := Get(a)
			// if err != nil {
			// 	log.Println(err)
			// 	return
			// }
			//_ = content
			log.Println("spidering " + a)
			//StorePage(siteStruct, urlStructTemp, content)
			StoreContentUrl(siteStruct, a)

			map2[url1] = true
			time.Sleep(time.Millisecond * 100)
			SpidePage(siteStruct, a)
		} else {
			log.Printf("none current url " + site + " " + a)
		}
	})

}

func StoreContentUrl(siteStruct *model.Site, url1 string) {

	count := 0
	var contentStruct model.Content
	db.Where("url = ?", url1).First(&contentStruct)

	db.Model(&model.Content{}).Where("url = ?", url1).Count(&count)
	// p("title", count)
	if count == 0 {
		newContent := model.Content{Url: url1, SiteId: siteStruct.ID, Status: 100, Code: 100}
		db.Create(&newContent)
	}
}

func StoreContent(siteStruct *model.Site, urlStruct *url.URL, content []byte) {
	url1 := urlStruct.String()

	count := 0
	var contentStruct model.Content
	db.Where("url = ?", url1).First(&contentStruct)

	db.Model(&model.Content{}).Where("url = ?", url1).Count(&count)
	// p("title", count)
	if count == 0 {
		newContent := model.Content{Url: url1, SiteId: siteStruct.ID, Content: content, Status: 200, Code: 200}
		db.Create(&newContent)
	}
}
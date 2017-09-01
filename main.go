package main

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/beeker1121/goque"
	"github.com/qiukeren/go-utils/common"
	. "github.com/qiukeren/go-utils/enca"
	iconv "github.com/qiukeren/go-utils/iconv"
	"github.com/remeh/sizedwaitgroup"
	"github.com/syndtr/goleveldb/leveldb"

	"bytes"
	"encoding/json"
	"errors"
	"log"
	"net/url"
	"strings"
	"time"
)

var q *goque.Queue
var db *leveldb.DB

type urlS struct {
	Url  string
	UrlP string
}

func main() {

	log.SetFlags(log.Lshortfile | log.Ltime | log.Ldate)
	var err error
	db, err = leveldb.OpenFile("./db", nil)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	q, err = goque.OpenQueue("./queue")
	if err != nil {
		panic(err)
	}
	defer q.Close()

	array := []string{
		"http://www.duwenzhang.com/wenzhang/shenghuosuibi/20140520/291739.html",
		"http://www.oschina.net/news/80475/bfs-0-5-0",
		"https://my.oschina.net/lujianing/blog/787745",
		"http://www.01happy.com/golang-tcp-socket-adhere/",
		"http://blog.studygolang.com/tag/golang_pkg/",
		"http://www.mike.org.cn/articles/some-classic-quotations-1-2/",
		"https://segmentfault.com/",
		"http://stackoverflow.com/questions/2635058/ibatis-get-executed-sql",
		"https://www.zhihu.com/question/27720523",
		"http://blog.dataman-inc.com/114-shurenyun-huodong/",
		"http://www.soomal.com/doc/10100005237.htm",
		"http://www.ruanyifeng.com/blog/2016/05/react_router.html?utm_source=tool.lu",
		"https://book.douban.com/review/5428330/",
		"http://www.l99.com/media_index.action",
		"http://www.dapenti.com/blog/more.asp?name=tupian&id=68524",
		"http://www.guoxue.com/?category_name=study&paged=3",
//		"http://www.nowamagic.net/librarys/veda/detail/2299",
		"http://blog.jobbole.com/811/",
		"http://blog.csdn.net/v_july_v/article/details/7382693",
		"http://www.cnblogs.com/yuuyuu/p/5180827.html",
		"http://www.ibm.com/developerworks/cn/linux/l-vim-script-2/",
		"http://tieba.baidu.com/p/2315349954",
		"https://zhuanlan.zhihu.com/",
		"http://www.socialbeta.com/",
		"https://www.zhihu.com/collection/43457862",
		"http://xh.5156edu.com/page/z5208m6176j18515.html",
		"http://g.yeeyan.org/",
		"http://www.dapenti.com/blog/more.asp?name=tupian&id=68524",
		"http://www.soomal.com/doc/10100005237.htm",
		"http://www.w2bc.com/Article/8951",
		"http://www.jianshu.com/p/941bfaf13be1",
		"http://www.juzimi.com/ju/144565",
		"http://limlee.blog.51cto.com/6717616/1223749",
		"https://www.zybuluo.com/Gestapo/note/32082",
		"http://rfyiamcool.blog.51cto.com/",
		"http://linux.chinaunix.net/techdoc/develop/2007/03/11/952015.shtml",
		"http://blog.chinaunix.net/uid-433083-id-3173212.html",
		"http://www.saltstack.cn/kb/salt-raet-01/#salt-raet-01",
		"https://book.douban.com/review/5428330/",
		"https://www.douban.com/photos/album/1651518758/",
		"http://www.postgres.cn/news/viewone/1/248",
		"http://www.bilibili.com/video/av3282944/",
		"http://bbs.paidai.com/topic/448017",
		"http://highscalability.com/blog/2013/5/13/the-secret-to-10-million-concurrent-connections-the-kernel-i.html",
		"http://www.qingshi.net/love/mingju/",
	}

	for _, v := range array {
		data, _ := json.Marshal(urlS{Url: v, UrlP: v})
		q.Enqueue(data)
	}
	GoQueue()
}
func GoQueue() {
	swg := sizedwaitgroup.New(200)
	for i := 0; ; i++ {
		swg.Add()
		go func(i int) {
			defer swg.Done()
			url1, err := PopQueue()
			if err != nil {
				time.Sleep(1 * time.Second)
				return
			}
			SpidePage(url1)
		}(i)
	}
	swg.Wait()

}

func SpidePage(url1 *urlS) {
	content, err := common.Get(url1.Url)
	if err != nil {
		log.Println("spder error: ", url1.Url, err)
		return
	}
	StoreContent(url1, content)

	reader := bytes.NewReader(content)
	doc, err := goquery.NewDocumentFromReader(reader)

	if err != nil {
		log.Println(err)
		return
	}
	log.Println("range page: ", url1.Url)
	defer log.Println("range done: ", url1.Url)
	doc.Find("a").Each(func(i int, contentSelection *goquery.Selection) {
		urlTemp, boolean := contentSelection.Attr("href")
		urls := &urlS{}
		urls.UrlP = url1.Url
		urls.Url = urlTemp
		urls, err := FormatUrl(urls)

		if err != nil {
			return
		}
		if boolean && IsCurrentSite(urls) && !Have(urls) {
			PushQueue(urls)
		}
	})
}

func Have(urls *urlS) bool {
	_, err := db.Get([]byte(CONTENT+urls.Url), nil)
	return err == nil
}

func IsCurrentSite(url1 *urlS) bool {
	if url1.Url == "" {
		return false
	}
	url11, err := url.Parse(url1.Url)
	if err != nil {
		return false
	}
	url22, err := url.Parse(url1.UrlP)
	if err != nil {
		return false
	}
	if url11.Scheme != "http" &&
		url11.Scheme != "https" &&
		url11.Scheme != "" {
		return false
	}
	if strings.HasPrefix(strings.TrimSpace(url1.Url), "javascript:") {
		return false
	}
	if url11.Host == "" {
		panic("不应该进入这段逻辑")
	}
	if url11.Host != "" {
		if url11.Host == url22.Host {
			return true
		}
	}
	return false

}

func FormatUrl(url1 *urlS) (url2 *urlS, err error) {
	url2 = &urlS{}
	url11, err := url.Parse(url1.Url)
	if err != nil {
		return nil, err
	}
	url22, err := url.Parse(url1.UrlP)
	if err != nil {
		return
	}
	if url11.Host == "" {
		if url22.Host == "" {
			err = errors.New("unknown host: " + url1.Url + "" + url1.UrlP)
			return
		}
		url11.Host = url22.Host
	}
	if url11.Scheme == "" {
		if url22.Scheme == "" {
			url11.Scheme = "http"
		} else {
			url11.Scheme = url22.Scheme
		}
	}
	if url11.Path == "" {
		url11.Path = "/"
	}
	url2.Url = url11.String()
	url2.UrlP = url1.UrlP
	return url2, nil
}

var default_encoding = "UTF-8"
var CONTENT = "CONTENT_"
var CONTENTUTF8 = "CONTENTUTF8_"
var ENCODING = "ENCODING_"
var QUEUE = "QUEUE_"

func StoreContent(url1 *urlS, content []byte) {
	encoding, _ := EncodingTest(&content)
	db.Put([]byte(ENCODING+url1.Url), []byte(encoding), nil)
	if encoding != default_encoding {
		data, _ := iconv.ConvString("UTF-8", "GBK", "123123123123啦啦啦")
		db.Put([]byte(CONTENTUTF8+url1.Url), []byte(data), nil)
		db.Put([]byte(CONTENT+url1.Url), content, nil)
	} else {
		db.Put([]byte(CONTENTUTF8+url1.Url), content, nil)
	}
}

func PushQueue(urls *urlS) {
	data, _ := json.Marshal(urls)
	_, err := db.Get([]byte(QUEUE+urls.Url), nil)
	if err == nil {
		return
	}
	log.Println("push queue: ", urls.Url)
	db.Put([]byte(QUEUE+urls.Url), []byte("1"), nil)
	q.Enqueue(data)
}

func PopQueue() (url1 *urlS, err error) {
	url1 = &urlS{}
	item, err := q.Dequeue()
	if err != nil {
		return
	}
	err = json.Unmarshal(item.Value, url1)
	log.Println("pop queue: ", url1.Url)
	return
}


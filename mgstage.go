package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/publicsuffix"
	"io/ioutil"
	"log"
	"net/http/cookiejar"
	"strings"

	//"github.com/jackdanger/collectlinks"
	"net/http"
	"net/url"
	"time"
)

var visited = make(map[string]bool)

func getMgstageToday() string{
	yyyymmdd := getTimeYYYYMMDD()
	url := "https://www.mgstage.com/search/cSearch.php?sale_start_range=%v-%v&type=top"
	return fmt.Sprintf(url,yyyymmdd,yyyymmdd)
}
func getTimeYYYYMMDD() string{
	t := time.Now()
	return fmt.Sprintf("%d-%02d-%02d",
		t.Year(), t.Month(), t.Day())
}
func getTimeYYYYMMDDHHMMSS() string{
	t := time.Now()
	return fmt.Sprintf("%d-%02d-%02dT %02d:%02d:%02d",
		t.Year(), t.Month(), t.Day(),
		t.Hour(), t.Minute(), t.Second())
}

func main() {
	fmt.Println("start mgstage")
	url := getMgstageToday()
	fmt.Println(url)

	queue := make(chan string, )
	go func() {
		queue <- url
	}()
	for uri := range queue {
		download(uri, queue)
	}
}

func download(url string, queue chan string) {
	visited[url] = true
	timeout := time.Duration(5 * time.Second)

	cookieJar, err := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	if err != nil {
		panic(err)
	}
	client := &http.Client{
		Jar: cookieJar,
		Timeout:timeout,
	}
	req, _ := http.NewRequest("GET", url, nil)
	// 自定义Header
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.4280.67 Safari/537.36")
	req.Header.Set("Cookie","_ga=GA1.2.1493108184.1606402051; __ulfpc=202011262247310860; _gid=GA1.2.930588390.1606921202; PHPSESSID=le4h3b9t3n5mol6osdoia065q3; bWdzdGFnZS5jb20%3D-_lr_uf_-r2icil=54562327-4eed-4df9-b0cb-1fa8d9146cfd; displayed_update_notice_flg=0; adc=1; _gat_UA-58252858-1=1; _gat_UA-158726521-1=1; uuid=EIQBKROJIH83HBAQ5MWK4X4FKA; aluid=d8f58764-e5cc-4cae-b0c0-6a1d2480221d; bWdzdGFnZS5jb20%3D-_lr_tabs_-r2icil%2Fmgs={%22sessionID%22:0%2C%22recordingID%22:%224-6982ac76-eaec-4dd5-aec6-083f41c8ec3e%22%2C%22lastActivity%22:1607008145749}; bWdzdGFnZS5jb20%3D-_lr_hb_-r2icil%2Fmgs={%22heartbeat%22:1607008145749}")
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("http get error", err)
		return
	}
	//函数结束后关闭相关链接
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusOK {
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}
		bodyString := string(bodyBytes)
		fmt.Println(bodyString)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	doc.Find("ol.grid_view li").Find(".hd").Each(func(index int, ele *goquery.Selection) {
		movieUrl, _ := ele.Find("a").Attr("href")
		fmt.Println(strings.Split(movieUrl, "/")[4], ele.Find(".title").Eq(0).Text())
	})



	//goqueryRR


	//find second php video links



	//links := collectlinks.All(resp.Body)
	//for _, link := range links {
	//	absolute := urlJoin(link, url)
	//	if url != " " {
	//		if !visited[absolute] {
	//			fmt.Println("parse url", absolute)
	//			go func() {
	//				queue <- absolute
	//			}()
	//		}
	//	}
	//}
}

func urlJoin(href, base string) string {
	uri, err := url.Parse(href)
	if err != nil {
		return " "
	}
	baseUrl, err := url.Parse(base)
	if err != nil {
		return " "
	}
	return baseUrl.ResolveReference(uri).String()
}
package main

import (
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io"
	"io/ioutil"
	"log"
	"myCrawler/mgstage"
	"os"
	"strings"
	"sync"

	//"github.com/jackdanger/collectlinks"
	"net/http"
	"net/url"
	"time"
)

func check(err error){
	if err != nil {
		fmt.Println("出错了 请查看",err)
		panic(err)
	}
}

func getMgstageTodayURl() string{
	yyyymmdd := getTimeYYYYMMDD()
	url := "https://www.mgstage.com/search/cSearch.php?search_word=&sale_start_range=%s-%s&sort=new&list_cnt=120&type=top"
	return fmt.Sprintf(url,yyyymmdd,yyyymmdd)
}
func getTimeYYYYMMDD() string{
	t := time.Now()
	return fmt.Sprintf("%d.%02d.%02d",
		t.Year(), t.Month(), t.Day())
}
func getTimeYYYYMMDDHHMMSS() string{
	t := time.Now()
	return fmt.Sprintf("%d-%02d-%02dT %02d:%02d:%02d",
		t.Year(), t.Month(), t.Day(),
		t.Hour(), t.Minute(), t.Second())
}



func getSavePath(url string,id string) string{
	urlSplits :=  strings.Split(url,"/")

	fileDir := fmt.Sprintf("assets/%s",id)
	os.Mkdir(fileDir,os.ModePerm)

	savePath := fmt.Sprintf("assets/%s/%s",id,  urlSplits[len(urlSplits)-1]  )
	return savePath
}
func downlaodThanSave(url string,path string){
	if(url==""){
		return
	}
	response, e := http.Get(url)
	if e != nil {
		log.Fatal(e)
	}
	defer response.Body.Close()

	file, err := os.Create(path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	_, err = io.Copy(file, response.Body)
	if err != nil {
		log.Fatal(err)
	}
	//fmt.Println("Success!")
}

func downloadImages(v *mgstage.Video){
	for _,imageUrl := range v.Images{
		downlaodThanSave(string(imageUrl),getSavePath(string(imageUrl),v.ID))
	}
}

var queue 		= make(chan *mgstage.Video, 3)
var queueVideos = make(chan *mgstage.Video, 3)
var queueImages = make(chan *mgstage.Video, 10)
var wg sync.WaitGroup

func main() {

	go func() {
		getList()
	}()


	go func() {
		DetailPageToImagesVideo()
	}()



	go func() {
		for video := range queueVideos {
			downloadVideo(video)
		}
	}()
	go func() {
		for video := range queueImages {
			downloadImages(video)
		}
	}()

	wg.Add(2)
	wg.Wait()

}

var videoMap 	= map[string] *mgstage.Video{}
var videoArray	= make([]*mgstage.Video,0)

func requestGetString(url string) string {
	timeout := time.Duration(100000 * time.Second)
	client := &http.Client{ Timeout:timeout,}
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("Cookie", "coc=1; PHPSESSID=q62v530bnrtfsqcfmnj2beuhp0; uuid=33e76d09fc86e676de2e0e8b887c8ee5; adc=1; __ulfpc=202012061038023125; _ga=GA1.2.98736236.1607222284; _gid=GA1.2.293888026.1607222284; bWdzdGFnZS5jb20%3D-_lr_uf_-r2icil=962de34a-89e7-413c-81cb-4d47111a2ac7; _gat_UA-58252858-1=1; _gat_UA-158726521-1=1; bWdzdGFnZS5jb20%3D-_lr_hb_-r2icil%2Fmgs={%22heartbeat%22:1607224305081}; bWdzdGFnZS5jb20%3D-_lr_tabs_-r2icil%2Fmgs={%22sessionID%22:0%2C%22recordingID%22:%224-e33ab609-ebc7-4322-a136-500d9641fe5f%22%2C%22lastActivity%22:1607224310786}//Cache-Control: no-cache; PHPSESSID=97n84d8jcorg6t45epdvbjvl02; uuid=33e76d09fc86e676de2e0e8b887c8ee5")
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("http get error", err)
		return ""
	}
	//函数结束后关闭相关链接
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}
		bodyString := string(bodyBytes)
		//fmt.Println(bodyString)
		return bodyString
	} else {
		fmt.Println("错误号码:")
		fmt.Println(resp.StatusCode)

		return ""
	}
}

func DetailPageToImagesVideo(){
	for v := range queue {
		//go func(v *mgstage.Video) {
			fmt.Println("=============")

			fmt.Println("queue url  full", v.FullUrl)
			detailDoc ,err 	:= requestGetDocument(v.FullUrl)
			if err != nil {
				log.Fatal(err)
			}
			getVideoModel(v,detailDoc)

			if err != nil {
				log.Fatal(err)
			}
			GetVideoMP4(v)
			fmt.Println("queue url video",v.VideoUrl)
			//为什么这里不需要使用go func呢{} 因为buffered channel吗
			queueVideos <- v
			queueImages <- v
		//}(Video)

	}
}
func getList() {

	mgstageTodayURl := getMgstageTodayURl()
	fmt.Println("start mgstage", mgstageTodayURl)
	//先不下载 搞定所有的原始数据
	bodyString := requestGetString(mgstageTodayURl)
	if(bodyString!=""){
		doc, err := goquery.NewDocumentFromReader(strings.NewReader(bodyString))
		if err != nil {
			log.Fatal(err)
		}
		doc.Find("div.search_list > div > ul > li").Each(func(index int, ele *goquery.Selection) {
			v := mgstage.Video{}
			v.Url, _ = getVideoUrl(ele)
			v.FullUrl = getFullUrl(v.Url)
			v.ID 	= getVideoId(v.Url)
			v.Price = getPrice(ele)
			videoMap[v.ID] = &v
			videoArray = append(videoArray, &v)
			//fmt.Println(v.FullUrl)
			queue <- &v
		})
	}
}
func downloadVideo(v *mgstage.Video){
	if(v.VideoUrl!=""){
		fmt.Printf("视频文件 %s",v.ID,v.VideoUrl)
		downlaodThanSave(string(v.VideoUrl),getSavePath(string(v.VideoUrl),v.ID))
	}
	//if(v.VideoUrl==""){
	//	fmt.Printf("没有找到视频文件 %s %s",v.ID,v.FullUrl)
	//	return
	//}
}
func GetVideoMP4(v *mgstage.Video) {
	toGetVideoUrl := fmt.Sprintf("https://www.mgstage.com/sampleplayer/sampleRespons.php?pid=%s", v.Pid)
	jsonString := requestGetString(toGetVideoUrl)
	var mgstageVideo JSONMgstageVideo
	err := json.Unmarshal([]byte(jsonString), &mgstageVideo)
	if err == nil {
		//fmt.Printf("%#v \n %#v \n", err, mgstageVideo.Url)
	}
	v.VideoUrl = getMp4UrlFromISMURL(mgstageVideo.Url)
}

type JSONMgstageVideo struct {
	Url string `json:"url"`
}
func getLinkByItemDoc(selectorTd string ,doc *goquery.Document) mgstage.Link {
	tdEle := doc.Find(selectorTd)
	linkElem := tdEle.Find("a")
	href,_  := linkElem.Attr("href")
	text := tdEle.Text()

	text = cleanText(text)
	href = cleanText(href)
	if(href!="") {
		href = getFullUrl(href)
	}
	return mgstage.Link{ text , href }
}
func getLinkByItemSelection(selectorTd string ,doc *goquery.Selection) mgstage.Link {
	linkElem := doc.Find(selectorTd)
	href,_  := linkElem.Attr("href")
	text := linkElem.Text()
	text = cleanText(text)
	href = cleanText(href)
	if(href!="") {
		href = getFullUrl(href)
	}
	return mgstage.Link{ text , href }
}
func cleanText(str string) string{
	s := strings.TrimLeft(str," ")
	s = strings.Trim(str,"\n")
	s = strings.Trim(str,"\r")
	s = strings.TrimSpace(str)
	return s
}

func getMp4UrlFromISMURL(url string)string{
	//url:="https://sample.mgstage.com/sample/nanpatv/200gana/2396/200gana-2396_20201202T125901.ism/request?uid=10000000-0000-0000-0000-00000000000a&amp;pid=16b6ae62-e6d6-412c-afda-c8b4709c86eb"
	end := strings.Index(url,"/request?")
	lastStr  := url[:end]
	lastStr = strings.Replace(lastStr,".ism",".mp4",1)
	//fmt.Println(lastStr)
	return lastStr
}

func getDocText (selector string ,doc *goquery.Document) string {
	t := doc.Find(selector).Text()
	return cleanText(t)
}
func getVideoModel(v *mgstage.Video,detailDoc *goquery.Document) (*mgstage.Video,error) {

	v.Title = getDocText("#center_column > div.common_detail_cover > h1",detailDoc)

	v.CountFavorite = getDocText("#playing > dl.detail_fav_cnt",detailDoc)
	v.CountPlay = detailDoc.Find("#playing > dl.playing").Text()


	v.Actor = getLinkByItemDoc("#center_column > div.common_detail_cover > div.detail_left > div > table:nth-child(3) > tbody > tr:nth-child(1) > td",detailDoc)
	v.Manufacturer =  getLinkByItemDoc("#center_column > div.common_detail_cover > div.detail_left > div > table:nth-child(3) > tbody > tr:nth-child(2) > td",detailDoc)

	v.TimeLong =  getDocText("#center_column > div.common_detail_cover > div.detail_left > div > table:nth-child(3) > tbody > tr:nth-child(3) > td",detailDoc)
	v.StartDate =  getDocText("#center_column > div.common_detail_cover > div.detail_left > div > table:nth-child(3) > tbody > tr:nth-child(5) > td",detailDoc)
	v.SaleDate =  getDocText("#center_column > div.common_detail_cover > div.detail_left > div > table:nth-child(3) > tbody > tr:nth-child(6) > td",detailDoc)

	v.Series = getLinkByItemDoc("#center_column > div.common_detail_cover > div.detail_left > div > table:nth-child(3) > tbody > tr:nth-child(7) > td",detailDoc)
	v.Company = getLinkByItemDoc("#center_column > div.common_detail_cover > div.detail_left > div > table:nth-child(3) > tbody > tr:nth-child(7) > td",detailDoc)

	detailDoc.Find("#center_column > div.common_detail_cover > div.detail_left > div > table:nth-child(3) > tbody > tr:nth-child(9) > td").Each(func(index int, ele *goquery.Selection) {
		link := getLinkByItemSelection("a",ele)
		v.Tags = append( v.Tags, link)
	})

	v.Rate =  getDocText("#center_column > div.common_detail_cover > div.detail_left > div > table:nth-child(3) > tbody > tr:nth-child(11) > td > span",detailDoc)

	v.ImageHead,_ = detailDoc.Find("#center_column > div.common_detail_cover > div.detail_left > div > div > h2 > img").Attr("src")
	v.Description =  getDocText("#introduction > dd > p.txt.introduction",detailDoc)

	//photos
	detailDoc.Find("#sample-photo > dd > ul > li > a").Each(func(index int, ele *goquery.Selection) {
		link, _ := ele.Attr("href")
		v.Images  = append(v.Images,link)
	})

	pidUrl, _ := detailDoc.Find(".button_sample").Attr("href")
	pidUrlSplits := strings.Split(pidUrl,"/")
	v.Pid = pidUrlSplits[len(pidUrlSplits)-1]
 	return v, nil
}

func requestGetDocument( url string ) (*goquery.Document, error)  {
	timeout := time.Duration(5 * time.Second)
	client := &http.Client{
		Timeout:timeout,
	}
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("Cookie", "coc=1; PHPSESSID=q62v530bnrtfsqcfmnj2beuhp0; uuid=33e76d09fc86e676de2e0e8b887c8ee5; adc=1; __ulfpc=202012061038023125; _ga=GA1.2.98736236.1607222284; _gid=GA1.2.293888026.1607222284; bWdzdGFnZS5jb20%3D-_lr_uf_-r2icil=962de34a-89e7-413c-81cb-4d47111a2ac7; _gat_UA-58252858-1=1; _gat_UA-158726521-1=1; bWdzdGFnZS5jb20%3D-_lr_hb_-r2icil%2Fmgs={%22heartbeat%22:1607224305081}; bWdzdGFnZS5jb20%3D-_lr_tabs_-r2icil%2Fmgs={%22sessionID%22:0%2C%22recordingID%22:%224-e33ab609-ebc7-4322-a136-500d9641fe5f%22%2C%22lastActivity%22:1607224310786}//Cache-Control: no-cache; PHPSESSID=97n84d8jcorg6t45epdvbjvl02; uuid=33e76d09fc86e676de2e0e8b887c8ee5")
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("error", err)
		return nil,err
	}
	//函数结束后关闭相关链接
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusOK {
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}
		bodyString := string(bodyBytes)
		//doc, err := goquery.NewDocumentFromReader(resp.Body)
		doc, err := goquery.NewDocumentFromReader(strings.NewReader(bodyString))
		if err != nil {
			fmt.Println("error", err)
			log.Fatal(err)
			return nil,err
		}
		return doc,nil
	} else {
		fmt.Println("错误号码:")
		fmt.Println(resp.StatusCode)
		panic(fmt.Sprintf("status code is %2d", resp.StatusCode))
		return nil,err
	}
}
func getVideoUrl(ele *goquery.Selection) (string, bool) {
	return ele.Find("a").Eq(0).Attr("href")
}


func getPrice(ele *goquery.Selection) string {
	return ele.Find(".price").Text()
}

func getVideoId(avUrl string) string {
	return strings.Split(avUrl, "/")[3]
}
func getFullUrl(avUrl string) string {
	return "https://www.mgstage.com" + avUrl
}
func getFullUrlJoin(href, base string) string {
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





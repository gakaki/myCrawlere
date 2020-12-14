package main

import (
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io/ioutil"
	"myCrawler/scanlibs"
	"myCrawler/utils"
	"net/http"
	"os"
	"reflect"
	"strconv"
	"strings"
	"sync"
)
func shareHeader(req *http.Request ) {
	//req.Header.Add("authority", "leetcode-cn.com")
}

func CheckError(err error){
	if err != nil {
		fmt.Printf("%d", err)
	}
}


var scanlibItems =  make([]*scanlibs.ScanLibItem,0)
//var queueFethItems = make(chan *scanlibs.ScanLibItem, 3)
var queueFethItems = make(chan int, 3)
var wg sync.WaitGroup

func getItemsModel(doc *goquery.Document)  {
	doc.Find("article").Each(func(index int, ele *goquery.Selection) {
		var item = scanlibs.ScanLibItem{}
		item.ID, _  = ele.Attr("id")
		item.URL, _ = ele.Find("a").Attr("href")
		item.CreatedAt, _ = ele.Find("time").Attr("datetime")
		item.Thumbnil, _ = ele.Find("img.aligncenter").Attr("src")
		item.Title = ele.Find("b").Text()
		item.SubTitle = ele.Find("p").Text()
		item.IsVideo = strings.Index(item.URL,"video") >= 0
		scanlibItems = append(scanlibItems, &item)
	})
}

func getLastPageIndex(doc *goquery.Document) int  {
	t := doc.Find(".page-numbers:nth-last-child(2)").Text()
	t = strings.Replace(t,",","",2)
	d,_ :=  strconv.Atoi(t)
	return d
}
func getOtherPage(){
	wg.Add(1)
	index := <- queueFethItems

	url := fmt.Sprintf("https://scanlibs.com/page/%d/", index )
	fmt.Println("request page ",index)
	//need retry
	doc, err := utils.RequestGetDocument(url)
	if(err!=nil){
		CheckError(err)
		return
	}
	getItemsModel(doc)

	data , err := json.Marshal(&scanlibItems)
	WriteToJSON([]byte(data),data)
	wg.Done()
}
func getFirstPage() int {
	url := "https://scanlibs.com/"
	doc, err := utils.RequestGetDocument(url)
	if(err!=nil){
		CheckError(err)
		panic(err)
	}
	getItemsModel(doc)
	return getLastPageIndex(doc)
}


func main() {
	lastPageIndex := getFirstPage()
	fmt.Println(lastPageIndex)
	for i := 2; i <= lastPageIndex; i++ {

		queueFethItems <- i



		//go func() {

			getOtherPage()

		//}()


	}
	//wg.Wait()
	//fmt.Println(
	//	len(questionJSON.StatStatusPairs),
	//	len(questionJSONTranslation.Data.Translations),
	//	len(favoritesJSON),
	//	len(tagsJSON.Topics),
	//	len(companyJSON.Data.InterviewHotCards),
	//	len(companyQuestionsJSON.Data.CompanyTag.Questions),
	//)

}
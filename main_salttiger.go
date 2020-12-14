package main

import (
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	Satiger "myCrawler/satiger"
	"myCrawler/utils"
	"regexp"
	"strings"
)

var salttigerItems = make([]*Satiger.SatigerItem, 0)

//var queueFethItems = make(chan *scanlibs.ScanLibItem, 3)
//var queueFethItems = make(chan int, 3)
//var wg sync.WaitGroup

func regexFind(regexStr string,str string)  string {
	reg := regexp.MustCompile(regexStr)
	res := reg.FindString(str)
	//fmt.Println(res)
	return res
}
func getDetailPage(item *Satiger.SatigerItem) {
	//wg.Add(1)
	//index := <- queueFethItems

	fmt.Println("request page ", item.URL)
	//need retry
	doc, err := utils.RequestGetDocument(item.URL)
	if err != nil {
		panic(err)
		return
	}
	doc.Find("article").Each(func(index int, ele *goquery.Selection) {

		item.ID, _ = ele.Attr("id")
		item.Thumbnil, _ = ele.Find("div > p:nth-child(1) > img").Attr("src")

		tmpText := ele.Find("strong").Text()
		item.PubDate = regexFind(`\d{4}.\d{1,2}`,tmpText) 	// 出版时间：2020.12
		//item.Yearmonth = regexFind(`\d{4}年\w{1,4}月`,item.Yearmonth) 	// 出版时间：2020.12

		officalA := ele.Find("div > p:nth-child(1) > strong > a:nth-child(2)")
		item.OfficalPress = officalA.Text()
		item.OfficalURL, _ = officalA.Attr("href")

		ele.Find("article strong > a[href*=ed2k]").Each(func(index int, it *goquery.Selection) {
			link,_ := it.Attr("href")
			item.OtherLinks = append(item.OtherLinks,link )
		})


		officalBaidu := ele.Find("article strong > a[href*=baidu]")
		item.BaiduURL = officalBaidu.Text()
		item.BaiduCode = regexFind(`提取码    ：\w{1,4}`,tmpText)
		item.BaiduCode = strings.Replace(item.BaiduCode,"提取码    ：","",1)

		item.Description,_ = ele.Find("div.entry-content").Html()
		item.Description = regexFind(`<p>内容简介([\s\S]*)`,item.Description)
		item.Description = strings.Replace(item.Description,"<p>内容简介：</p>","",1)
		item.CreatedAt,_ = ele.Find("footer > a:nth-child(1) > time").Attr("datetime")

		ele.Find("footer > a[rel*=tag]").Each(func(index int, e *goquery.Selection) {
			tag := Satiger.Tag{}
			tag.URL,_ 	= e.Attr("href")
			tag.Name 	= e.Text()
			item.Tags = append(item.Tags,tag)
		})

		jsonStr, _ := json.Marshal(salttigerItems)
		body := []byte(jsonStr)
		utils.WriteToJSONByFileName(body, "saltigeryitian")


	})
	//wg.Done()
}
func getArchives() {
	url := "https://salttiger.com/archives/"
	doc, err := utils.RequestGetDocument(url)
	if err != nil {
		panic(err)
	}
	doc.Find("ul.car-list li").Each(func(index int, ele *goquery.Selection) {

		createdAt := ele.Find("span.car-yearmonth").Text()
		ele.Find("ul.car-monthlisting li").Each(func(index int, ele *goquery.Selection) {
			var item = Satiger.SatigerItem{}

			item.Yearmonth = createdAt
			itemA := ele.Find("a")
			item.Title = itemA.Text()
			item.URL, _ = itemA.Attr("href")
			salttigerItems = append(salttigerItems, &item)
		})

	})
}

func main() {
	getArchives()
	fmt.Println(&salttigerItems)
	for i := 0; i <= len(salttigerItems); i++ {
		//go func() {
			getDetailPage(salttigerItems[i])
		//}()
	}
	//wg.Wait()
	fmt.Println(
		len(salttigerItems),
	)


}

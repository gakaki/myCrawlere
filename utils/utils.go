package utils

import (
	fmt "fmt"
	"github.com/PuerkitoBio/goquery"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"reflect"
	"strings"
	"time"
)

func WriteToFile(body []byte, name string){
	fileName := fmt.Sprintf("json/%s",name)
	_ = os.Mkdir("json", os.ModePerm)
	_ = ioutil.WriteFile(fileName, body, os.ModePerm)
}
func WriteToJSONByFileName(body []byte, name string){
	fileName := fmt.Sprintf("json/%s.json",name)
	_ = os.Mkdir("json", os.ModePerm)
	_ = ioutil.WriteFile(fileName, body, os.ModePerm)
}
func WriteToJSON(body []byte,json interface{}){
	typeOfA := reflect.TypeOf(json)
	//fmt.Println(typeOfA.Name(), typeOfA.Kind())
	fileName := fmt.Sprintf("json/%s.json",typeOfA.Name())
	_ = os.Mkdir("json", os.ModePerm)

	//jsonByte, _ := json.([]byte)
	_ = ioutil.WriteFile(fileName, body, os.ModePerm)
}

func RequestGetDocument( url string ) (*goquery.Document, error)  {
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

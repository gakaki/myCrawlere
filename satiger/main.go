package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"reflect"
	"strings"
)
func shareHeader(req *http.Request ) {
	req.Header.Add("authority", "leetcode-cn.com")
	req.Header.Add("pragma", "no-cache")
	req.Header.Add("cache-control", "no-cache")
	req.Header.Add("accept", "*/*")
	req.Header.Add("origin", "https://leetcode-cn.com")
	req.Header.Add("x-csrftoken", "kxbbHPzz4OE5jNNZ2hbTIB5i1sjsltwDgjWESMR5tiKCkiJzW8UazVqvPv1qfIRh")
	req.Header.Add("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/81.0.4033.2 Safari/537.36")
	req.Header.Add("content-type", "application/json")
	req.Header.Add("sec-fetch-site", "same-origin")
	req.Header.Add("sec-fetch-mode", "cors")
	req.Header.Add("sec-fetch-dest", "empty")
	req.Header.Add("referer", "https://leetcode-cn.com/problemset/all/")
	req.Header.Add("accept-language", "zh-CN,zh;q=0.9")
	req.Header.Add("cookie", "__auc=43e9012a175ed9b43d6d34e3041; gr_user_id=0a2d3940-43ae-42bc-9af9-880c481ecc57; _ga=GA1.2.1218529367.1606009179; grwng_uid=80cc6ea3-cfbb-41ce-8830-56144952e39a; a2873925c34ecbd2_gr_last_sent_cs1=si-quan-jia; _gid=GA1.2.337157349.1607263523; Hm_lvt_fa218a3ff7179639febdb15e372f411c=1606480937,1606480940,1606978596,1607268523; _gat_gtag_UA_131851415_1=1; a2873925c34ecbd2_gr_session_id=7e648219-90ff-4733-8cc4-3e36e37d7fe9; a2873925c34ecbd2_gr_last_sent_sid_with_cs1=7e648219-90ff-4733-8cc4-3e36e37d7fe9; a2873925c34ecbd2_gr_session_id_7e648219-90ff-4733-8cc4-3e36e37d7fe9=true; __asc=23756b4d17638cbf8371f644e23; Hm_lpvt_fa218a3ff7179639febdb15e372f411c=1607270669; a2873925c34ecbd2_gr_cs1=si-quan-jia; csrftoken=XBNJgVxePgLqRMItVofQWfAgTqyQcrqCycVfDwXOlMFUJjiMWcriQj2eUhV7oILH; LEETCODE_SESSION=eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJfYXV0aF91c2VyX2lkIjoiMjA5MTU0OSIsIl9hdXRoX3VzZXJfYmFja2VuZCI6ImF1dGhlbnRpY2F0aW9uLmF1dGhfYmFja2VuZHMuUGhvbmVBdXRoZW50aWNhdGlvbkJhY2tlbmQiLCJfYXV0aF91c2VyX2hhc2giOiIzZWMxMWNkZjM5NjM5NzYwNDM5MzZlNDg2YzdjYTdiOGZhNjMyNjA2MzhmMzJkYzgxOGY4MjhkOTU4MTYyM2VmIiwiaWQiOjIwOTE1NDksImVtYWlsIjoiamNnbHFtb3l4QDEyNi5jb20iLCJ1c2VybmFtZSI6InNpLXF1YW4tamlhIiwidXNlcl9zbHVnIjoic2ktcXVhbi1qaWEiLCJhdmF0YXIiOiJodHRwczovL2Fzc2V0cy5sZWV0Y29kZS1jbi5jb20vYWxpeXVuLWxjLXVwbG9hZC91c2Vycy9zaS1xdWFuLWppYS9hdmF0YXJfMTYwNTU4NDk5MC5wbmciLCJwaG9uZV92ZXJpZmllZCI6dHJ1ZSwiX3RpbWVzdGFtcCI6MTYwNzI3MDY3OC4yNzAyNTcyfQ.x_nRmwZl-ETgjZWKTKOoc8jJT6H5FQOB64sZfVYvXAA")
	req.Header.Add("user-agent", " Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.4280.67 Safari/537.36")
	req.Header.Add("x-csrftoken", " 3extWqvqKpBWX0EWmqMTkUeDhxp2RxtwKkSzTqygBbFy2Uqp84sKuQRu0OHvyhoB")
	req.Header.Add("x-definition-name", " interviewCard,interviewCompanyOptions,companyTag,jobsCompany")
	req.Header.Add("x-operation-name", " companyTag")
	req.Header.Add("x-timezone", " Asia/Shanghai")


}
func request(url string ,query string) ([]byte,error){
	method := "GET"
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		fmt.Println(err)
		return nil,err
	}
	shareHeader(req)
	return baseRequest(req)
}
func baseRequest( req *http.Request )([]byte,error){
	shareHeader(req)

	client := &http.Client {}
	res, err := client.Do(req)

	if err != nil {
		fmt.Println(err)
		return nil,err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return nil,err
	}
	//fmt.Println(string(body))
	return body,nil
}
func graphql(postQL string) ([]byte,error){
	url := "https://leetcode-cn.com/graphql"
	method := "POST"
	payload := strings.NewReader(postQL)
	req, err := http.NewRequest(method, url, payload)
	if err != nil {
		fmt.Println(err)
		return nil,err
	}
	return baseRequest(req)
}

func checkError(err error){
	if err != nil {
		fmt.Printf("%d", err)
	}
}

func writeToJSON(body []byte,json interface{}){
	typeOfA := reflect.TypeOf(json)
	//fmt.Println(typeOfA.Name(), typeOfA.Kind())
	fileName := fmt.Sprintf("json/%s.json",typeOfA.Name())
	_ = os.Mkdir("json", os.ModePerm)

	//jsonByte, _ := json.([]byte)
	_ = ioutil.WriteFile(fileName, body, os.ModePerm)
}

func getAllQuestionsTranslation() ( QuestionCHNJSON,  error) {
	body, _ := graphql(`{"operationName":"getQuestionTranslation","variables":{},"query":"query getQuestionTranslation($lang: String) {\n  translations: allAppliedQuestionTranslations(lang: $lang) {\n    title\n    questionId\n    __typename\n  }\n}\n"}`)
	var questionCHNJSON QuestionCHNJSON
	err := json.Unmarshal(body, &questionCHNJSON)
	checkError(err)
	writeToJSON(body,questionCHNJSON)
	return questionCHNJSON,err
}

func getAllQuestions() (QuestionJSON, error) {
	body, _ := request("https://leetcode-cn.com/api/problems/all/", "")
	var questionJSON QuestionJSON
	err := json.Unmarshal(body, &questionJSON)
	checkError(err)
	writeToJSON(body,questionJSON)
	return questionJSON, err
}
func getFavorites() (FavoriteJSON, error) {
	body, _ := request("https://leetcode-cn.com/problems/api/favorites/", "")
	var favoriteJSON FavoriteJSON
	err := json.Unmarshal(body, &favoriteJSON)
	checkError(err)
	writeToJSON(body,favoriteJSON)
	return favoriteJSON, err
}
func getTags() (TagJSON, error) {
	body, _ := request("https://leetcode-cn.com/problems/api/tags/", "")
	var tagJSON TagJSON
	err := json.Unmarshal(body, &tagJSON)
	checkError(err)
	writeToJSON(body,tagJSON)
	return tagJSON, err
}


func getCompanys() (CompanyJSON, error) {
	body, _ := graphql(`{"operationName":"getHotCompanies","variables":{"onlyCompanyCards":true},"query":"query getHotCompanies($onlyCompanyCards: Boolean!) {\n  interviewHotCards(onlyCompanyCards: $onlyCompanyCards) {\n    id\n    numQuestions\n    company {\n      name\n      slug\n      imgUrl\n      __typename\n    }\n    __typename\n  }\n}\n"}`)
	var companyJSON CompanyJSON
	err := json.Unmarshal(body, &companyJSON)
	checkError(err)
	writeToJSON(body,companyJSON)
	return companyJSON,err
}

func getCompnayQuestions() (CompanyQuestionJSON, error) {
	body, _ := graphql(`{"operationName":"companyTag","variables":{"slug":"bytedance"},"query":"query companyTag($slug: String!) {\n  interviewCard(companySlug: $slug) {\n    id\n    isFavorite\n    isPremiumOnly\n    privilegeExpiresAt\n    jobsCompany {\n      name\n      jobPostingNum\n      isVerified\n      description\n      logo\n      logoPath\n      postingTypeCounts {\n        count\n        postingType\n        __typename\n      }\n      industryDisplay\n      scaleDisplay\n      financingStageDisplay\n      website\n      legalName\n      __typename\n    }\n    __typename\n  }\n  interviewCompanyOptions(query: $slug) {\n    id\n    __typename\n  }\n  companyTag(slug: $slug) {\n    name\n    id\n    imgUrl\n    translatedName\n    frequencies\n    questions {\n      title\n      translatedTitle\n      titleSlug\n      questionId\n      stats\n      status\n      questionFrontendId\n      difficulty\n      frequencyTimePeriod\n      topicTags {\n        id\n        name\n        slug\n        translatedName\n        __typename\n      }\n      __typename\n    }\n    __typename\n  }\n  jobsCompany(companySlug: $slug) {\n    name\n    legalName\n    logo\n    description\n    website\n    industryDisplay\n    scaleDisplay\n    financingStageDisplay\n    isVerified\n    __typename\n  }\n}\n"}`)
	var companyQuestionJSON CompanyQuestionJSON
	//fmt.Println(string(body))
	err := json.Unmarshal(body, &companyQuestionJSON)
	checkError(err)
	writeToJSON(body,companyQuestionJSON)
	return companyQuestionJSON,err
}



func main() {

	//获得所有的题目
	questionJSON, _ := getAllQuestions()
	//获得所有题目的中文名字
	questionJSONTranslation , _ := getAllQuestionsTranslation()
	//所有比赛题目合集
	favoritesJSON,_ := getFavorites()
	//所有分类合集
	tagsJSON,_ := getTags()

	companyJSON,_ := getCompanys()
	//model数据重新组织
	//写入elasticsearch 和 dgraphql数据库

	companyQuestionsJSON,_ := getCompnayQuestions()

	fmt.Println(
		len(questionJSON.StatStatusPairs),
		len(questionJSONTranslation.Data.Translations),
		len(favoritesJSON),
		len(tagsJSON.Topics),
		len(companyJSON.Data.InterviewHotCards),
		len(companyQuestionsJSON.Data.CompanyTag.Questions),
	)

}
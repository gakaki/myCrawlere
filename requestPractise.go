package main

import (
	"fmt"
	"sync"
	"time"
)



var channelUrls chan string = make(chan string,1)
var wgc sync.WaitGroup
func SendUrlToChannel(){
	for i:=0 ; i< 4; i++{
		fmt.Println("入 ",i)
		channelUrls <- fmt.Sprintf("%d",i)
	}
}
func UrlToRequest(){
	for x := range channelUrls{
		go func(s string) {
			time.Sleep(time.Second)
			fmt.Println("出 ",s)
		}(x)
	}
}
func main(){
	//go SendUrlToChannel()
	//go UrlToRequest()
	go func() {
		defer wgc.Done()
		SendUrlToChannel()
	}()
	go func() {
		defer wgc.Done()
		UrlToRequest()
	}()
	wgc.Add(2)
	wgc.Wait()
}

package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)
var num = 0

func curl(){
	client := &http.Client{}
	req, err := http.NewRequest("GET", "http://127.0.0.1:8082/product/order?productID=1&userId=1", nil)
	if err != nil {
		fmt.Println(err.Error())
	}
	req.AddCookie(&http.Cookie{
		Name:       "userId",
		Value:      "iWl3MIbY-I-6VlCidbb6hw",
		Path:       "/",
	})
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err.Error())
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err.Error())
	}
	if string(body) == "false" {
		num++
	}
}

func main() {
	ch := make(chan int, 10000)
	for i:= 0;i<5000; i++ {
		go curl()
		ch <- 1
	}
	for i:= 0;i<5000; i++ {
		num++
		<- ch
	}
	time.Sleep(100*time.Second)
	fmt.Println(num)
}

package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

var baseURL string = "https://kr.indeed.com/%EC%B7%A8%EC%97%85?q=python&l="

func main() {
	getPages()
}

func getPages() int {
	res, err := http.Get(baseURL)
	checkErr(err)
	checkCode(res)

	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	checkCode(res)

	fmt.Println(doc)

	return 0
}

func checkErr(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

func checkCode(res *http.Response) {
	if res.StatusCode != 200 {
		log.Fatalln("Request Failed with Status: ", res.StatusCode)
	}
}

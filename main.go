package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/PuerkitoBio/goquery"
)

type extractedJob struct {
	id       string
	title    string
	location string
	salary   string
	summary  string
	company  string
}

var baseURL string = "https://kr.indeed.com/%EC%B7%A8%EC%97%85?q=python"

func main() {
	totalPages := getPages()
	for i := 0; i < totalPages; i++ {
		getPage(i)
	}
}

func getPage(page int) {
	pageURL := baseURL + "&start=" + strconv.Itoa(page*10)
	fmt.Println("Requesting", pageURL)
	res, err := http.Get(pageURL)
	checkErr(err)
	checkCode(res)

	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)

	searchCards := doc.Find(".tapItem")
	searchCards.Each(func(i int, card *goquery.Selection) {
		id, _ := card.Attr("data-jk")

		title := card.Find(".jobTitle>span").Text()
		company := card.Find(".companyName").Text()
		location := card.Find(".companyLocation").Text()
		summary := card.Find(".job-snippet").Text()
		fmt.Println("id:", id)
		fmt.Println("title:", title)
		fmt.Println("company:", company)
		fmt.Println("location:", location)
		fmt.Println("summary:", summary+"\n")
	})
}

func getPages() int {
	pages := 0
	res, err := http.Get(baseURL)
	checkErr(err)
	checkCode(res)

	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	checkCode(res)

	doc.Find(".pagination").Each(func(i int, s *goquery.Selection) {
		pages = s.Find("a").Length()
	})

	return pages
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

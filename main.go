package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

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
var viewURL string = "https://kr.indeed.com/%EC%B1%84%EC%9A%A9%EB%B3%B4%EA%B8%B0?jk="

func main() {
	var jobs []extractedJob
	c := make(chan []extractedJob)
	totalPages := getPages()

	for i := 0; i < totalPages; i++ {
		go getPage(i, c)
	}

	for i := 0; i < totalPages; i++ {
		extractedJobs := <-c
		jobs = append(jobs, extractedJobs...)
	}

	wirteJobs(jobs)
	logJobs(jobs)
	fmt.Println("Done, extracted", len(jobs))
}

func getPage(page int, mainC chan<- []extractedJob) {
	var jobs []extractedJob
	c := make(chan extractedJob)
	pageURL := baseURL + "&start=" + strconv.Itoa(page*10)
	fmt.Println("Requesting", pageURL)
	res, err := http.Get(pageURL)
	checkErr(err)
	checkCode(res)

	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)

	searchCards := doc.Find(".tapItem")
	searchCards.Each(func(i int, card *goquery.Selection) {
		go extractJob(card, c)
	})

	for i := 0; i < searchCards.Length(); i++ {
		job := <-c
		jobs = append(jobs, job)
	}
	mainC <- jobs
}

func extractJob(card *goquery.Selection, c chan<- extractedJob) {
	id, _ := card.Attr("data-jk")
	title := cleanString(card.Find(".jobTitle>span").Text())
	company := cleanString(card.Find(".companyName").Text())
	location := cleanString(card.Find(".companyLocation").Text())
	salary := cleanString(card.Find(".salary-snippet").Text())
	summary := cleanString(card.Find(".job-snippet").Text())
	c <- extractedJob{
		id:       id,
		title:    title,
		company:  company,
		location: location,
		salary:   salary,
		summary:  summary,
	}
}

func cleanString(str string) string {
	return strings.Join(strings.Fields(strings.TrimSpace(str)), " ")
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

func wirteJobs(jobs []extractedJob) {
	file, err := os.Create("jobs.csv")
	checkErr(err)

	w := csv.NewWriter(file)
	defer w.Flush()

	headers := []string{"ID", "TITLE", "COMPANY", "LOCATION", "SALARY", "SUMMARY"}
	wErr := w.Write(headers)
	checkErr(wErr)

	for _, job := range jobs {
		jobLink := "<a href='" + viewURL + job.id + "'>" + "</a>"
		jobSlice := []string{jobLink, job.title, job.company, job.location, job.salary, job.summary}
		jwErr := w.Write(jobSlice)
		checkErr(jwErr)
	}
}

func logJobs(jobs []extractedJob) {
	for idx, job := range jobs {
		fmt.Println("[", idx, "]", "id:", job.id)
		fmt.Println("title:", job.title)
		fmt.Println("company:", job.company)
		fmt.Println("location:", job.location)
		fmt.Println("salary:", job.salary)
		fmt.Println("summary:", job.summary+"\n")
	}
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

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
}

var baseURL string = "https://kr.indeed.com/jobs?q=unreal"

func main() {
	var jobs []extractedJob
	totalPages := getPages()

	for i := 0; i < totalPages; i++ {
		extractedJobs := getPage(i)
		jobs = append(jobs, extractedJobs...)
	}

	writeJobs(jobs)
	fmt.Println("extracted: ", len(jobs))
}

func checkError(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

func checkCode(res *http.Response) {
	if res.StatusCode != 200 {
		log.Fatalln("status code error: %d %s", res.StatusCode, res.Status)
	}
}

func cleanString(str string) string {
	return strings.Join(strings.Fields(strings.TrimSpace(str)), " ")
}

func getPages() int {
	pages := 0
	res, err := http.Get(baseURL)
	checkError(err)
	checkCode(res)

	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	checkError(err)
	doc.Find(".pagination").Each(func(i int, s *goquery.Selection) {
		pages = s.Find("a").Length()
	})
	fmt.Println("Total Page Number: ", pages)
	return pages
}

func extractJob(job *goquery.Selection) extractedJob {
	id, _ := job.Attr("data-jk")
	title := cleanString(job.Find(".title>a").Text())
	location := cleanString(job.Find(".sjcl").Text())
	salary := cleanString(job.Find(".salaryText").Text())
	summary := cleanString(job.Find(".summary").Text())
	return extractedJob{
		id:       id,
		title:    title,
		location: location,
		salary:   salary,
		summary:  summary,
	}
}

func getPage(page int) []extractedJob {
	var jobs []extractedJob
	pageUrl := baseURL + "&start=" + strconv.Itoa(page*10)
	res, err := http.Get(pageUrl)
	checkError(err)
	checkCode(res)
	fmt.Println(pageUrl)

	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	checkError(err)

	jobCards := doc.Find(".jobsearch-SerpJobCard")
	jobCards.Each(func(i int, s *goquery.Selection) {
		job := extractJob(s)
		jobs = append(jobs, job)
	})

	return jobs
}

func writeJobs(jobs []extractedJob) {
	file, err := os.Create("jobs.csv")
	checkError(err)

	w := csv.NewWriter(file)
	defer w.Flush()

	headers := []string{"ID", "Title", "Location", "Salary", "Summary"}
	wErr := w.Write(headers)
	checkError(wErr)

	for _, job := range jobs {
		jobSlice := []string{"https://kr.indeed.com/viewjob?jk=" + job.id, job.title, job.location, job.salary, job.summary}
		jwErr := w.Write(jobSlice)
		checkError(jwErr)
	}
}

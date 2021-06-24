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
}

var baseURL string = "https://kr.indeed.com/jobs?q=unreal"

func main() {
	totalPages := getPages()

	for i := 0; i < totalPages; i++ {
		getPage(i)
	}
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

func getPage(page int) {
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
		id, exists := s.Attr("data-jk")
		title := s.Find(".title>a").Text()
		location := s.Find(".sjcl").Text()
		fmt.Println(id, exists, title, location)
	})
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

func cleanString(str string) {

}

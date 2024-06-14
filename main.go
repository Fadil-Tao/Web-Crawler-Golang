package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"

	"github.com/PuerkitoBio/goquery"
)

var target string
var wg sync.WaitGroup

func main() {
	if len(os.Args) == 1 {
		fmt.Println("Plis input url")
		return
	}
	target = os.Args[1]

	queue := make(chan string)
	filteredqueue := make(chan string)
	wg.Add(1)
	go func() { queue <- target }()
	go func() {
		var urls = make(map[string]bool)
		for uri := range queue {
			if !urls[uri] {
				urls[uri] = true
				filteredqueue <- uri
			} else {
				wg.Done()
			}
		}
	}()
	for i := 0 ; i < 10 ; i++{
		go func(){
			for NextLink := range filteredqueue{
				Crawling(NextLink,queue)
				wg.Done()
			}
		}()
	} 
	wg.Wait()
}

func Crawling(uri string, queue chan string) {
	if uri == "" {
		return
	}

	fmt.Println("Fetching:.......", uri)
	response, err := http.Get(uri)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer response.Body.Close()

	doc, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		log.Fatal(err)
		return
	}

	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		val, exist := s.Attr("href")
		if exist {
			NextLink := TrimUrl(val)
			wg.Add(1)
			go func() {
				queue <- NextLink
			}()
		}
	})

}

func TrimUrl(uri string) string {
	uri = strings.TrimSuffix(uri, "/")
	validUrl, err := url.Parse(uri)
	if err != nil {
		log.Fatal(err)
		return ""
	}

	targetUrl, _ := url.Parse(target)
	if strings.Contains(validUrl.String(), targetUrl.Host) {
		return uri
	}
	return ""

}

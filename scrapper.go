package gowebgrep

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/gocolly/colly"
)

const (
	DefaultRPS = 1
)

type Scrapper struct {
	url              string
	Filter           *Filter
	RequestPerSecond *string         // WILL BE 1 PER DEFAULT NOT TO D DOS THE WEBSITE
	FoundPool        []string        // MAIN SLICE, ALL THE MATCHED STRINGS WILL BE KEEPED
	EndpointsChecked map[string]bool // TO PREVENT REPEATING ENDPOINTS WE ALREADY CHECKED
	File             *os.File
	collector        *colly.Collector
	timer            *Timer
}

func InitializeScrapper(filter *Filter, tm *Timer, rps *string) *Scrapper {
	urlEntered := os.Args[1]
	fmt.Println(urlEntered)
	fi, _ := os.Create("output.txt")

	return &Scrapper{
		url:              urlEntered,
		Filter:           filter,
		RequestPerSecond: rps,
		FoundPool:        make([]string, 10),
		EndpointsChecked: make(map[string]bool),
		File:             fi,
		collector:        colly.NewCollector(),
		timer:            tm,
	}
}

func (s *Scrapper) StartScrapping(endpoint string) {

	s.collector.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	s.collector.OnHTML("a[href]", func(e *colly.HTMLElement) {

		val, _ := e.DOM.Attr("href")
		doesStartWithSlash := s.checkPrefix(val, "/")

		if doesStartWithSlash {

			wasCheckedBefore := s.checkHasControlledBefore(val)

			if !wasCheckedBefore {
				s.timer.Mutex.Lock()
				s.EndpointsChecked[val] = true
				time.Sleep(time.Millisecond * time.Duration(s.timer.Milliseconds))
				s.timer.Mutex.Unlock()
				s.collector.Visit(s.url + val)
			}
		} else {
			doesStartWithHttp := s.checkPrefix(val, "http")
			if doesStartWithHttp {
				s.collector.Visit(val)
			}
		}

	})

	s.collector.OnResponse(func(r *colly.Response) {
		body := r.Body
		allFound := s.Filter.ApplyFilter(body)
		for i := 0; i < len(allFound); i++ {
			mailWithDownLine := append(allFound[i], 10)
			s.File.Write(mailWithDownLine)
			fmt.Println("We found an email ; ", string(allFound[i]))
		}
	})

	s.collector.Visit(s.url)
}

func (s *Scrapper) checkHasControlledBefore(endpoint string) bool {
	urlWithEndpoint := s.url + endpoint
	fmt.Println(urlWithEndpoint)
	_, ok := s.EndpointsChecked[urlWithEndpoint]
	return ok
}

func (s *Scrapper) checkPrefix(endpoint, prefix string) bool {
	return strings.HasPrefix(endpoint, prefix)
}

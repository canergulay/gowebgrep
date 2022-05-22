package gowebgrep

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/gocolly/colly"
)

const (
	DefaultRPS = 1
)

type Scrapper struct {
	url              string
	Filter           *Filter
	RequestPerSecond *string         // WILL BE 1 PER DEFAULT NOT TO D DOS THE WEBSITE
	FoundPool        map[string]bool // MAIN SLICE, ALL THE MATCHED STRINGS WILL BE KEEPED
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
		FoundPool:        make(map[string]bool),
		EndpointsChecked: make(map[string]bool),
		File:             fi,
		collector:        colly.NewCollector(colly.URLFilters(regexp.MustCompile(urlEntered))),
		timer:            tm,
	}
}

func (s *Scrapper) StartScrapping(endpoint string) {

	s.collector.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	// s.collector.OnHTML("link", func(e *colly.HTMLElement) {
	// 	val, _ := e.DOM.Attr("href")
	// 	s.checkAndLog(val)
	// })

	// s.collector.OnHTML("a[href]", func(e *colly.HTMLElement) {
	// 	val, _ := e.DOM.Attr("href")
	// 	s.checkAndLog(val)
	// })

	s.collector.OnResponse(func(r *colly.Response) {
		body := r.Body
		allFound := s.Filter.ApplyFilter(body, s.Filter.RegExp)

		for i := 0; i < len(allFound); i++ {
			s.checkAndPersistFoundValue(allFound[i])
		}

		allHrefs := s.Filter.ApplyFilter(body, HREF_REGEXP)
		for k := 0; k < len(allHrefs); k++ {
			parsedHref := s.Filter.HrefParser(allHrefs[k])
			s.checkAndLog(parsedHref)
		}
	})

	s.collector.Visit(s.url)
}

func (s *Scrapper) checkAndLog(val string) {

	if s.Filter.CheckBannedExtensions(val) {
		return
	}

	doesStartWithSlash := s.checkPrefix(val, "/")

	if doesStartWithSlash {

		wasCheckedBefore := s.checkHasControlledBefore(val)

		if !wasCheckedBefore {
			s.EndpointsChecked[val] = true
			s.collector.Visit(s.url + val)
		}

	} else {
		doesStartWithHttp := s.checkPrefix(val, "http")
		if doesStartWithHttp {

			s.collector.Visit(val)
		}
	}
}

func (s *Scrapper) checkHasControlledBefore(endpoint string) bool {
	urlWithEndpoint := s.url + endpoint
	//fmt.Println(urlWithEndpoint)
	_, ok := s.EndpointsChecked[urlWithEndpoint]
	return ok
}

func (s *Scrapper) checkPrefix(endpoint, prefix string) bool {
	return strings.HasPrefix(endpoint, prefix)
}

func (s *Scrapper) checkAndPersistFoundValue(value []byte) {
	valString := string(value)
	_, doWeHave := s.FoundPool[valString]
	if !doWeHave {
		fmt.Println("We found a matched string ; ", string(valString))
		s.FoundPool[valString] = true
		mailWithDownLine := append(value, 10)
		s.File.Write(mailWithDownLine)
	}
}

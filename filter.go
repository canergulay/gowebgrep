package gowebgrep

import (
	"fmt"
	"os"
	"regexp"
)

type Filter struct {
	RegExp       string
	SecondFilter *string // WILL NOT BE USED IF NOT SPECIFIED
}

const (
	DEFAULT_REGEXP = `[a-zA-Z0-9+._-]+@[a-zA-Z0-9._-]+\.[a-zA-Z0-9_-]+`
	HREF_REGEXP    = `href=['"]([^'"]+?)['"]`
)

const (
	JPG = "jpg"
	PNG = "png"
	PDF = "pdf"
)

func NewFilter() *Filter {

	args := os.Args[1:]

	switch len(args) {
	case 0:
		fmt.Println(" YOU SHOULD SPECIFY A URL AS COMMAND LINE ARGUMENT!")
		os.Exit(1)
	case 1:
		return &Filter{
			RegExp:       DEFAULT_REGEXP,
			SecondFilter: nil,
		}
	case 2:
		return &Filter{
			RegExp:       args[1],
			SecondFilter: nil,
		}
	}
	return &Filter{
		RegExp:       args[1],
		SecondFilter: &args[2],
	}
}

func (f *Filter) ApplyFilter(data []byte, rxp string) [][]byte {
	all := regexpHelper(data, rxp)
	if f.SecondFilter == nil {
		return all // IF THERE ISN'T A SECONDFILTER SPECIFIED, WE WILL RETURN EVERYTHING WE FOUND FOR THE GIVEN REGEXP.
	}

	// TODO : IMPLEMENT SECOND FILTER LATER ON
	return all
}

func regexpHelper(data []byte, rxp string) [][]byte {
	re := regexp.MustCompile(rxp)
	allMatches := re.FindAll(data, 1000) // ASSUMING THAT THERE COULD HAVE BEEN 1000 EMAILS IN A WEB PAGE AT MOST.
	return allMatches
}

func (f *Filter) CheckBannedExtensions(data string) bool {
	if len(data) < 3 {
		return false
	}
	extension := data[len(data)-3:]
	return extension == JPG || extension == PNG || extension == PDF
}

func (f *Filter) HrefParser(data []byte) string {
	data = data[:len(data)-1]
	data = data[6:]
	return string(data)
}

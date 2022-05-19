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

func (f *Filter) ApplyFilter(data []byte) [][]byte {
	all := regexpHelper(data, f.RegExp)
	if f.SecondFilter == nil {
		return all // IF THERE ISN'T A SECONDFILTER SPECIFIED, WE WILL RETURN EVERYTHING WE FOUND FOR THE GIVEN REGEXP.
	}

	// TODO : IMPLEMENT SECOND FILTER LATER ON
	return all
}

func regexpHelper(data []byte, rxp string) [][]byte {
	re := regexp.MustCompile(`[a-zA-Z0-9+._-]+@[a-zA-Z0-9._-]+\.[a-zA-Z0-9_-]+`)
	allMatches := re.FindAll(data, 1000) // ASSUMING THAT THERE COULD HAVE BEEN 1000 EMAILS IN A WEB PAGE AT MOST.
	return allMatches
}

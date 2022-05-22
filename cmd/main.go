package main

import (
	"github.com/canergulay/gowebgrep"
)

func main() {
	filter := gowebgrep.NewFilter()
	timer := gowebgrep.CreateTimer(750)
	scrapper := gowebgrep.InitializeScrapper(filter, timer, nil)
	scrapper.StartScrapping("/")
}

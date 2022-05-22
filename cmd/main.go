package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/canergulay/gowebgrep"
)

func main() {
	filter := gowebgrep.NewFilter()
	timer := gowebgrep.CreateTimer(750)
	scrapper := gowebgrep.InitializeScrapper(filter, timer, nil)
	scrapper.StartScrapping("/")
	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')
	fmt.Print("-> ", text)
}

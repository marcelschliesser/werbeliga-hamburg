package main

import (
	"fmt"
	"log"

	"github.com/marcelschliesser/werbeliga-hamburg/crawler"
)

const url string = "https://werbeliga.de/de/Spielplan,%20Tabelle%20&%20Torsch%C3%BCtzen"

func main() {
	matches, err := crawler.ParseMatchResults(url)
	if err != nil {
		log.Println(err.Error())
	}
	for i, m := range matches {
		fmt.Println(i, m)
	}
}

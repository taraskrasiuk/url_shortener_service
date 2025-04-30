package main

import (
	"flag"
	"fmt"
	"log"
	"taraskrasiuk/url_shortener_service/internal/shortener"
)

func main() {
	flag.Parse()
	shortUrl, err := shortener.NewShortLinker(10, "http", "localhost").Create(flag.Arg(0))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(shortUrl)
}

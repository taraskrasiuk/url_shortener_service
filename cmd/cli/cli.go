package main

import (
	"flag"
	"fmt"
	"log"
	"taraskrasiuk/url_shortener_service/internal/shortener"
	"taraskrasiuk/url_shortener_service/internal/storage"
)

func main() {
	flag.Parse()
	shortUrl, err := shortener.NewShortLinker(10).Create(flag.Arg(0))
	if err != nil {
		log.Fatal(err)
	}

	fStorage := storage.NewFileStorage("file_storage.db")
	fStorage.Write(flag.Arg(0), shortUrl)
	fmt.Println("the pair link and shorten version has been written.")
}

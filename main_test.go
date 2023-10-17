package main

import (
	"log"
	"os"
	"testing"

	"github.com/PuerkitoBio/goquery"
)

func BenchmarkParse(b *testing.B) {

	c, err := loadConfig("config.json")
	if err != nil {
		log.Fatal("Could not read configuration file: ", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {

		b.StopTimer()

		f, err := os.Open("pizza.html")
		if err != nil {
			log.Fatal(err)
		}

		doc, err := goquery.NewDocumentFromReader(f)
		if err != nil {
			log.Fatal(err)
		}

		b.StartTimer()
		_ = parseDoc(doc, c)

		f.Close()
	}
}

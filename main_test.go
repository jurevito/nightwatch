package main

import (
	"encoding/json"
	"log"
	"os"
	"testing"

	"github.com/PuerkitoBio/goquery"
)

func BenchmarkParse(b *testing.B) {

	// Read CSS selector configuration.
	data, err := os.ReadFile("config.json")
	if err != nil {
		log.Fatal(err)
	}

	var config Config

	if err := json.Unmarshal(data, &config); err != nil {
		log.Fatal(err)
	}

	// Read HTML file.
	file, err := os.Open("pizza.html")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	for i := 0; i < b.N; i++ {
		doc, err := goquery.NewDocumentFromReader(file)
		if err != nil {
			log.Fatal(err)
		}

		b.ResetTimer()
		_ = parseDoc(doc, &config)
	}
}

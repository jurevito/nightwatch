package main

import (
	"log"
	"os"
	"testing"

	"github.com/PuerkitoBio/goquery"
	"github.com/stretchr/testify/require"
)

func TestParse(t *testing.T) {
	c, err := loadConfig("config.json")
	if err != nil {
		log.Fatal("Could not read configuration file: ", err)
	}

	f, err := os.Open("pizza.html")
	if err != nil {
		log.Fatal("Could not read HTML file: ", err)
	}
	defer f.Close()

	doc, err := goquery.NewDocumentFromReader(f)
	if err != nil {
		log.Fatal("File cannot be parsed as HTML: ", err)
	}

	res := parseDoc(doc, c)

	// Check if number of parsed links is correct.
	require.Equal(t, 98, len(res.Organic))
	require.Equal(t, 28, len(res.Panel))
	require.Equal(t, 3, len(res.Local))

	expected := []int{5, 3, 6, 4}
	require.Equal(t, len(expected), len(res.Caurasel))

	for i := 0; i < len(res.Caurasel); i++ {
		require.Equal(t, expected[i], len(res.Caurasel[i]))
	}
}

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

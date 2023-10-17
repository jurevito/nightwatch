package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/PuerkitoBio/goquery"
)

type Link struct {
	Title string   `json:"title"`
	Url   string   `json:"url"`
	Type  RankType `json:"type"`
}

type Result struct {
	Organic  []Link   `json:"organic"`
	Caurasel [][]Link `json:"caurasel"`
	Panel    []Link   `json:"knowledge_panel"`
	Local    []Link   `json:"local"`
}

type Config struct {
	MainDiv string `json:"main_div"`

	OrganicDiv  string `json:"organic_div"`
	CauraselDiv string `json:"caurasel_div"`

	PhotoCaurasel string `json:"photo_caurasel"`

	FindCaurasel string `json:"find_caurasel"`
	FindTitle    string `json:"find_title"`

	RecipeCaurasel string `json:"recipe_caurasel"`
	RecipeTitle    string `json:"recipe_title"`

	VideoCaurasel string `json:"video_caurasel"`
	VideoTitle    string `json:"video_title"`

	LocalLinks string `json:"local_links"`
	LocalTitle string `json:"local_title"`

	PanelLinks string `json:"panel_links"`
}

func main() {

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

	doc, err := goquery.NewDocumentFromReader(file)
	if err != nil {
		log.Fatal(err)
	}

	res := parseDoc(doc, &config)

	// Save Result struct into JSON file.
	data, err = json.Marshal(res)
	if err != nil {
		log.Fatal(err)
	}

	output, err := os.Create("output.json")
	if err != nil {
		log.Fatal(err)
	}
	defer output.Close()

	_, err = output.Write(data)
	if err != nil {
		log.Fatal(err)
	}
}

func parseDoc(doc *goquery.Document, config *Config) Result {

	res := Result{}
	parent := doc.Find(config.MainDiv).First()

	// Organic.
	parent.Find(config.OrganicDiv).Each(func(index int, elem *goquery.Selection) {
		url, exist := elem.Find("a").First().Attr("href")
		title := elem.Find("h3").First().Text()

		if !exist {
			fmt.Println("Error parsing organic links.")
		}

		res.Organic = append(res.Organic, Link{
			Title: title,
			Url:   url,
			Type:  Organic,
		})
	})

	// Caurasel.
	parent.Find(config.CauraselDiv).Each(func(index int, elem *goquery.Selection) {

		local := false
		links := []Link{}

		// Photo Caurasel.
		elem.Find(config.PhotoCaurasel).Each(func(index int, elem *goquery.Selection) {
			src, _ := elem.Attr("data-src")
			title, _ := elem.Attr("alt")

			if len(title) != 0 && len(src) != 0 {
				links = append(links, Link{
					Title: title,
					Url:   src,
					Type:  Carousel,
				})
			}
		})

		// 'Find results on' Caurasel.
		elem.Find(config.FindCaurasel).Each(func(index int, elem *goquery.Selection) {
			url, exist := elem.Attr("href")
			title := elem.Find(config.FindTitle).First().Text()

			if !exist {
				fmt.Println("Error parsing 'Find results on' caurasel links.")
			}

			if len(title) != 0 && len(url) != 0 {
				links = append(links, Link{
					Title: title,
					Url:   url,
					Type:  Carousel,
				})
			}
		})

		// Recipe Caurasel.
		elem.Find(config.RecipeCaurasel).Each(func(index int, elem *goquery.Selection) {
			url, exist := elem.Attr("href")
			title := elem.Find(config.RecipeTitle).First().Text()

			if !exist {
				fmt.Println("Error parsing 'Recipes' caurasel links.")
			}

			if len(title) != 0 && len(url) != 0 {
				links = append(links, Link{
					Title: title,
					Url:   url,
					Type:  Carousel,
				})
			}
		})

		// Video Caurasel.
		elem.Find(config.VideoCaurasel).Each(func(index int, elem *goquery.Selection) {
			url, exist := elem.Attr("href")
			title := elem.Find(config.VideoTitle).First().Text()

			if !exist {
				fmt.Println("Error parsing 'Videos' caurasel links.")
			}

			links = append(links, Link{
				Title: title,
				Url:   url,
				Type:  Carousel,
			})
		})

		// Local.
		elem.Find(config.LocalLinks).Each(func(index int, elem *goquery.Selection) {
			title := elem.Find(config.LocalTitle).First().Text()
			local = true

			links = append(links, Link{
				Title: title,
				Url:   "",
				Type:  Local,
			})
		})

		if local {
			res.Local = links
		} else {
			res.Caurasel = append(res.Caurasel, links)
		}
	})

	// Knowledge Panel.
	panelLinks := parent.Find(config.PanelLinks)
	panelLinks = panelLinks.FilterFunction(func(i int, link *goquery.Selection) bool {

		// Filter out empty links.
		url, exist := link.Attr("href")
		return exist && len(url) > 1
	})

	panelLinks.Each(func(index int, elem *goquery.Selection) {
		link, exist := elem.Attr("href")
		title := elem.Text()

		if !exist {
			fmt.Println("Error parsing knowledge panel links.")
		}

		res.Panel = append(res.Panel, Link{
			Title: title,
			Url:   link,
			Type:  Panel,
		})
	})

	return res
}

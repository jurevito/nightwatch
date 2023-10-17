package main

import (
	"encoding/json"
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
	MainElem     string `json:"main_elem"`
	CauraselElem string `json:"caurasel_elem"`

	OrganicElem  string `json:"organic_elem"`
	OrganicLink  string `json:"organic_link"`
	OrganicTitle string `json:"organic_title"`

	PhotoImg string `json:"photo_img"`

	FindElem  string `json:"find_elem"`
	FindLink  string `json:"find_link"`
	FindTitle string `json:"find_title"`

	RecipeElem  string `json:"recipe_elem"`
	RecipeLink  string `json:"recipe_link"`
	RecipeTitle string `json:"recipe_title"`

	VideoElem  string `json:"video_elem"`
	VideoLink  string `json:"video_link"`
	VideoTitle string `json:"video_title"`

	LocalElem  string `json:"local_elem"`
	LocalLink  string `json:"local_link"`
	LocalTitle string `json:"local_title"`

	PanelElem  string `json:"panel_elem"`
	PanelLink  string `json:"panel_link"`
	PanelTitle string `json:"panel_title"`
}

func loadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var c Config

	if err := json.Unmarshal(data, &c); err != nil {
		return nil, err
	}

	return &c, nil
}

func saveJSON(path string, res *Result) error {
	data, err := json.Marshal(res)
	if err != nil {
		return err
	}

	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	if _, err = f.Write(data); err != nil {
		return err
	}

	return nil
}

func main() {

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

	if err := saveJSON("output.json", res); err != nil {
		log.Fatal("Could not save results into JSON file: ", err)
	}
}

func extract(ls, ts string, elem *goquery.Selection, t RankType) (*Link, bool) {
	url, exist := elem.Find(ls).AddSelection(elem).First().Attr("href")
	title := elem.Find(ts).AddSelection(elem).First().Text()

	return &Link{
		Title: title,
		Url:   url,
		Type:  t,
	}, exist
}

func isValid(link *Link) bool {
	return len(link.Url) > 1 && len(link.Title) != 0
}

func parseDoc(doc *goquery.Document, c *Config) *Result {

	res := Result{}
	parent := doc.Find(c.MainElem).First()

	// Organic.
	parent.Find(c.OrganicElem).Each(func(index int, elem *goquery.Selection) {
		link, _ := extract(c.OrganicLink, c.OrganicTitle, elem, Organic)

		if isValid(link) {
			res.Organic = append(res.Organic, *link)
		}
	})

	// Knowledge Panel.
	parent.Find(c.PanelElem).Each(func(index int, elem *goquery.Selection) {
		link, _ := extract(c.PanelLink, c.PanelTitle, elem, Panel)

		if isValid(link) {
			res.Panel = append(res.Panel, *link)
		}
	})

	// Caurasel.
	parent.Find(c.CauraselElem).Each(func(index int, elem *goquery.Selection) {

		isLocal := false
		links := []Link{}

		// Photo Caurasel.
		elem.Find(c.PhotoImg).Each(func(index int, elem *goquery.Selection) {
			url, _ := elem.Attr("data-src")
			title, _ := elem.Attr("alt")

			link := &Link{
				Title: title,
				Url:   url,
				Type:  Carousel,
			}

			if isValid(link) {
				links = append(links, *link)
			}
		})

		// 'Find results on' Caurasel.
		elem.Find(c.FindElem).Each(func(index int, elem *goquery.Selection) {
			link, _ := extract(c.FindLink, c.FindTitle, elem, Carousel)

			if isValid(link) {
				links = append(links, *link)
			}
		})

		// Recipe Caurasel.
		elem.Find(c.RecipeElem).Each(func(index int, elem *goquery.Selection) {
			link, _ := extract(c.RecipeLink, c.RecipeTitle, elem, Carousel)

			if isValid(link) {
				links = append(links, *link)
			}
		})

		// Video Caurasel.
		elem.Find(c.VideoElem).Each(func(index int, elem *goquery.Selection) {
			link, _ := extract(c.VideoLink, c.VideoTitle, elem, Carousel)

			if isValid(link) {
				links = append(links, *link)
			}
		})

		// Local.
		elem.Find(c.LocalElem).Each(func(index int, elem *goquery.Selection) {
			link, _ := extract(c.LocalLink, c.LocalTitle, elem, Local)
			isLocal = true

			if len(link.Title) != 0 {
				links = append(links, *link)
			}
		})

		if isLocal {
			res.Local = links
		} else {
			res.Caurasel = append(res.Caurasel, links)
		}
	})

	return &res
}

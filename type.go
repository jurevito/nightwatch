package main

import "encoding/json"

type RankType int

const (
	Organic RankType = iota
	Local
	Carousel
	Panel
	Snippet
)

func (t RankType) String() string {
	switch t {
	case Organic:
		return "organic"
	case Local:
		return "local"
	case Carousel:
		return "carousel"
	case Panel:
		return "knowledge_panel"
	case Snippet:
		return "featured_snippet"
	}
	return "unknown"
}

func (t RankType) MarshalJSON() ([]byte, error) {
	var s string

	switch t {
	case Organic:
		s = "organic"
	case Local:
		s = "local"
	case Carousel:
		s = "carousel"
	case Panel:
		s = "knowledge_panel"
	case Snippet:
		s = "featured_snippet"
	}

	return json.Marshal(s)
}

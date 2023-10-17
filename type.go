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

func (s RankType) String() string {
	switch s {
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

func (s RankType) MarshalJSON() ([]byte, error) {
	var rType string

	switch s {
	case Organic:
		rType = "organic"
	case Local:
		rType = "local"
	case Carousel:
		rType = "carousel"
	case Panel:
		rType = "knowledge_panel"
	case Snippet:
		rType = "featured_snippet"
	}

	return json.Marshal(rType)
}

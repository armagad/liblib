package search

import (
	"fmt"
	"net/http"
	"time"
)

type corpus struct{}

var data corpus

func init() {
	data = corpus{}
}

func HandleHttp() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query().Get("q")
		for _, result := range algorithm(query) {
			fmt.Fprintf(w,
				"<h6>%s</h6><h1>%s</h1><p>%s...%s</p>",
				result.URL,
				result.Title,
				result.Text,
				result.Date,
			)
		}
	}
}

func AddItems(items interface{}) {
}

func PurgeItems(items interface{}) {
}

func algorithm(q string) []SearchResult {
	return []SearchResult{}
}

type SearchResult struct {
	Matches []string
	Title   string
	Date    time.Time
	URL     string
	Text    string
}

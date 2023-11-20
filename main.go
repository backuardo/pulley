package main

import (
	"encoding/json"
	"fmt"
	"index/suffixarray"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

func main() {
	searcher := Searcher{}
	err := searcher.Load("completeworks.txt")
	if err != nil {
		log.Fatal(err)
	}

	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", fs)

	http.HandleFunc("/search", handleSearch(searcher))

	port := os.Getenv("PORT")
	if port == "" {
		port = "3001"
	}

	fmt.Printf("shakesearch available at http://localhost:%s...", port)
	err = http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
	if err != nil {
		log.Fatal(err)
	}
}

type Searcher struct {
	CompleteWorks string
	LowerCompleteWorks string
	SuffixArray   *suffixarray.Index
}

func handleSearch(searcher Searcher) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query().Get("q")
		if query == "" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("missing query parameter"))
			return
		}

		offsetStr := r.URL.Query().Get("offset")
		limitStr := r.URL.Query().Get("limit")

		offset, err := strconv.Atoi(offsetStr)
		if err != nil {
			offset = 0
		}

		limit, err := strconv.Atoi(limitStr)
		if err != nil || limit <= 0 {
			limit = 20
		}

		results := searcher.Search(query, offset, limit)

		jsonResponse, err := json.Marshal(results)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("error encoding results"))
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonResponse)
	}
}

func (s *Searcher) Load(filename string) error {
	dat, err := ioutil.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("Load: %w", err)
	}
	s.CompleteWorks = string(dat)
	s.LowerCompleteWorks = strings.ToLower(s.CompleteWorks)
	s.SuffixArray = suffixarray.New([]byte(s.LowerCompleteWorks))
	return nil
}

func (s *Searcher) Search(query string, offset, limit int) []string {
	lowerQuery := strings.ToLower(query)
	idxs := s.SuffixArray.Lookup([]byte(lowerQuery), -1)

	results := []string{}
	for i, idx := range idxs {
		if i < offset {
			continue
		}
		if len(results) >= limit {
			break
		}
		start := max(0, idx-250)
		end := min(len(s.CompleteWorks), idx+250)
		results = append(results, s.CompleteWorks[start:end])
	}

	return results
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if b < a {
		return a
	}
	return b
}
package models

import (
	"encoding/json"
	"os"
	"strings"
)

// saveToFile saves the Articles slice to the JSON file.
func (m *ArticleModel) saveToFile() error {
	file, err := os.Create("internal/data/articles.json")
	if err != nil {
		return err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			panic(err)
		}
	}(file)

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ") // Optional: to format the JSON for better readability
	if err := encoder.Encode(m.Articles); err != nil {
		return err
	}

	return nil
}

// loadFromFile loads articles from the JSON file into the Articles slice.
func (m *ArticleModel) loadFromFile() error {
	file, err := os.Open("internal/data/articles.json")
	if err != nil {
		return err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			panic(err)
		}
	}(file)

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&m.Articles); err != nil {
		return err
	}

	return nil
}

// Helper function to match search query
func matchSearchQuery(article Article, query string) bool {
	// Convert both title and query to lowercase for case-insensitive search
	title := strings.ToLower(article.Title)
	content := strings.ToLower(article.Content)
	query = strings.ToLower(query)

	// Check if query is in title or content
	return strings.Contains(title, query) || strings.Contains(content, query)
}

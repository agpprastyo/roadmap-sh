package models

import (
	"errors"
	"sort"

	"sync"
	"time"
)

// Article represents a single article.
type Article struct {
	ID        int        `json:"id"`
	Title     string     `json:"title"`
	Content   string     `json:"content"`
	Published bool       `json:"published"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
	DeleteAt  *time.Time `json:"delete_at"`
}

type SortBy struct {
	TitleAsc      string `json:"title_asc"`
	TitleDesc     string `json:"title_desc"`
	CreatedAtAsc  string `json:"created_at_asc"`
	CreatedAtDesc string `json:"created_at_desc"`
}

// ArticleModel holds the data for your articles.
type ArticleModel struct {
	Articles []Article
	SortBy   SortBy
	mu       sync.Mutex // Mutex to protect access to Articles
}

// GetAll retrieves articles with optional sorting and searching
func (m *ArticleModel) GetAll(sortBy string, searchQuery string) ([]Article, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Load articles from JSON file
	if err := m.loadFromFile(); err != nil {
		return nil, err
	}

	var publishedArticles []Article

	// Filter articles
	for _, article := range m.Articles {
		if article.Published && article.DeleteAt == nil {
			// If search query is provided, perform search
			if searchQuery != "" {
				if !matchSearchQuery(article, searchQuery) {
					continue
				}
			}
			publishedArticles = append(publishedArticles, article)
		}
	}

	// Sort articles based on the sortBy parameter
	switch sortBy {
	case "title_asc":
		sort.Slice(publishedArticles, func(i, j int) bool {
			return publishedArticles[i].Title < publishedArticles[j].Title
		})
	case "title_desc":
		sort.Slice(publishedArticles, func(i, j int) bool {
			return publishedArticles[i].Title > publishedArticles[j].Title
		})
	case "created_at_asc":
		sort.Slice(publishedArticles, func(i, j int) bool {
			return publishedArticles[i].CreatedAt.Before(publishedArticles[j].CreatedAt)
		})
	case "created_at_desc":
		sort.Slice(publishedArticles, func(i, j int) bool {
			return publishedArticles[i].CreatedAt.After(publishedArticles[j].CreatedAt)
		})
	default:
		// Default sorting (e.g., by creation date, latest first)
		sort.Slice(publishedArticles, func(i, j int) bool {
			return publishedArticles[i].CreatedAt.After(publishedArticles[j].CreatedAt)
		})
	}

	return publishedArticles, nil
}

// GetByID retrieves an article by its ID, ensuring it is published, published before now, and not deleted.
func (m *ArticleModel) GetByID(id int) (*Article, error) {
	m.mu.Lock()         // Lock the mutex to prevent concurrent access
	defer m.mu.Unlock() // Ensure the mutex is unlocked after the function returns

	// Load articles from JSON file
	if err := m.loadFromFile(); err != nil {
		return nil, err
	}

	for _, article := range m.Articles {
		if article.ID == id && article.Published && article.DeleteAt == nil {
			return &article, nil
		}
	}

	// Return nil if no matching article is found
	return nil, nil
}

// GetByIDAdmin read article by ID for admin access
func (m *ArticleModel) GetByIDAdmin(id int) (*Article, error) {
	m.mu.Lock()         // Lock the mutex to prevent concurrent access
	defer m.mu.Unlock() // Ensure the mutex is unlocked after the function returns

	// Load articles from JSON file
	if err := m.loadFromFile(); err != nil {
		return nil, err
	}

	for _, article := range m.Articles {
		if article.ID == id {
			return &article, nil
		}
	}

	// Return nil if no matching article is found
	return nil, nil
}

// GetAllAdmin retrieves all articles for admin access.
func (m *ArticleModel) GetAllAdmin(sortBy string, searchQuery string) ([]Article, error) {
	m.mu.Lock()         // Lock the mutex to prevent concurrent access
	defer m.mu.Unlock() // Ensure the mutex is unlocked after the function returns

	// Load articles from JSON file
	if err := m.loadFromFile(); err != nil {
		return nil, err
	}

	var allArticles []Article

	// Filter articles
	for _, article := range m.Articles {
		if searchQuery != "" {
			if !matchSearchQuery(article, searchQuery) {
				continue
			}
		}
		allArticles = append(allArticles, article)

	}

	// Sort articles based on the sortBy parameter if provided
	switch sortBy {
	case "title-asc": // Sort by title ascending (A-Z)
		sort.Slice(allArticles, func(i, j int) bool {
			return allArticles[i].Title < allArticles[j].Title
		})
	case "title-desc": // Sort by title descending (Z-A)
		sort.Slice(allArticles, func(i, j int) bool {
			return allArticles[i].Title > allArticles[j].Title
		})
	case "created_at-asc": // Sort by created_at ascending (oldest first)
		sort.Slice(allArticles, func(i, j int) bool {
			return allArticles[i].CreatedAt.Before(allArticles[j].CreatedAt)
		})
	case "created_at-desc": // Sort by created_at descending (latest first)
		sort.Slice(allArticles, func(i, j int) bool {
			return allArticles[i].CreatedAt.After(allArticles[j].CreatedAt)
		})
	}

	return allArticles, nil
}

// CreateArticle creates a new article and appends it to the Articles slice.
func (m *ArticleModel) CreateArticle(article Article) (Article, string, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Load existing articles from the JSON file to get the latest ID
	if err := m.loadFromFile(); err != nil {
		return Article{}, "", err
	}

	// Validate title and content
	if article.Title == "" || article.Content == "" {
		return Article{}, "", errors.New("title and content must be provided")
	}

	// Set the ID to be one more than the highest existing ID
	var maxID int
	for _, a := range m.Articles {
		if a.ID > maxID {
			maxID = a.ID
		}
	}
	article.ID = maxID + 1 // Set to the next available ID

	// Set CreatedAt to the current time
	now := time.Now()
	article.CreatedAt = now

	// Handle Published status
	// This will respect the Published value from the input
	// If not specified, it will default to false (Go's zero value for bool)

	// Set default values for UpdatedAt and DeleteAt
	article.UpdatedAt = nil
	article.DeleteAt = nil

	// Append the new article to the Articles slice
	m.Articles = append(m.Articles, article)

	// Save the updated articles back to the JSON file
	if err := m.saveToFile(); err != nil {
		return Article{}, "", err
	}

	// Prepare message based on published status
	message := "Article created successfully,"
	if article.Published {
		message += " it is published."
	} else {
		message += " but not published."
	}

	return article, message, nil
}

// UpdateArticle updates an existing article by ID, only modifying the fields that are provided
func (m *ArticleModel) UpdateArticle(id int, updatedArticle Article) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Load existing articles from the JSON file
	if err := m.loadFromFile(); err != nil {
		return err
	}

	// Find the article by ID
	for i, article := range m.Articles {
		if article.ID == id {
			// Update title if not empty
			if updatedArticle.Title != "" {
				m.Articles[i].Title = updatedArticle.Title
			}

			// Update content if not empty
			if updatedArticle.Content != "" {
				m.Articles[i].Content = updatedArticle.Content
			}

			// Always update published status (even if false)
			// This allows changing from true to false
			m.Articles[i].Published = updatedArticle.Published

			// Update the UpdatedAt timestamp
			now := time.Now()
			m.Articles[i].UpdatedAt = &now

			// Save the updated articles back to the file
			return m.saveToFile()
		}
	}

	// If no article was found with the given ID, return an error
	return errors.New("article not found")
}

// DeleteArticle marks an article as deleted by ID.
func (m *ArticleModel) DeleteArticle(id int) error {
	m.mu.Lock()         // Lock the mutex to prevent concurrent access
	defer m.mu.Unlock() // Ensure the mutex is unlocked after the function returns

	// Load existing articles from the JSON file
	if err := m.loadFromFile(); err != nil {
		return err
	}

	// Find the article by ID
	for i, article := range m.Articles {
		if article.ID == id {
			// Mark the article as deleted
			now := time.Now()
			m.Articles[i].DeleteAt = &now // Set DeleteAt to now
			// Save the updated articles back to the file
			return m.saveToFile()
		}
	}

	// If no article was found with the given ID, return an error
	return errors.New("article not found")
}

// RestoreArticle restores a deleted article by ID.
func (m *ArticleModel) RestoreArticle(id int) error {
	m.mu.Lock()         // Lock the mutex to prevent concurrent access
	defer m.mu.Unlock() // Ensure the mutex is unlocked after the function returns

	// Load existing articles from the JSON file
	if err := m.loadFromFile(); err != nil {
		return err
	}

	// Find the article by ID
	for i, article := range m.Articles {
		if article.ID == id {
			// Mark the article as not deleted
			m.Articles[i].DeleteAt = nil // Set DeleteAt to nil
			// Save the updated articles back to the file
			return m.saveToFile()
		}
	}

	// If no article was found with the given ID, return an error
	return errors.New("article not found")
}

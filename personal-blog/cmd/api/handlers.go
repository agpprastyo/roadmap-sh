package main

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"golang.org/x/crypto/bcrypt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"personal-blog/internal/models"
	"personal-blog/internal/request"
	"personal-blog/internal/response"
	"personal-blog/internal/sessions"
	"strconv"
)

func (app *application) status(w http.ResponseWriter, r *http.Request) {
	data := map[string]string{
		"Status": "OK",
	}

	err := response.JSON(w, http.StatusOK, data)
	if err != nil {
		app.serverError(w, r, err)
	}
}

func (app *application) protected(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte("This is a protected handler"))
	if err != nil {
		app.serverError(w, r, err)
	}
}

func (app *application) signIn(w http.ResponseWriter, r *http.Request) {
	// Log the incoming request details
	log.Printf("Sign-in attempt from IP: %s", r.RemoteAddr)

	// Get the username and password from the Authorization header
	username, password, ok := r.BasicAuth()
	if !ok {
		log.Printf("Failed sign-in: No credentials provided")
		http.Error(w, "Unauthorized: No credentials", http.StatusUnauthorized)
		return
	}

	// Log the username (but never log the password)
	log.Printf("Sign-in attempt for username: %s", username)

	// Validate the username
	if username != app.config.basicAuth.username {
		log.Printf("Failed sign-in: Incorrect username. Expected: %s, Got: %s",
			app.config.basicAuth.username, username)
		http.Error(w, "Unauthorized: Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Compare passwords
	passwordMatch := bcrypt.CompareHashAndPassword(
		[]byte(app.config.basicAuth.hashedPassword),
		[]byte(password),
	) == nil

	if !passwordMatch {
		log.Printf("Failed sign-in: Incorrect password for username: %s", username)
		http.Error(w, "Unauthorized: Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Authentication successful, create a session
	_, err := sessions.CreateSession(w, r, username, app.config.cookie.secretKey)
	if err != nil {
		log.Printf("Error creating session: %v", err)
		app.serverError(w, r, err)
		return
	}

	// Log successful authentication
	log.Printf("Successful sign-in for user: %s", username)

	// Prepare response
	responseData := map[string]interface{}{
		"message": "Sign in successful",
	}

	// Set CORS headers explicitly
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
	w.Header().Set("Access-Control-Allow-Credentials", "true")

	// Return success response
	err = response.JSON(w, http.StatusOK, responseData)
	if err != nil {
		log.Printf("Error sending JSON response: %v", err)
		app.serverError(w, r, err)
	}
}

func (app *application) signOut(w http.ResponseWriter, r *http.Request) {
	// Invalidate the session, clearing the cookie
	if err := sessions.DestroySession(w, r); err != nil {
		app.serverError(w, r, err) // Handle session destruction error
		return
	}

	// return a success response
	err := response.JSON(w, http.StatusOK, map[string]string{"message": "Sign out successful"})
	if err != nil {
		app.serverError(w, r, err)
	}
}

// Handler for reading articles
func (app *application) readArticles(w http.ResponseWriter, r *http.Request) {
	// Define constants
	const defaultPage = 1
	const defaultPageSize = 10

	// Validate page sizes
	validPageSizes := map[int]bool{
		10: true,
		20: true,
		50: true,
	}

	// Determine the active sort option with priority
	var activeSortOption string
	switch {
	case r.URL.Query().Get("title_asc") == "true":
		activeSortOption = "title_asc"
	case r.URL.Query().Get("title_desc") == "true":
		activeSortOption = "title_desc"
	case r.URL.Query().Get("created_at_asc") == "true":
		activeSortOption = "created_at_asc"
	case r.URL.Query().Get("created_at_desc") == "true":
		activeSortOption = "created_at_desc"
	default:
		// No sorting specified, use default (latest first)
		activeSortOption = "created_at_desc"
	}

	// Parse query parameters
	pageStr := r.URL.Query().Get("page")
	pageSizeStr := r.URL.Query().Get("page_size")
	searchQuery := r.URL.Query().Get("search")

	// Validate and convert page
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = defaultPage
	}

	// Validate and convert page size
	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || !validPageSizes[pageSize] {
		pageSize = defaultPageSize
	}

	// Read all articles with optional search and sorting
	articles, err := app.articleModel.GetAll(activeSortOption, searchQuery)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	// Pagination logic
	start := (page - 1) * pageSize
	end := start + pageSize

	// Handle pagination edge cases
	if start > len(articles) {
		articles = app.articleModel.Articles
	} else if end > len(articles) {
		end = len(articles)
	}

	// Slice the articles for the current page
	paginatedArticles := articles[start:end]

	// Prepare response metadata
	responseData := map[string]interface{}{
		"articles":  paginatedArticles,
		"page":      page,
		"page_size": pageSize,
		"total":     len(articles),
		"sort_by":   activeSortOption,
		"search":    searchQuery,
		"next_url":  nil, // Initialize as nil
	}

	// Generate next URL if there are more pages
	if end < len(articles) {
		nextPage := page + 1
		nextURL := fmt.Sprintf("%s/api/v1/articles?page=%d&page_size=%d&%s=true&search=%s",
			r.Host,
			nextPage,
			pageSize,
			activeSortOption,
			url.QueryEscape(searchQuery))
		responseData["next_url"] = nextURL
	}

	// Send the response
	err = response.JSON(w, http.StatusOK, responseData)
	if err != nil {
		app.serverError(w, r, err)
	}
}

func (app *application) readArticleByID(w http.ResponseWriter, r *http.Request) {
	// Extract the article ID from the URL
	idStr := chi.URLParam(r, "id")

	// Convert the ID from string to int
	id, err := strconv.Atoi(idStr)
	if err != nil {
		// If the conversion fails, return a 400 Bad Request error
		app.badRequest(w, r, err)
		return
	}

	// Fetch the article by ID
	article, err := app.articleModel.GetByID(id)
	if err != nil {
		// If there was an error retrieving the article, return a 500 Internal Server Error
		app.serverError(w, r, err)
		return
	}

	// Check if the article was found
	if article == nil {
		// If no article is found, return a 404 Not Found error
		app.notFound(w, r)
		return
	}

	// Return the article as JSON
	err = response.JSON(w, http.StatusOK, article)
	if err != nil {
		// If there was an error marshaling the article to JSON, return a 500 Internal Server Error
		app.serverError(w, r, err)
		return
	}
}

func (app *application) readArticleByIDAdmin(w http.ResponseWriter, r *http.Request) {
	// Extract the article ID from the URL
	idStr := chi.URLParam(r, "id")

	// Convert the ID from string to int
	id, err := strconv.Atoi(idStr)
	if err != nil {
		// If the conversion fails, return a 400 Bad Request error
		http.Error(w, "Invalid article ID", http.StatusBadRequest)
		return
	}

	// Fetch the article by ID
	article, err := app.articleModel.GetByIDAdmin(id)
	if err != nil {
		// If there was an error retrieving the article, return a 500 Internal Server Error
		app.serverError(w, r, err)
		return
	}

	// Check if the article was found
	if article == nil {
		// If no article is found, return a 404 Not Found error
		app.notFound(w, r)
		return
	}

	// Return the article as JSON
	err = response.JSON(w, http.StatusOK, article)
	if err != nil {
		// If there was an error marshaling the article to JSON, return a 500 Internal Server Error
		app.serverError(w, r, err)
		return
	}
}

func (app *application) readArticlesAll(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters for pagination and sorting
	pageStr := r.URL.Query().Get("page")
	pageSizeStr := r.URL.Query().Get("page_size")
	sortBy := r.URL.Query().Get("sort_by")
	search := r.URL.Query().Get("search")

	// Set default pagination values if not provided
	const defaultPage = 1
	const defaultPageSize = 10
	page, err := strconv.Atoi(pageStr)

	if err != nil || page < 1 {
		page = defaultPage
	}

	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || (pageSize != 10 && pageSize != 20 && pageSize != 50) {
		pageSize = defaultPageSize
	}

	// Read all articles
	articles, err := app.articleModel.GetAllAdmin(sortBy, search)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	// Pagination logic
	start := (page - 1) * pageSize
	end := start + pageSize

	if start > len(articles) {
		// If the start index is greater than the length, return an empty slice
		articles = []models.Article{}
	} else if end > len(articles) {
		end = len(articles) // Make sure we do not exceed length
	}

	paginatedArticles := articles[start:end]

	// Prepare the response metadata
	responseData := map[string]interface{}{
		"articles":  paginatedArticles,
		"page":      page,
		"page_size": pageSize,
		"total":     len(articles),
		"next_url":  "", // URL for the next page if applicable
		"sort_by":   sortBy,
	}

	// Generate next URL if there are more pages
	if end < len(articles) { // Check if there is a next page
		nextPage := page + 1
		responseData["next_url"] = fmt.Sprintf("%s/api/v1/admin?page=%d&page_size=%d&sort_by=%s", r.Host, nextPage, pageSize, sortBy)
	} else {
		responseData["next_url"] = nil
	}

	// Send the response
	err = response.JSON(w, http.StatusOK, responseData)
	if err != nil {
		app.serverError(w, r, err)
	}
}

func (app *application) createArticle(w http.ResponseWriter, r *http.Request) {
	// Define a variable to hold the incoming article data
	var article models.Article

	// Use the custom DecodeJSONStrict method
	err := request.DecodeJSONStrict(w, r, &article)
	if err != nil {
		// Use a method to handle the error appropriately
		app.badRequest(w, r, err)
		return
	}

	// Validate required fields
	if article.Title == "" || article.Content == "" {
		//app.badRequest(w, r, http.StatusBadRequest, "Title and Content must be provided")
		app.badRequest(w, r, fmt.Errorf("title and content must be provided"))
		return
	}

	// Create the article using the ArticleModel's CreateArticle method
	createdArticle, message, err := app.articleModel.CreateArticle(article)
	if err != nil {
		// Handle errors from creating the article
		app.serverError(w, r, err)
		return
	}

	// Respond with the created article and a 201 Created status
	responseData := map[string]interface{}{
		"message": message,
		"article": createdArticle,
	}

	err = response.JSON(w, http.StatusCreated, responseData)
	if err != nil {
		app.serverError(w, r, err)
	}
}

func (app *application) updateArticle(w http.ResponseWriter, r *http.Request) {
	// Extract the article ID from the URL
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid article ID", http.StatusBadRequest)
		return
	}

	// Decode the incoming JSON request to an Article struct
	var updateData models.Article
	err = json.NewDecoder(r.Body).Decode(&updateData)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Call the UpdateArticle function from the ArticleModel
	err = app.articleModel.UpdateArticle(id, updateData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Respond with a success message
	successResponse := map[string]string{"message": "Article updated successfully"}

	err = response.JSON(w, http.StatusOK, successResponse)
	if err != nil {
		app.serverError(w, r, err)
	}
}

// deleteArticle marks an article as deleted by ID.
func (app *application) deleteArticle(w http.ResponseWriter, r *http.Request) {
	// Extract the article ID from the URL
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr) // Convert the string ID to an integer
	if err != nil {
		http.Error(w, "Invalid article ID", http.StatusBadRequest)
		return
	}

	// Call the DeleteArticle function from the ArticleModel
	err = app.articleModel.DeleteArticle(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Respond with a success message
	successResponse := map[string]string{"message": "Article deleted successfully"}

	err = response.JSON(w, http.StatusOK, successResponse)
	if err != nil {
		app.serverError(w, r, err)
	}
}

// restoreArticle marks an article as restored by ID.
func (app *application) restoreArticle(w http.ResponseWriter, r *http.Request) {
	// Extract the article ID from the URL
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr) // Convert the string ID to an integer
	if err != nil {
		http.Error(w, "Invalid article ID", http.StatusBadRequest)
		return
	}

	// Call the RestoreArticle function from the ArticleModel
	err = app.articleModel.RestoreArticle(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Respond with a success message
	successResponse := map[string]string{"message": "Article restored successfully"}

	err = response.JSON(w, http.StatusOK, successResponse)
	if err != nil {
		app.serverError(w, r, err)
	}
}

// serveOpenAPIDocs handles the request to serve the OpenAPI documentation.
func (app *application) serveOpenAPIDocs(w http.ResponseWriter, r *http.Request) {
	// Read the OpenAPI JSON file
	data, err := ioutil.ReadFile("../docs/personal-blog.openapi.json")
	if err != nil {
		http.Error(w, "Unable to read OpenAPI documentation", http.StatusInternalServerError)
		return
	}

	// Set the Content-Type header to application/json
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

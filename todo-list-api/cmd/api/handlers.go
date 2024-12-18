package main

import (
	"fmt"
	"github.com/agpprastyo/todo-list-api/internal/database"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
	"time"

	"github.com/agpprastyo/todo-list-api/internal/password"
	"github.com/agpprastyo/todo-list-api/internal/request"
	"github.com/agpprastyo/todo-list-api/internal/response"
	"github.com/agpprastyo/todo-list-api/internal/validator"
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

func (app *application) createUser(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Email     string              `json:"Email"`
		Password  string              `json:"Password"`
		Name      string              `json:"Name"`
		Validator validator.Validator `json:"-"`
	}

	err := request.DecodeJSON(w, r, &input)
	if err != nil {
		app.badRequest(w, r, err)
		return
	}

	_, found, err := app.db.GetUserByEmail(input.Email)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	input.Validator.CheckField(input.Email != "", "Email", "Email is required")
	input.Validator.CheckField(validator.Matches(input.Email, validator.RgxEmail), "Email", "Must be a valid email address")
	input.Validator.CheckField(!found, "Email", "Email is already in use")

	input.Validator.CheckField(input.Password != "", "Password", "Password is required")
	input.Validator.CheckField(len(input.Password) >= 8, "Password", "Password is too short")
	input.Validator.CheckField(len(input.Password) <= 72, "Password", "Password is too long")
	input.Validator.CheckField(validator.NotIn(input.Password, password.CommonPasswords...), "Password", "Password is too common")

	input.Validator.CheckField(input.Name != "", "Name", "Name is required")
	input.Validator.CheckField(len(input.Name) >= 3, "Name", "Name is too short")
	input.Validator.CheckField(len(input.Name) <= 100, "Name", "Name is too long")

	if input.Validator.HasErrors() {
		app.failedValidation(w, r, input.Validator)
		return
	}

	hashedPassword, err := password.Hash(input.Password)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	_, err = app.db.InsertUser(input.Email, hashedPassword, input.Name)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	data := map[string]interface{}{
		"Message": "User created",
		"User": map[string]string{
			"Email": input.Email,
			"Name":  input.Name,
		},
	}
	err = response.JSON(w, http.StatusCreated, data)
}

func (app *application) createAuthenticationToken(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Email     string              `json:"Email"`
		Password  string              `json:"Password"`
		Validator validator.Validator `json:"-"`
	}

	err := request.DecodeJSON(w, r, &input)
	if err != nil {
		app.badRequest(w, r, err)
		return
	}

	user, found, err := app.db.GetUserByEmail(input.Email)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	input.Validator.CheckField(input.Email != "", "Email", "Email is required")
	input.Validator.CheckField(found, "Email", "Email address could not be found")

	if found {
		passwordMatches, err := password.Matches(input.Password, user.HashedPassword)
		if err != nil {
			app.serverError(w, r, err)
			return
		}

		input.Validator.CheckField(input.Password != "", "Password", "Password is required")
		input.Validator.CheckField(passwordMatches, "Password", "Password is incorrect")
	}

	if input.Validator.HasErrors() {
		app.failedValidation(w, r, input.Validator)
		return
	}
	accessToken, accessExpiry, err := app.generateAccessToken(user)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	// Generate refresh token
	refreshToken, err := app.db.CreateRefreshToken(user.ID, 30*24*time.Hour) // 30 days
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	data := map[string]interface{}{
		"AccessToken":        accessToken,
		"AccessTokenExpiry":  accessExpiry.Format(time.RFC3339),
		"RefreshToken":       refreshToken.Token,
		"RefreshTokenExpiry": refreshToken.Expiry.Format(time.RFC3339),
	}

	err = response.JSON(w, http.StatusOK, data)
	if err != nil {
		app.serverError(w, r, err)
	}

}

// New refresh token handler
func (app *application) refreshToken(w http.ResponseWriter, r *http.Request) {
	var input struct {
		RefreshToken string              `json:"RefreshToken"`
		Validator    validator.Validator `json:"-"`
	}

	err := request.DecodeJSON(w, r, &input)
	if err != nil {
		app.badRequest(w, r, err)
		return
	}

	input.Validator.CheckField(input.RefreshToken != "", "RefreshToken", "Refresh token is required")
	if input.Validator.HasErrors() {
		app.failedValidation(w, r, input.Validator)
		return
	}

	// Validate refresh token
	refreshToken, err := app.db.GetRefreshToken(input.RefreshToken)
	if err != nil {
		fmt.Println("validate refresh token error")
		app.invalidRefreshToken(w, r)
		return
	}

	// Get user
	user, found, err := app.db.GetUser(refreshToken.UserID)
	if err != nil || !found {
		fmt.Println("get user error")
		app.invalidRefreshToken(w, r)
		return
	}

	// Generate new access token
	accessToken, accessExpiry, err := app.generateAccessToken(user)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	// Generate new refresh token
	newRefreshToken, err := app.db.CreateRefreshToken(user.ID, 30*24*time.Hour)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	// Revoke old refresh token
	err = app.db.RevokeRefreshToken(input.RefreshToken)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	data := map[string]interface{}{
		"AccessToken":        accessToken,
		"AccessTokenExpiry":  accessExpiry.Format(time.RFC3339),
		"RefreshToken":       newRefreshToken.Token,
		"RefreshTokenExpiry": newRefreshToken.Expiry.Format(time.RFC3339),
	}

	err = response.JSON(w, http.StatusOK, data)
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

// create a new todo
func (app *application) createTodo(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Title       string              `json:"Title"`
		Description string              `json:"Description"`
		Validator   validator.Validator `json:"-"`
	}

	err := request.DecodeJSON(w, r, &input)
	if err != nil {
		app.badRequest(w, r, err)
		return
	}

	input.Validator.CheckField(input.Title != "", "Title", "Title is required")
	input.Validator.CheckField(len(input.Title) <= 100, "Title", "Title is too long")

	if input.Validator.HasErrors() {
		app.failedValidation(w, r, input.Validator)
		return
	}

	authenticatedUser := contextGetAuthenticatedUser(r)
	// get user iD
	userID := authenticatedUser.ID

	_, err = app.db.InsertTodo(input.Title, input.Description, userID)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	data := map[string]interface{}{
		"Message": "Todo created",
		"Todo": map[string]interface{}{
			"Title":       input.Title,
			"Description": input.Description,
		},
	}
	err = response.JSON(w, http.StatusCreated, data)
}

// get all todos
func (app *application) getAllTodos(w http.ResponseWriter, r *http.Request) {
	authenticatedUser := contextGetAuthenticatedUser(r)
	userID := authenticatedUser.ID

	// Parse query parameters
	filters := database.TodoFilters{
		Search:    r.URL.Query().Get("search"),
		SortBy:    r.URL.Query().Get("sort_by"),
		SortOrder: r.URL.Query().Get("sort_order"),
		PageSize:  10, // Default page size
		Page:      1,  // Default page
	}

	// Parse page size
	if pageSizeStr := r.URL.Query().Get("page_size"); pageSizeStr != "" {
		if pageSize, err := strconv.Atoi(pageSizeStr); err == nil && pageSize > 0 {
			filters.PageSize = pageSize
		}
	}

	// Parse page number
	if pageStr := r.URL.Query().Get("page"); pageStr != "" {
		if page, err := strconv.Atoi(pageStr); err == nil && page > 0 {
			filters.Page = page
		}
	}

	todos, totalCount, err := app.db.GetTodos(userID, filters)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	// Calculate pagination metadata
	totalPages := (totalCount + filters.PageSize - 1) / filters.PageSize
	hasNextPage := filters.Page < totalPages
	hasPrevPage := filters.Page > 1

	data := map[string]interface{}{
		"Todos": todos,
		"Pagination": map[string]interface{}{
			"CurrentPage": filters.Page,
			"PageSize":    filters.PageSize,
			"TotalCount":  totalCount,
			"TotalPages":  totalPages,
			"HasNextPage": hasNextPage,
			"HasPrevPage": hasPrevPage,
		},
	}

	err = response.JSON(w, http.StatusOK, data)
	if err != nil {
		app.serverError(w, r, err)
	}
}

// delete a todo
func (app *application) deleteTodo(w http.ResponseWriter, r *http.Request) {
	authenticatedUser := contextGetAuthenticatedUser(r)

	userID := authenticatedUser.ID

	idStr := chi.URLParam(r, "id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		app.badRequest(w, r, err)
		return
	}

	todo, found, err := app.db.GetTodoByID(id)
	if !found {
		app.notFound(w, r)
		return
	}
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	if todo.UserID != userID {
		app.notFound(w, r)
		return
	}

	err = app.db.DeleteTodoByID(id)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	err = response.JSONWithHeaders(w, http.StatusNoContent, nil, nil)
	if err != nil {
		app.serverError(w, r, err)
	}
}

func (app *application) updateTodo(w http.ResponseWriter, r *http.Request) {
	authenticatedUser := contextGetAuthenticatedUser(r)
	userID := authenticatedUser.ID

	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		app.badRequest(w, r, err)
		return
	}

	todo, found, err := app.db.GetTodoByID(id)
	if !found {
		app.notFound(w, r)
		return
	}
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	if todo.UserID != userID {
		app.notFound(w, r)
		return
	}

	var input struct {
		Title       *string             `json:"Title,omitempty"`
		Description *string             `json:"Description,omitempty"`
		Validator   validator.Validator `json:"-"`
	}

	err = request.DecodeJSON(w, r, &input)
	if err != nil {
		app.badRequest(w, r, err)
		return
	}

	// Use existing values if not provided in the request
	title := todo.Title
	description := todo.Description

	if input.Title != nil {
		input.Validator.CheckField(*input.Title != "", "Title", "Title is required")
		input.Validator.CheckField(len(*input.Title) <= 100, "Title", "Title is too long")
		title = *input.Title
	}

	if input.Description != nil {
		description = *input.Description
	}

	if input.Validator.HasErrors() {
		app.failedValidation(w, r, input.Validator)
		return
	}

	err = app.db.UpdateTodo(id, title, description)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	data := map[string]interface{}{
		"Message": "Todo updated",
		"Todo": map[string]interface{}{
			"Title":       title,
			"Description": description,
		},
	}
	err = response.JSON(w, http.StatusOK, data)
	if err != nil {
		app.serverError(w, r, err)
	}
}

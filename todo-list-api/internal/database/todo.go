package database

import (
	"context"
	"strconv"
	"strings"
	"time"
)

type Todo struct {
	ID          int    `db:"id" json:"id"`
	CreatedAt   string `db:"created_at" json:"created_at"`
	Title       string `db:"title" json:"title"`
	Description string `db:"description" json:"description"`
	UserID      int    `db:"user_id" json:"user_id"`
}

type TodoFilters struct {
	Search    string
	SortBy    string
	SortOrder string
	Page      int
	PageSize  int
}

// InsertTodo inserts a new todo into the database
func (db *DB) InsertTodo(title, description string, userID int) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	var id int

	query := `
		INSERT INTO todo (created, title, description, user_id)
		VALUES ($1, $2, $3, $4)
		RETURNING id`

	err := db.GetContext(ctx, &id, query, time.Now(), title, description, userID)
	if err != nil {
		return 0, err
	}
	return id, err
}

// GetTodo retrieves a todo from the database
func (db *DB) GetTodo(id int) (*Todo, bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	var todo Todo

	query := `SELECT * FROM todo WHERE id = $1`

	err := db.GetContext(ctx, &todo, query, id)
	if err != nil {
		return nil, false, err
	}

	return &todo, true, err
}

// GetTodos retrieves all todos from the database
// GetTodos retrieves filtered, sorted and paginated todos from the database
func (db *DB) GetTodos(userID int, filters TodoFilters) ([]Todo, int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	// Base query
	query := `SELECT * FROM todo WHERE user_id = $1`
	countQuery := `SELECT COUNT(*) FROM todo WHERE user_id = $1`
	args := []interface{}{userID}
	argPosition := 2

	// Add search filter if provided
	if filters.Search != "" {
		query += ` AND (title ILIKE $` + strconv.Itoa(argPosition) +
			` OR description ILIKE $` + strconv.Itoa(argPosition) + `)`
		countQuery += ` AND (title ILIKE $` + strconv.Itoa(argPosition) +
			` OR description ILIKE $` + strconv.Itoa(argPosition) + `)`
		args = append(args, "%"+filters.Search+"%")
		argPosition++
	}

	// Add sorting
	validSortFields := map[string]bool{
		"title":       true,
		"description": true,
		"created_at":  true,
	}

	if filters.SortBy != "" && validSortFields[strings.ToLower(filters.SortBy)] {
		query += ` ORDER BY ` + strings.ToLower(filters.SortBy)

		// Add sort order if valid
		if strings.ToUpper(filters.SortOrder) == "DESC" {
			query += " DESC"
		} else {
			query += " ASC"
		}
	} else {
		// Default sorting
		query += ` ORDER BY created_at DESC`
	}

	// Add pagination
	if filters.Page < 1 {
		filters.Page = 1
	}
	if filters.PageSize < 1 {
		filters.PageSize = 10 // Default page size
	}

	// Get total count before pagination
	var totalCount int
	err := db.GetContext(ctx, &totalCount, countQuery, args...)
	if err != nil {
		return nil, 0, err
	}

	// Add pagination to query
	offset := (filters.Page - 1) * filters.PageSize
	query += ` LIMIT $` + strconv.Itoa(argPosition) + ` OFFSET $` + strconv.Itoa(argPosition+1)
	args = append(args, filters.PageSize, offset)

	// Execute final query
	var todos []Todo
	err = db.SelectContext(ctx, &todos, query, args...)
	if err != nil {
		return nil, 0, err
	}

	return todos, totalCount, nil
}

// DeleteTodoByID  deletes a todo from the database
func (db *DB) DeleteTodoByID(id int) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	query := `DELETE FROM todo WHERE id = $1`

	_, err := db.ExecContext(ctx, query, id)
	return err
}

// GetTodoByID get todo by ID
func (db *DB) GetTodoByID(id int) (*Todo, bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	var todo Todo

	query := `SELECT * FROM todo WHERE id = $1`

	err := db.GetContext(ctx, &todo, query, id)
	if err != nil {
		return nil, false, err
	}

	return &todo, true, err
}

// UpdateTodo updates a todo in the database
func (db *DB) UpdateTodo(id int, title, description string) error {
	db.mu.Lock()
	defer db.mu.Unlock()
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	query := `UPDATE todo SET title = $1, description = $2 WHERE id = $3`
	_, err := db.ExecContext(ctx, query, title, description, id)
	return err
}

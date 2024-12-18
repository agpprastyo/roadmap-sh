package database

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

type User struct {
	ID             int       `db:"id" json:"id"`
	Created        time.Time `db:"created" json:"created"`
	Name           string    `db:"name" json:"name"`
	Email          string    `db:"email" json:"email"`
	HashedPassword string    `db:"hashed_password" json:"hashedPassword"`
}

func (db *DB) InsertUser(email, hashedPassword string, name string) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	var id int

	query := `
		INSERT INTO users (created, email, hashed_password, name)
		VALUES ($1, $2, $3, $4)
		RETURNING id`

	err := db.GetContext(ctx, &id, query, time.Now(), email, hashedPassword, name)
	if err != nil {
		return 0, err
	}

	return id, err
}

func (db *DB) GetUser(id int) (*User, bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	var user User

	query := `SELECT * FROM users WHERE id = $1`

	err := db.GetContext(ctx, &user, query, id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, false, nil
	}

	return &user, true, err
}

func (db *DB) GetUserByEmail(email string) (*User, bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	var user User

	query := `SELECT * FROM users WHERE email = $1`

	err := db.GetContext(ctx, &user, query, email)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, false, nil
	}

	return &user, true, err
}

func (db *DB) UpdateUserHashedPassword(id int, hashedPassword string) error {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	query := `UPDATE users SET hashed_password = $1 WHERE id = $2`

	_, err := db.ExecContext(ctx, query, hashedPassword, id)
	return err
}

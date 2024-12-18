package database

import (
	"context"

	"github.com/agpprastyo/todo-list-api/internal/validator"
	"time"
)

type RefreshToken struct {
	ID        int       `db:"id"`
	UserID    int       `db:"user_id"`
	Token     string    `db:"token"`
	Expiry    time.Time `db:"expiry"`
	CreatedAt time.Time `db:"created_at"`
	Revoked   bool      `db:"revoked"`
}

// CreateRefreshToken database/refresh_tokens.go
func (db *DB) CreateRefreshToken(userID int, ttl time.Duration) (*RefreshToken, error) {
	token := &RefreshToken{
		UserID: userID,
		Token:  validator.GenerateSecureToken(32), // implement this helper function
		Expiry: time.Now().Add(ttl),
	}

	query := `
        INSERT INTO refresh_tokens (user_id, token, expiry)
        VALUES ($1, $2, $3)
        RETURNING id, created_at`

	err := db.QueryRowContext(context.Background(), query,
		token.UserID, token.Token, token.Expiry).Scan(&token.ID, &token.CreatedAt)

	if err != nil {
		return nil, err
	}
	return token, nil
}

func (db *DB) GetRefreshToken(token string) (*RefreshToken, error) {
	var refreshToken RefreshToken

	query := `
        SELECT id, user_id, token, expiry, created_at, revoked
        FROM refresh_tokens
        WHERE token = $1 AND revoked = false AND expiry > NOW()`

	err := db.GetContext(context.Background(), &refreshToken, query, token)
	if err != nil {
		return nil, err
	}
	return &refreshToken, nil
}

func (db *DB) RevokeRefreshToken(token string) error {
	query := `UPDATE refresh_tokens SET revoked = true WHERE token = $1`
	_, err := db.ExecContext(context.Background(), query, token)
	return err
}

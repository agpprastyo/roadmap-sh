package sessions

import (
	"errors"
	"github.com/google/uuid"
	"net/http"
	"personal-blog/internal/cookies"
)

// ErrSessionNotFound Define a custom error for session handling
var ErrSessionNotFound = errors.New("session not found")

func CreateSession(w http.ResponseWriter, r *http.Request, userID string, secretKey string) (error, error) {
	// Generate a secure, unique session token
	sessionToken := uuid.New().String()

	// Create a secure, HTTP-only cookie
	cookie := http.Cookie{
		Name:     "session-hmm",
		Value:    sessionToken,
		Path:     "/",
		MaxAge:   60 * 60 * 24,          // 24 hours
		HttpOnly: true,                  // Prevents JavaScript access
		Secure:   true,                  // Only sent over HTTPS
		SameSite: http.SameSiteNoneMode, // Crucial for cross-site cookies
	}

	// Store session information (e.g., in a database or in-memory store)
	err := storeSessionData(sessionToken, userID)
	if err != nil {
		return err, nil
	}

	// Write the cookie
	http.SetCookie(w, &cookie)
	return nil, nil
}

func storeSessionData(sessionToken, userID string) error {
	// Implement session storage logic
	// This could be in-memory, Redis, database, etc.
	// Example with in-memory map (not production-ready)
	sessionStore[sessionToken] = userID
	return nil
}

// Global in-memory session store (for demonstration)
var sessionStore = make(map[string]string)

// GetSession retrieves the session value for the authenticated user
func GetSession(r *http.Request) (string, error) {
	// Use the Read function from cookies.go to get the session cookie value
	sessionValue, err := cookies.Read(r, "session-hmm")
	if err != nil {
		return "", ErrSessionNotFound // Return custom error if cookie not found
	}
	return sessionValue, nil
}

// DestroySession invalidates the session
func DestroySession(w http.ResponseWriter, r *http.Request) error {
	// Create a cookie that expires immediately
	cookie := http.Cookie{
		Name:   "session-hmm",
		Value:  "",
		Path:   "/",
		MaxAge: -1, // Expire the cookie
	}

	// Use the Write function from cookies.go to clear the session cookie
	return cookies.Write(w, cookie)
}

package main

import (
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"log/slog"
	"net/http"
	"personal-blog/internal/sessions"

	"personal-blog/internal/response"

	"github.com/tomasen/realip"
)

func (app *application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			err := recover()
			if err != nil {
				app.serverError(w, r, fmt.Errorf("%s", err))
			}
		}()

		next.ServeHTTP(w, r)
	})
}

func (app *application) logAccess(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mw := response.NewMetricsResponseWriter(w)
		next.ServeHTTP(mw, r)

		var (
			ip     = realip.FromRequest(r)
			method = r.Method
			url    = r.URL.String()
			proto  = r.Proto
		)

		userAttrs := slog.Group("user", "ip", ip)
		requestAttrs := slog.Group("request", "method", method, "url", url, "proto", proto)
		responseAttrs := slog.Group("repsonse", "status", mw.StatusCode, "size", mw.BytesCount)

		app.logger.Info("access", userAttrs, requestAttrs, responseAttrs)
	})
}

func (app *application) requireAuthSession(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Attempt to retrieve the session value
		sessionValue, err := sessions.GetSession(r)
		if err != nil {
			app.unauthorized(w, r)
			return
		}

		// Check if the session value is valid (not empty)
		if sessionValue == "" {
			app.unauthorized(w, r)
			return
		}

		// If the session is valid, proceed to the next handler
		next.ServeHTTP(w, r)
	})
}

func (app *application) requireBasicAuthentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		username, plaintextPassword, ok := r.BasicAuth()
		if !ok {
			app.basicAuthenticationRequired(w, r)
			return
		}

		if app.config.basicAuth.username != username {
			app.basicAuthenticationRequired(w, r)
			return
		}

		err := bcrypt.CompareHashAndPassword([]byte(app.config.basicAuth.hashedPassword), []byte(plaintextPassword))
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			app.basicAuthenticationRequired(w, r)
			return
		case err != nil:
			app.serverError(w, r, err)
			return
		}

		// Authentication successful, create a session
		if err, _ := sessions.CreateSession(w, r, username, app.config.cookie.secretKey); err != nil {
			app.serverError(w, r, err) // Handle session creation error
			return
		}

		next.ServeHTTP(w, r)
	})
}

//func (app *application) requireBasicAuthentication(next http.Handler) http.Handler {
//	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//		username, plaintextPassword, ok := r.BasicAuth()
//		if !ok {
//			app.basicAuthenticationRequired(w, r)
//			return
//		}
//
//		if app.config.basicAuth.username != username {
//			app.basicAuthenticationRequired(w, r)
//			return
//		}
//
//		err := bcrypt.CompareHashAndPassword([]byte(app.config.basicAuth.hashedPassword), []byte(plaintextPassword))
//		switch {
//		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
//			app.basicAuthenticationRequired(w, r)
//			return
//		case err != nil:
//			app.serverError(w, r, err)
//			return
//		}
//
//		next.ServeHTTP(w, r)
//	})
//}

package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (app *application) routes() http.Handler {
	mux := chi.NewRouter()

	mux.NotFound(app.notFound)
	mux.MethodNotAllowed(app.methodNotAllowed)

	mux.Use(app.logAccess)
	mux.Use(app.recoverPanic)
	mux.Use(app.authenticate)

	mux.Get("/status", app.status)
	mux.Post("/users", app.createUser)
	mux.Post("/authentication-tokens", app.createAuthenticationToken)

	mux.Group(func(mux chi.Router) {
		mux.Use(app.requireAuthenticatedUser)

		mux.Get("/protected", app.protected)
	})

	mux.Group(func(mux chi.Router) {
		mux.Use(app.requireBasicAuthentication)

		mux.Get("/basic-auth-protected", app.protected)
	})

	return mux
}

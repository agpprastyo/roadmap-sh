package main

import (
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/httprate"
)

func (app *application) routes() http.Handler {
	mux := chi.NewRouter()

	mux.NotFound(app.notFound)
	mux.MethodNotAllowed(app.methodNotAllowed)

	mux.Use(app.logAccess)
	mux.Use(app.recoverPanic)
	mux.Use(app.authenticate)
	mux.Use(middleware.Throttle(15))

	mux.Use(httprate.LimitByIP(100, 1*time.Minute))

	mux.Get("/status", app.status)
	mux.Post("/users", app.createUser)
	mux.Post("/authentication-tokens", app.createAuthenticationToken)
	mux.Post("/refresh-tokens", app.refreshToken)

	mux.Group(func(mux chi.Router) {
		mux.Use(app.requireAuthenticatedUser)

		mux.Get("/protected", app.protected)

		mux.Route("/todo", func(mux chi.Router) {
			mux.Get("/", app.getAllTodos)
			mux.Post("/", app.createTodo)
			mux.Delete("/{id:[1-9][0-9]*}", app.deleteTodo)
			mux.Patch("/{id:[1-9][0-9]*}", app.updateTodo)
		})
	})

	mux.Group(func(mux chi.Router) {
		mux.Use(app.requireBasicAuthentication)

		mux.Get("/basic-auth-protected", app.protected)
	})

	return mux
}

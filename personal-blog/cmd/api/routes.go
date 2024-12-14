package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	_ "github.com/swaggo/files"
	httpSwagger "github.com/swaggo/http-swagger"
	_ "github.com/swaggo/http-swagger/example/go-chi/docs"
	"net/http"
	"path/filepath"
)

func (app *application) routes() http.Handler {
	mux := chi.NewRouter()
	// Add CORS middleware
	// CORS middleware
	mux.Use(cors.Handler(cors.Options{

		AllowedOrigins:   []string{"http://localhost:5173"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	mux.NotFound(app.notFound)
	mux.MethodNotAllowed(app.methodNotAllowed)

	mux.Use(app.logAccess)
	mux.Use(app.recoverPanic)

	// Add OPTIONS handler
	mux.Options("/*", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, PATCH")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Authorization, Content-Type, X-CSRF-Token")
		w.WriteHeader(http.StatusOK)
	})

	// Add Swagger UI route
	mux.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"), // The URL pointing to API definition
		httpSwagger.UIConfig(map[string]string{
			"displayRequestDuration":   "true",
			"defaultModelsExpandDepth": "-1",
			"deepLinking":              "true",
			"displayOperationId":       "true",
		}),
	))

	// Serve the swagger.yaml file
	mux.Get("/swagger/doc.json", func(w http.ResponseWriter, r *http.Request) {
		yamlFile := filepath.Join("cmd", "docs", "personal-blog.swagger.yaml")
		http.ServeFile(w, r, yamlFile)
	})

	mux.Get("/status", app.status)

	mux.Group(func(mux chi.Router) {
		mux.Use(app.requireBasicAuthentication)

		mux.Get("/basic-auth-protected", app.protected)
	})

	// Group for the API version 1 routes
	mux.Route("/api/v1", func(mux chi.Router) {

		mux.Get("/articles", app.readArticles)
		mux.Get("/article/{id:[1-9][0-9]*}", app.readArticleByID)

		//docs

		mux.Group(func(mux chi.Router) {
			mux.Use(app.requireAuthSession)

			mux.Get("/admin", app.readArticlesAll)
			mux.Get("/admin/article/{id:[1-9][0-9]*}", app.readArticleByIDAdmin)

			mux.Post("/create", app.createArticle)
			mux.Patch("/edit/{id:[1-9][0-9]*}", app.updateArticle)
			mux.Delete("/delete/{id:[1-9][0-9]*}", app.deleteArticle)
			// restore article
			mux.Post("/restore/{id:[1-9][0-9]*}", app.restoreArticle)
		})

		mux.Post("/sign-in", app.signIn)   // Admin sign-in
		mux.Post("/sign-out", app.signOut) // Admin sign-out
	})

	return mux
}

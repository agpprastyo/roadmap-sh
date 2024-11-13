package server

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (s *Server) RegisterRoutes() http.Handler {
	r := gin.Default()

	r.GET("/", s.HomeHandler)
	r.GET("/article/:id", s.ArticleHandler)

	r.GET("/admin", s.AdminHandler)

	r.GET("/login", s.LoginHandler)

	//r.GET("/health", s.healthHandler)

	r.Static("/assets", "./cmd/web/assets")

	return r
}

package server

import (
	"github.com/a-h/templ"
	"github.com/gin-gonic/gin"
	"personal-blog/cmd/web"
)

func (s *Server) HomeHandler(c *gin.Context) {
	c.Writer.Header().Set("Content-Type", "text/html")
	templ.Handler(web.Home()).ServeHTTP(c.Writer, c.Request)
}

func (s *Server) ArticleHandler(c *gin.Context) {

}

func (s *Server) AdminHandler(c *gin.Context) {
	c.Writer.Header().Set("Content-Type", "text/html")
	templ.Handler(web.Admin()).ServeHTTP(c.Writer, c.Request)
}

func (s *Server) LoginHandler(c *gin.Context) {
	c.Writer.Header().Set("Content-Type", "text/html")
	templ.Handler(web.Login()).ServeHTTP(c.Writer, c.Request)
}

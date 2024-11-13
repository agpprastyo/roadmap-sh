package server

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
	"unit-converter/cmd/web"
)

func (s *Server) RegisterRoutes() http.Handler {
	r := httprouter.New()
	r.HandlerFunc(http.MethodGet, "/", s.HomeHandler)
	r.HandlerFunc(http.MethodPost, "/weight", s.WeightHandler)
	r.HandlerFunc(http.MethodPost, "/temperature", s.TempHandler)
	r.HandlerFunc(http.MethodPost, "/length", s.LengthHandler)

	fileServer := http.FileServer(http.FS(web.Files))
	r.Handler(http.MethodGet, "/assets/*filepath", fileServer)
	//r.Handler(http.MethodGet, "/web", templ.Handler(web.HelloForm()))
	//r.HandlerFunc(http.MethodPost, "/hello", web.HelloWebHandler)

	return r
}

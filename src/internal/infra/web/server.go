package web

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Handler struct {
	Path        string
	Method      string
	HandlerFunc http.HandlerFunc
}

type Middleware struct {
	Name    string
	Handler func(next http.Handler) http.Handler
}

type Server struct {
	Router   chi.Router
	Port     int
	Handlers []Handler
	// Middlewares []Middleware
}

func NewServer(serverPort int, handlers []Handler) *Server {
	return &Server{
		Router:   chi.NewRouter(),
		Port:     serverPort,
		Handlers: handlers,
		// Middlewares: middlewares,
	}
}

func (s *Server) Run() {
	// for _, m := range s.Middlewares {
	// 	s.Router.Use(m.Handler)
	// }
	for _, h := range s.Handlers {
		s.Router.MethodFunc(h.Method, h.Path, h.HandlerFunc)
	}
	fmt.Printf("Starting server on port [%d]\n", s.Port)
	http.ListenAndServe(fmt.Sprintf(":%d", s.Port), s.Router)
}

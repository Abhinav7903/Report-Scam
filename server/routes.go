package server

import (
	"abuse/db/postgres"
	"net/http"

	"github.com/gorilla/mux"
)

type ResponseMsg struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type Server struct {
	router   *mux.Router
	postgres *postgres.Postgres
}

func (s *Server) RegisterRoutes() {
	s.router.HandleFunc("/ping", s.HandlePong())

}

func (s *Server) HandlePong() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.respond(
			w,
			"pong",
			http.StatusOK,
			nil,
		)
	}
}

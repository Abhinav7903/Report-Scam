package server

import (
	"abuse/db/postgres"
	"abuse/pkg/mail"
	sessmanager "abuse/pkg/sessmanger"
	users "abuse/pkg/user"
	"net/http"

	"github.com/gorilla/mux"
)

type ResponseMsg struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type Server struct {
	router      *mux.Router
	postgres    *postgres.Postgres
	user        users.Repository
	mail        mail.Repository
	sessmanager sessmanager.Repository
}

func (s *Server) RegisterRoutes() {
	s.router.HandleFunc("/ping", s.HandlePong())
	//  User routes
	s.router.HandleFunc("/users", s.handleCreateUser()).Methods(http.MethodPost, http.MethodOptions)
	s.router.HandleFunc("/users", s.handleGetUser()).Methods(http.MethodGet, http.MethodOptions)
	s.router.HandleFunc("/verify", s.handleVerify()).Methods(http.MethodGet, http.MethodOptions)
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

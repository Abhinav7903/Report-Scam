package server

import (
	"abuse/factory"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/sirupsen/logrus"
)

func (s *Server) handleCreateUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var user factory.User
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			logrus.Error("Error decoding user:", err)
			s.respond(w, nil, http.StatusBadRequest, err)
			return
		}

		_, err := s.user.CreateUser(user)
		if err != nil {
			logrus.Error("Error creating user:", err)
			s.respond(w, nil, http.StatusInternalServerError, err)
			return
		}

		// Generate an email hash
		hash, err := s.sessmanager.StoreEmailHash(user.Email)
		if err != nil {
			s.respond(w, nil, http.StatusInternalServerError, fmt.Errorf("failed to process email verification"))
			return
		}

		// Send verification email to the user
		err = s.mail.SendMail(
			user.Email,
			"Verify your email",
			"Click the link to verify your email: http://localhost:8194/verify?ehash="+hash,
		)
		if err != nil {
			s.respond(w, nil, http.StatusInternalServerError, fmt.Errorf("verification email failed"))
			return
		}

		// Notify admin about the new user
		err = s.mail.SendMail(
			"abhinavashish4@gmail.com",
			"New user signed up",
			"New user signed up with email: "+user.Email+" and name: "+user.Name,
		)
		if err != nil {
			s.respond(w, nil, http.StatusInternalServerError, fmt.Errorf("admin notification failed"))
			return
		}

		logrus.Info("User created successfully:", user.Email)
		s.respond(w, ResponseMsg{Message: "user created successfully"}, http.StatusCreated, nil)
	}
}

func (s *Server) handleVerify() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Extract the email hash from the query parameters
		hash := r.URL.Query().Get("ehash")
		if hash == "" {
			s.respond(w, nil, http.StatusBadRequest, fmt.Errorf("invalid verification link"))
			return
		}

		// Retrieve email associated with the hash
		email, err := s.sessmanager.GetEmailFromHash(hash)
		if err != nil {
			logrus.Error("Error retrieving email from hash:", err)
			s.respond(w, nil, http.StatusInternalServerError, fmt.Errorf("email verification failed"))
			return
		}

		// Mark the email as verified in the database
		if err := s.user.VerifyEmail(email); err != nil {
			logrus.Error("failed to update email verification status:", err)
			s.respond(w, nil, http.StatusInternalServerError, fmt.Errorf("email verification update failed"))
			return
		}

		// Respond with success message
		s.respond(w, ResponseMsg{Message: "Email verified successfully"}, http.StatusOK, nil)
	}
}

func (s *Server) handleGetUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		email := r.URL.Query().Get("email")
		if email == "" {
			logrus.Error("Email query parameter is required")
			s.respond(w, ResponseMsg{Message: "Email query parameter is required"}, http.StatusBadRequest, nil)
			return
		}

		user, err := s.user.GetUser(email)
		if err != nil {
			logrus.Error("Error getting user:", err)
			s.respond(w, ResponseMsg{Message: "Error getting user"}, http.StatusInternalServerError, nil)
			return
		}

		logrus.Info("User retrieved successfully:", user.Email)
		s.respond(w, user, http.StatusOK, nil)
	}
}

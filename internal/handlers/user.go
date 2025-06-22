package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-kit/log"
	"github.com/gorilla/mux"

	"github.com/shahabeshaghi67/Golang-Microservice-Template/internal/model"
)

// NewUserHandler returns a new instance of UserHandler.
func NewUserHandler(soService UserService, logger log.Logger) *UserHandler {
	return &UserHandler{
		UserService: soService,
		logger:      logger,
	}
}

// UserHandler underlying data.
type UserHandler struct {
	UserService UserService
	logger      log.Logger
}

// CreateUser returns a handler that creates a new user from the request body.
func (h *UserHandler) CreateUser() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var user model.User
		var err error

		// Decode request body
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			h.logger.Log("err", fmt.Sprintf("failed to decode user: %+v", err))
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}

		// Validate user
		if err := user.Validate(); err != nil {
			h.logger.Log("err", fmt.Sprintf("user validation failed: %+v", err))
			http.Error(w, "user validation failed", http.StatusUnprocessableEntity)
			return
		}

		// Process user creation
		if user, err = h.UserService.Create(r.Context(), user); err != nil {
			h.logger.Log("err", fmt.Sprintf("failed to create user: %+v", err))
			http.Error(w, "failed to create user", http.StatusInternalServerError)
			return
		}

		h.logger.Log("info", fmt.Sprintf("user created successfully: %v", user.ID))
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(user)
	})
}

// GetUserByID returns a handler that retrieves a user by ID from the URL.
func (h *UserHandler) GetUserByID() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id, found := mux.Vars(r)["id"]
		if id == "" || !found {
			h.logger.Log("err", "missing user ID in URL")
			http.Error(w, "missing user ID", http.StatusBadRequest)
			return
		}

		user, err := h.UserService.GetByID(r.Context(), id)
		if err != nil {
			h.logger.Log("err", fmt.Sprintf("failed to get user by ID: %+v", err))
			http.Error(w, "failed to get user by ID", http.StatusInternalServerError)
			return
		}

		h.logger.Log("info", fmt.Sprintf("user retrieved successfully: %v", user.ID))
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(user)
	})
}

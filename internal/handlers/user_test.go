package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-kit/log"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/shahabeshaghi67/Golang-Microservice-Template/internal/handlers"
	"github.com/shahabeshaghi67/Golang-Microservice-Template/internal/mocks"
	"github.com/shahabeshaghi67/Golang-Microservice-Template/internal/model"
)

const userHandlerPath = "/users"

func setupUsersRouter() (http.Handler, *mocks.UserService) {
	router := mux.NewRouter()
	userService := &mocks.UserService{}
	userHandler := handlers.NewUserHandler(
		userService,
		log.NewNopLogger(),
	)

	router.Methods("POST").
		Path(userHandlerPath).
		Handler(userHandler.CreateUser())

	router.Methods("Get").
		Path(userHandlerPath + "/{id}").
		Handler(userHandler.GetUserByID())

	return router, userService
}

func TestUserHandler_CreateUser(t *testing.T) {
	t.Run("invalid request body", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/users", bytes.NewBufferString("invalid-json"))
		w := httptest.NewRecorder()
		router, _ := setupUsersRouter()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("validation error", func(t *testing.T) {
		invalidUser := model.User{}
		body, _ := json.Marshal(invalidUser)
		req := httptest.NewRequest("POST", "/users", bytes.NewBuffer(body))
		w := httptest.NewRecorder()
		router, _ := setupUsersRouter()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
	})

	t.Run("service error", func(t *testing.T) {
		user := testUser[0]
		body, _ := json.Marshal(user)
		req := httptest.NewRequest("POST", "/users", bytes.NewBuffer(body))
		w := httptest.NewRecorder()
		router, userService := setupUsersRouter()

		userService.On("Create", mock.Anything, user).Return(model.User{}, assert.AnError)

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("success", func(t *testing.T) {
		user := testUser[0]
		body, _ := json.Marshal(user)
		req := httptest.NewRequest("POST", "/users", bytes.NewBuffer(body))
		w := httptest.NewRecorder()
		router, userService := setupUsersRouter()

		userService.On("Create", mock.Anything, user).Return(user, nil)

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
	})
}

func TestUserHandler_GetUserByID(t *testing.T) {
	user := testUser[0]
	userID := user.ID.String()

	t.Run("service error", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/users/"+userID, nil)
		w := httptest.NewRecorder()
		router, userService := setupUsersRouter()

		userService.On("GetByID", mock.Anything, userID).Return(model.User{}, assert.AnError)

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("success", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/users/"+userID, nil)
		w := httptest.NewRecorder()
		router, userService := setupUsersRouter()

		userService.On("GetByID", mock.Anything, userID).Return(user, nil)

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var got model.User
		err := json.NewDecoder(w.Body).Decode(&got)
		assert.NoError(t, err)
		assert.Equal(t, user, got)
	})
}

package api

import (
	"net/http"

	"github.com/go-playground/validator/v10"

	"github.com/shahabeshaghi67/Golang-Microservice-Template/internal/utils"
)

var (
	routeValidator *validator.Validate
)

func init() {
	v := validator.New()
	v.RegisterTagNameFunc(utils.JsonTagName)
	routeValidator = v
}

// API is a struct that holds all the information needed to create a new API endpoint.
type API struct {
	Path           string            `validate:"required,uri"`
	Handler        http.Handler      `validate:"required"`
	Method         string            `validate:"required,oneof=GET POST PUT DELETE PATCH OPTIONS HEAD TRACE CONNECT"`
	RequestType    any               `validate:"-"`
	ResponseType   []APIResponseType `validate:"required,dive"`
	Tags           []string
	Summary        string
	Description    string
	SecurityName   string
	SecurityScopes []string
	Deprecated     bool
}

// APIResponseType is a struct that holds the response type for an API endpoint.
type APIResponseType struct {
	StatusCode int `validate:"required,oneof=100 101 200 201 202 204 400 401 403 404 405 409 422 500"`
	Body       any `validate:"-"`
}

// ServeHTTP is wrapper for calling handler method.
func (a *API) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.Handler.ServeHTTP(w, r)
}

// Validate is used to check if struct fields meet data requirements.
func (a *API) Validate() error {
	return routeValidator.Struct(a)
}

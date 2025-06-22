package app

import (
	"net/http"

	"github.com/shahabeshaghi67/Golang-Microservice-Template/internal/api"
	"github.com/shahabeshaghi67/Golang-Microservice-Template/internal/model"
)

func (a *App) generateRoutes() []*api.API {
	// write routes here
	a.routes = []*api.API{
		{
			Path:         "/users",
			Method:       http.MethodPost,
			Handler:      a.handlers["users.create"],
			Tags:         []string{"user"},
			Summary:      "Validates and creates a new user",
			Description:  "Validates and creates a new user",
			ResponseType: model.User{},
			RequestType:  model.User{},
		},
		{
			Path:         "/users/{id}",
			Method:       http.MethodGet,
			Handler:      a.handlers["users.getByID"],
			Tags:         []string{"user"},
			Summary:      "Get user by ID",
			Description:  "Get user by ID",
			ResponseType: model.User{},
			
		},
	}
	return a.routes
}

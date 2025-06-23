package app

import (
	"net/http"

	"github.com/shahabeshaghi67/Golang-Microservice-Template/internal/api"
	"github.com/shahabeshaghi67/Golang-Microservice-Template/internal/model"
)

func (a *App) generateRoutes() ([]*api.API, error) {
	// write routes here
	a.Routes = []*api.API{
		{
			Path:        "/users",
			Method:      http.MethodPost,
			Handler:     a.handlers["users.create"],
			Tags:        []string{"user"},
			Summary:     "Validates and creates a new user",
			Description: "Validates and creates a new user",
			ResponseType: []api.APIResponseType{
				{
					StatusCode: http.StatusCreated,
					Body:       model.User{},
				},
				{
					StatusCode: http.StatusBadRequest,
				},
				{
					StatusCode: http.StatusUnprocessableEntity,
				},
			},
			RequestType: model.User{},
		},
		{
			Path:        "/users/{id}",
			Method:      http.MethodGet,
			Handler:     a.handlers["users.getByID"],
			Tags:        []string{"user"},
			Summary:     "Get user by ID",
			Description: "Get user by ID",
			RequestType: struct {
				ID string `path:"id"`
			}{},
			ResponseType: []api.APIResponseType{
				{
					StatusCode: http.StatusOK,
					Body:       model.User{},
				},
				{
					StatusCode: http.StatusBadRequest,
				},
				{
					StatusCode: http.StatusNotFound,
				},
			},
		},
	}

	for _, route := range a.Routes {
		if err := route.Validate(); err != nil {
			return nil, err
		}
	}

	return a.Routes, nil
}

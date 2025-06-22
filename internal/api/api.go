package api

import (
	"net/http"

	"github.com/swaggest/openapi-go"
)

// API is a struct that holds all the information needed to create a new API endpoint.
type API struct {
	Path           string
	Handler        http.Handler
	Method         string
	RequestType    any
	ResponseType   any
	Tags           []string
	Summary        string
	Description    string
	SecurityName   string
	SecurityScopes []string
}

// SetupOpenAPIOperation declares OpenAPI schema for the handler.
func (a *API) SetupOpenAPIOperation(oc openapi.OperationContext) error {
	for _, tag := range a.Tags {
		oc.SetTags(tag)
	}
	if a.Summary != "" {
		oc.SetSummary(a.Summary)
	}
	if a.Description != "" {
		oc.SetDescription(a.Description)
	}
	if a.RequestType != nil {
		oc.AddReqStructure(a.RequestType)
	}
	if a.ResponseType != nil {
		oc.AddRespStructure(a.ResponseType)
	}

	oc.AddSecurity(a.SecurityName, a.SecurityScopes...)

	return nil
}

// ServeHTTP is wrapper for calling handler method.
func (a *API) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.Handler.ServeHTTP(w, r)
}

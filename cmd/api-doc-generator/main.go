package main

import (
	"fmt"
	"os"

	"github.com/go-kit/log/level"
	"github.com/swaggest/openapi-go"
	"github.com/swaggest/openapi-go/openapi3"

	"github.com/shahabeshaghi67/Golang-Microservice-Template/internal/app"
	"github.com/shahabeshaghi67/Golang-Microservice-Template/internal/config"
	"github.com/shahabeshaghi67/Golang-Microservice-Template/pkg/logger"
)

var apiYmlFile = "doc/api.yaml"

func main() {
	logger := logger.NewLogger()
	cfg := config.Load()

	app := app.NewApp(cfg, logger, true)
	_, err := app.Wire()
	if err != nil {
		_ = level.Error(logger).Log("msg", "app wiring failed", "error:", err)
		panic(err)
	}

	// Setup OpenAPI schema.

	reflector := openapi3.Reflector{}
	reflector.Spec = &openapi3.Spec{Openapi: "3.0.3"}

	productionDescription := "production server"
	stagingDescription := "staging server"
	refl := openapi3.NewReflector()
	refl.Spec.Info.WithTitle("golang API service").
		WithVersion("0.0.1").
		WithDescription("golang API service handles users Data") // ATTENTION: add description
	refl.Spec.WithServers(
		openapi3.Server{
			URL:         "https://production.example.com",
			Description: &productionDescription,
		},
		openapi3.Server{
			URL:         "https://staging.example.com",
			Description: &stagingDescription,
		},
	)
	teamName := ""
	teamURL := ""
	teamEmail := ""
	oauth2Description := `Note: In addition to the required scopes services accessing this API must be explicitly allowed. Please make contact directly if you need aservice to be added to the allowed list.`

	refl.Spec.WithComponents(openapi3.Components{
		SecuritySchemes: &openapi3.ComponentsSecuritySchemes{
			MapOfSecuritySchemeOrRefValues: map[string]openapi3.SecuritySchemeOrRef{
				"OAuth2": {
					SecurityScheme: &openapi3.SecurityScheme{
						OAuth2SecurityScheme: &openapi3.OAuth2SecurityScheme{
							Description: &oauth2Description,
							Flows: openapi3.OAuthFlows{
								ClientCredentials: &openapi3.ClientCredentialsFlow{
									TokenURL: "",
									Scopes:   map[string]string{
										//ATTENTION: add oauth2 scopes here
									},
								},
							},
							MapOfAnything: map[string]interface{}{},
						},
					},
				},
			},
		},
	})
	refl.Spec.Info.WithContact(openapi3.Contact{
		Name:  &teamName,
		URL:   &teamURL,
		Email: &teamEmail,
	})
	refl.Spec.Info.MapOfAnything = map[string]interface{}{}
	refl.Spec.Info.MapOfAnything["x-audience"] = ""
	refl.Spec.Info.MapOfAnything["x-api-id"] = "" //ATTENTION: add api id

	for _, route := range app.Routes {
		op, err := reflector.NewOperationContext(route.Method, route.Path)
		if err != nil {
			panic(fmt.Sprintf("Unable to create operation context for %s %s: %v", route.Method, route.Path, err))
		}
		for _, tag := range route.Tags {
			op.SetTags(tag)
		}
		if route.Summary != "" {
			op.SetSummary(route.Summary)
		}
		if route.Description != "" {
			op.SetDescription(route.Description)
		}
		if route.RequestType != nil {
			op.AddReqStructure(route.RequestType)
		}
		for _, resp := range route.ResponseType {
			op.AddRespStructure(resp.Body, func(cu *openapi.ContentUnit) { cu.HTTPStatus = resp.StatusCode })
		}
		op.SetIsDeprecated(route.Deprecated)
		op.AddSecurity(route.SecurityName, route.SecurityScopes...)

		refl.AddOperation(op)
	}

	// Get the resulting schema.
	apiYml, _ := refl.Spec.MarshalYAML()
	err = os.WriteFile(apiYmlFile, apiYml, 0600)
	if err != nil {
		panic("Unable to write data into the file")
	}
}

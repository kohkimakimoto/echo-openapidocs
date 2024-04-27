package main

import (
	_ "embed"
	"github.com/kohkimakimoto/echo-openapidocs"
	"github.com/labstack/echo/v4"
	"log"
)

// OpenAPISpecGithub is the OpenAPI specifications for GitHub v3 REST API.
// I got the OpenAPI specifications from https://github.com/github/rest-api-description
//
//go:embed ghes-3.0.yaml
var OpenAPISpecGithub string

func main() {
	e := echo.New()

	// ElementsDocuments
	openapidocs.ElementsDocuments(e, "/docs/elements/github", openapidocs.ElementsConfig{
		Spec:  OpenAPISpecGithub,
		Title: "GitHub v3 REST API",
	})
	openapidocs.ElementsDocuments(e, "/docs/elements/openai", openapidocs.ElementsConfig{
		SpecUrl: "https://raw.githubusercontent.com/openai/openai-openapi/master/openapi.yaml",
		Title:   "OpenAI API",
	})

	// ScalarDocuments
	openapidocs.ScalarDocuments(e, "/docs/scalar/github", openapidocs.ScalarConfig{
		Spec:  OpenAPISpecGithub,
		Title: "GitHub v3 REST API",
	})
	openapidocs.ScalarDocuments(e, "/docs/scalar/openai", openapidocs.ScalarConfig{
		SpecUrl: "https://raw.githubusercontent.com/openai/openai-openapi/master/openapi.yaml",
		Title:   "OpenAI API",
	})

	// SwaggerUIDocuments
	openapidocs.SwaggerUIDocuments(e, "/docs/swagger-ui/github", openapidocs.SwaggerUIConfig{
		Spec:  OpenAPISpecGithub,
		Title: "GitHub v3 REST API",
	})
	openapidocs.SwaggerUIDocuments(e, "/docs/swagger-ui/openai", openapidocs.SwaggerUIConfig{
		SpecUrl: "https://raw.githubusercontent.com/openai/openai-openapi/master/openapi.yaml",
		Title:   "OpenAI API",
	})

	// RedocDocuments
	openapidocs.RedocDocuments(e, "/docs/redoc/github", openapidocs.RedocConfig{
		Spec:  OpenAPISpecGithub,
		Title: "GitHub v3 REST API",
	})
	openapidocs.RedocDocuments(e, "/docs/redoc/openai", openapidocs.RedocConfig{
		SpecUrl: "https://raw.githubusercontent.com/openai/openai-openapi/master/openapi.yaml",
		Title:   "OpenAI API",
	})

	// Start the server
	if err := e.Start(":8080"); err != nil {
		log.Fatal(err)
	}
}

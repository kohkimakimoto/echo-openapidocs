package main

import (
	_ "embed"
	"log"

	openapidocs "github.com/kohkimakimoto/echo-openapidocs/v4"
	"github.com/kohkimakimoto/echo-openapidocs/v4/embedded/elements"
	"github.com/kohkimakimoto/echo-openapidocs/v4/embedded/redoc"
	"github.com/labstack/echo/v4"
)

// OpenAPISpecGithub is the OpenAPI specifications for GitHub v3 REST API.
// I got the OpenAPI specifications from https://github.com/github/rest-api-description
//
//go:embed api.github.com.yaml
var OpenAPISpecGithub string

func main() {
	e := echo.New()

	// ElementsDocuments (CDN assets)
	openapidocs.ElementsDocuments(e, "/docs/elements/github", openapidocs.ElementsConfig{
		Spec:  OpenAPISpecGithub,
		Title: "GitHub v3 REST API",
	})
	openapidocs.ElementsDocuments(e, "/docs/elements/openai", openapidocs.ElementsConfig{
		SpecUrl: "https://raw.githubusercontent.com/openai/openai-openapi/master/openapi.yaml",
		Title:   "OpenAI API",
	})

	// ElementsDocuments with embedded assets (no CDN)
	elements.Documents(e, "/docs/elements-embedded/github", openapidocs.ElementsConfig{
		Spec:  OpenAPISpecGithub,
		Title: "GitHub v3 REST API (embedded assets)",
	})
	elements.Documents(e, "/docs/elements-embedded/openai", openapidocs.ElementsConfig{
		SpecUrl: "https://raw.githubusercontent.com/openai/openai-openapi/master/openapi.yaml",
		Title:   "OpenAI API (embedded assets)",
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

	// RedocDocuments (CDN assets)
	openapidocs.RedocDocuments(e, "/docs/redoc/github", openapidocs.RedocConfig{
		Spec:  OpenAPISpecGithub,
		Title: "GitHub v3 REST API",
	})
	openapidocs.RedocDocuments(e, "/docs/redoc/openai", openapidocs.RedocConfig{
		SpecUrl: "https://raw.githubusercontent.com/openai/openai-openapi/master/openapi.yaml",
		Title:   "OpenAI API",
	})

	// RedocDocuments with embedded assets (no CDN)
	redoc.Documents(e, "/docs/redoc-embedded/github", openapidocs.RedocConfig{
		Spec:  OpenAPISpecGithub,
		Title: "GitHub v3 REST API (embedded assets)",
	})
	redoc.Documents(e, "/docs/redoc-embedded/openai", openapidocs.RedocConfig{
		SpecUrl: "https://raw.githubusercontent.com/openai/openai-openapi/master/openapi.yaml",
		Title:   "OpenAI API (embedded assets)",
	})

	// Start the server
	if err := e.Start(":8080"); err != nil {
		log.Fatal(err)
	}
}

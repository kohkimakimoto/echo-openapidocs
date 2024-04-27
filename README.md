# Echo OpenAPI Docs

[![MIT License](https://img.shields.io/badge/license-MIT-blue.svg)](https://github.com/kohkimakimoto/echo-openapidocs/blob/master/LICENSE)
[![Go Reference](https://pkg.go.dev/badge/github.com/kohkimakimoto/echo-openapidocs.svg)](https://pkg.go.dev/github.com/kohkimakimoto/echo-openapidocs)

Echo OpenAPI Docs is a Go library that allows you to use OpenAPI documentation generators in your [Echo](https://github.com/labstack/echo) application.
It enables you to easily host interactive API documentation based on OpenAPI specifications.

It supports several generators. Please see the [Supported Documentation Generators](#supported-documentation-generators) section.
If you are not familiar with OpenAPI, please visit the [OpenAPI Specification](https://swagger.io/specification/) page.

## Installation

```sh
go get github.com/kohkimakimoto/echo-openapidocs
```

## Minimum Example

```go
package main

import (
	"github.com/kohkimakimoto/echo-openapidocs"
	"github.com/labstack/echo/v4"
	"log"
)

func main() {
	e := echo.New()

	// Register the GitHub v3 REST API documentation at /docs
	openapidocs.ElementsDocuments(e, "/docs", openapidocs.ElementsConfig{
		SpecUrl:  "https://raw.githubusercontent.com/github/rest-api-description/main/descriptions/ghes-3.0/ghes-3.0.yaml",
		Title: "GitHub v3 REST API",
	})

	// Start the server
	if err := e.Start(":8080"); err != nil {
		log.Fatal(err)
	}
}
```

You can see the documentation at `http://localhost:8080/docs`.

You can also see other examples in [examples/main.go](examples/main.go) file.

## Supported Documentation Generators

### Spotlight Elements

[Spotlight Elements](https://github.com/stoplightio/elements): Build beautiful, interactive API Docs with embeddable React or Web Components, powered by OpenAPI and Markdown.

```go
// Register the Spotlight Elements documentation with OpenAPI Spec url.
openapidocs.ElementsDocuments(e, "/docs", openapidocs.ElementsConfig{
	// The following is the example for generating the OpenAI API documentation.
	// You can replace the SpecUrl with your OpenAPI Spec url.
	SpecUrl: "https://raw.githubusercontent.com/openai/openai-openapi/master/openapi.yaml",
})
```

The generated documentation looks like this:

![](https://raw.githubusercontent.com/kohkimakimoto/echo-openapidocs/main/images/example-elements.png)

You can also specify the OpenAPI Spec as a string.

```go
// Using the `go:embed` directive is a great way to embed the OpenAPI Spec file as a string.
//go:embed openai-openapi.yaml
var OpenAIAPISpec string

// Register the Spotlight Elements documentation with OpenAPI Spec string.
openapidocs.ElementsDocuments(e, "/docs", openapidocs.ElementsConfig{
	Spec: OpenAIAPISpec,
})
```

The `ElementsDocuments` function takes a configuration from the `ElementsConfig` struct.
For more details, please refer to the documentation [here](https://pkg.go.dev/github.com/kohkimakimoto/echo-openapidocs#ElementsConfig).

### Scalar API Reference

[Scalar API Reference](https://github.com/scalar/scalar): Beautiful API references from OpenAPI/Swagger files ✨

```go
// Register the Scalar documentation with OpenAPI Spec url.
openapidocs.ScalarDocuments(e, "/docs", openapidocs.ScalarConfig{
	// The following is the example for generating the OpenAI API documentation.
	// You can replace the SpecUrl with your OpenAPI Spec url.
	SpecUrl: "https://raw.githubusercontent.com/openai/openai-openapi/master/openapi.yaml",
})
```

The generated documentation looks like this:

![](https://raw.githubusercontent.com/kohkimakimoto/echo-openapidocs/main/images/example-scalar.png)

You can also specify the OpenAPI Spec as a string.

```go
// Using the `go:embed` directive is a great way to embed the OpenAPI Spec file as a string.
//go:embed openai-openapi.yaml
var OpenAIAPISpec string

// Register the Scalar documentation with OpenAPI Spec string.
openapidocs.ScalarDocuments(e, "/docs", openapidocs.ScalarConfig{
	Spec: OpenAIAPISpec,
})
```

The `ScalarDocuments` function takes a configuration from the `ScalarConfig` struct.
For more details, please refer to the documentation [here](https://pkg.go.dev/github.com/kohkimakimoto/echo-openapidocs#ScalarConfig).

### Swagger UI

[Swagger UI](https://github.com/swagger-api/swagger-ui): Swagger UI is a collection of HTML, JavaScript, and CSS assets that dynamically generate beautiful documentation from a Swagger-compliant API.

```go
// Register the SwaggerUI documentation with OpenAPI Spec url.
openapidocs.SwaggerUIDocuments(e, "/docs", openapidocs.SwaggerUIConfig{
	// The following is the example for generating the OpenAI API documentation.
	// You can replace the SpecUrl with your OpenAPI Spec url.
	SpecUrl: "https://raw.githubusercontent.com/openai/openai-openapi/master/openapi.yaml",
})
```

The generated documentation looks like this:

![](https://raw.githubusercontent.com/kohkimakimoto/echo-openapidocs/main/images/example-swaggerui.png)

You can also specify the OpenAPI Spec as a string.

```go
// Using the `go:embed` directive is a great way to embed the OpenAPI Spec file as a string.
//go:embed openai-openapi.yaml
var OpenAIAPISpec string

// Register the SwaggerUI documentation with OpenAPI Spec string.
openapidocs.SwaggerUIDocuments(e, "/docs", openapidocs.SwaggerUIConfig{
	Spec: OpenAIAPISpec,
})
```

The `SwaggerUIDocuments` function takes a configuration from the `SwaggerUIConfig` struct.
For more details, please refer to the documentation [here](https://pkg.go.dev/github.com/kohkimakimoto/echo-openapidocs#SwaggerUIConfig).

### ReDoc

[ReDoc](https://github.com/Redocly/redoc): 📘 OpenAPI/Swagger-generated API Reference Documentation.

```go
// Register the Redoc documentation with OpenAPI Spec url.
openapidocs.RedocDocuments(e, "/docs", openapidocs.RedocConfig{
	// The following is the example for generating the OpenAI API documentation.
	// You can replace the SpecUrl with your OpenAPI Spec url.
	SpecUrl: "https://raw.githubusercontent.com/openai/openai-openapi/master/openapi.yaml",
})
```

The generated documentation looks like this:

![](https://raw.githubusercontent.com/kohkimakimoto/echo-openapidocs/main/images/example-redoc.png)

You can also specify the OpenAPI Spec as a string.

```go
// Using the `go:embed` directive is a great way to embed the OpenAPI Spec file as a string.
//go:embed openai-openapi.yaml
var OpenAIAPISpec string

// Register the Redoc documentation with OpenAPI Spec string.
openapidocs.RedocDocuments(e, "/docs", openapidocs.RedocConfig{
	Spec: OpenAIAPISpec,
})
```

The `RedocDocuments` function takes a configuration from the `RedocConfig` struct.
For more details, please refer to the documentation [here](https://pkg.go.dev/github.com/kohkimakimoto/echo-openapidocs#RedocConfig).

## Author

Kohki Makimoto <kohki.makimoto@gmail.com>

## License

The MIT License (MIT)

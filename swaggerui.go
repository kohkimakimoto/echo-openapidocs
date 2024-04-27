package openapidocs

import (
	"bytes"
	"encoding/json"
	"github.com/labstack/echo/v4"
	htmltemplate "html/template"
	"net/http"
	"path"
	"strings"
)

// SwaggerUIConfig is the configuration for SwaggerUIDocumentsHandler to generate the OpenAPI documentation with Swagger UI.
// Some fields are Swagger UI configuration options.
// See https://github.com/swagger-api/swagger-ui/blob/master/docs/usage/configuration.md
type SwaggerUIConfig struct {
	// Spec is the OpenAPI specification.
	Spec string
	// SpecUrl is the URL of the OpenAPI specification. If Spec is not empty, SpecUrl is ignored.
	SpecUrl string
	// Title is the title of the page.
	Title string
	// Template is a template string for rendering the page with html/template.
	Template string

	// DeepLinking is the Swagger UI `deepLinking` configuration.
	DeepLinking bool
	// DisplayOperationId is the Swagger UI `DisplayOperationId` configuration.
	DisplayOperationId bool

	// TODO: Add more Redoc configuration options...
}

type swaggerUITemplateParams struct {
	SwaggerUIConfig
	BasePath               string
	SwaggerUIConfiguration htmltemplate.JS
}

type swaggerUIConfiguration struct {
	Url                string `json:"url"`
	DomId              string `json:"dom_id"`
	DeepLinking        bool   `json:"deepLinking,omitempty"`
	DisplayOperationId bool   `json:"displayOperationId,omitempty"`
}

var DefaultSwaggerUIConfig = SwaggerUIConfig{
	Spec:               "",
	SpecUrl:            "",
	Title:              "API documentation with Swagger UI",
	Template:           defaultSwaggerUITemplate,
	DeepLinking:        false,
	DisplayOperationId: false,
}

const defaultSwaggerUITemplate = `<html lang="en">
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <title>{{ .Title }}</title>
  <link rel="stylesheet" href="https://unpkg.com/swagger-ui-dist/swagger-ui.css" />
</head>
<body>
  <div id="swagger-ui"></div>
  <script src="https://unpkg.com/swagger-ui-dist/swagger-ui-bundle.js" crossorigin></script>
  <script>
	var configuration = {{ .SwaggerUIConfiguration }};
    window.onload = () => {
	  window.ui = SwaggerUIBundle(configuration);
    };
  </script>
</body>
</html>
`

func SwaggerUIDocumentsHandler(config SwaggerUIConfig) echo.HandlerFunc {
	if config.Template == "" {
		config.Template = DefaultSwaggerUIConfig.Template
	}
	if config.Title == "" {
		config.Title = DefaultSwaggerUIConfig.Title
	}

	useSpecUrl := false
	if config.Spec == "" {
		if config.SpecUrl == "" {
			panic("either Spec or SpecUrl must be set")
		}
		useSpecUrl = true
	}

	pageTmpl := htmltemplate.Must(htmltemplate.New("T").Parse(config.Template))

	return func(c echo.Context) error {
		p := c.Request().URL.Path

		// determine the base path
		relPath := c.Param("*")
		basePath := strings.TrimSuffix(p, relPath)

		var specUrl string
		if !useSpecUrl {
			specUrl = path.Join(basePath, "openapi-spec")
			if strings.HasSuffix(p, specUrl) {
				return c.Blob(http.StatusOK, "text/plain; charset=utf-8", []byte(config.Spec))
			}
		} else {
			specUrl = config.SpecUrl
		}

		if relPath != "" {
			// The document site only works with the base path.
			return c.Redirect(http.StatusFound, basePath)
		}

		swaggerUIConfiguration := swaggerUIConfiguration{
			Url:                specUrl,
			DomId:              "#swagger-ui",
			DeepLinking:        config.DeepLinking,
			DisplayOperationId: config.DisplayOperationId,
		}

		jsonDate, err := json.Marshal(swaggerUIConfiguration)
		if err != nil {
			return err
		}

		params := swaggerUITemplateParams{
			SwaggerUIConfig:        config,
			BasePath:               basePath,
			SwaggerUIConfiguration: htmltemplate.JS(jsonDate),
		}

		buf := new(bytes.Buffer)
		if err := pageTmpl.Execute(buf, params); err != nil {
			panic(err)
		}
		return c.HTML(http.StatusOK, buf.String())
	}
}

// SwaggerUIDocuments registers a handler for serving Swagger UI documents.
func SwaggerUIDocuments(e *echo.Echo, pathPrefix string, config SwaggerUIConfig) {
	e.GET(pathPrefix+"*", SwaggerUIDocumentsHandler(config))
}

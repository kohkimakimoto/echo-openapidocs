package openapidocs

import (
	"bytes"
	"encoding/json"
	htmltemplate "html/template"
	"io/fs"
	"net/http"
	"path"
	"strings"

	"github.com/labstack/echo/v5"
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
	// If empty, CDNSwaggerUITemplate is used.
	Template string
	// Assets is an optional filesystem that serves Swagger UI static files
	// (swagger-ui.css and swagger-ui-bundle.js) under "{pathPrefix}/_assets/*".
	// When set, SwaggerUIDocuments also registers AssetsHandler.
	// Prefer github.com/kohkimakimoto/echo-openapidocs/v5/embedded/swaggerui which sets
	// Assets and Template together.
	Assets fs.FS

	// DeepLinking is the Swagger UI `deepLinking` configuration.
	DeepLinking bool
	// DisplayOperationId is the Swagger UI `DisplayOperationId` configuration.
	DisplayOperationId bool

	// TODO: Add more Swagger UI configuration options...
}

type swaggerUITemplateParams struct {
	SwaggerUIConfig
	BasePath               string
	AssetsBaseURL          string
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
	Template:           CDNSwaggerUITemplate,
	DeepLinking:        false,
	DisplayOperationId: false,
}

// CDNSwaggerUITemplate loads Swagger UI from the CDN (unpkg).
const CDNSwaggerUITemplate = `<html lang="en">
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <title>{{ .Title }}</title>
  <link rel="stylesheet" href="https://unpkg.com/swagger-ui-dist@5.32.14/swagger-ui.css" />
</head>
<body>
  <div id="swagger-ui"></div>
  <script src="https://unpkg.com/swagger-ui-dist@5.32.14/swagger-ui-bundle.js" crossorigin></script>
  <script>
	var configuration = {{ .SwaggerUIConfiguration }};
    window.onload = () => {
	  window.ui = SwaggerUIBundle(configuration);
    };
  </script>
</body>
</html>
`

// SwaggerUIDocumentsHandler returns an echo.HandlerFunc to serve the OpenAPI documentation with Swagger UI.
func SwaggerUIDocumentsHandler(config SwaggerUIConfig) echo.HandlerFunc {
	if config.Template == "" {
		config.Template = CDNSwaggerUITemplate
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

	return func(c *echo.Context) error {
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
			AssetsBaseURL:          path.Join(basePath, AssetsPath),
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
// When config.Assets is set, it also registers AssetsHandler at "{pathPrefix}/_assets/*".
// If you use SwaggerUIDocumentsHandler directly, register AssetsHandler yourself when using Assets.
func SwaggerUIDocuments(e *echo.Echo, pathPrefix string, config SwaggerUIConfig) {
	registerAssets(e, pathPrefix, config.Assets)
	e.GET(pathPrefix+"*", SwaggerUIDocumentsHandler(config))
}

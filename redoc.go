package openapidocs

import (
	"bytes"
	"github.com/labstack/echo/v4"
	htmltemplate "html/template"
	"net/http"
	"path"
	"strings"
)

// RedocConfig is the configuration for RedocDocumentsHandler to generate the OpenAPI documentation with Redoc.
// Some fields are Redoc configuration options.
// See https://github.com/Redocly/redoc/blob/main/docs/config.md
type RedocConfig struct {
	// Spec is the OpenAPI specification.
	Spec string
	// SpecUrl is the URL of the OpenAPI specification. If Spec is not empty, SpecUrl is ignored.
	SpecUrl string
	// Title is the title of the page.
	Title string
	// Template is a template string for rendering the page with html/template.
	Template string

	// DisableSearch is the Redoc `disableSearch` configuration.
	DisableSearch bool
	// MinCharacterLengthToInitSearch is the Redoc `minCharacterLengthToInitSearch` configuration.
	MinCharacterLengthToInitSearch int

	// TODO: Add more Redoc configuration options...
}

type redocTemplateParams struct {
	RedocConfig
	BasePath string
	SpecUrl  string
}

var DefaultRedocConfig = RedocConfig{
	Spec:                           "",
	SpecUrl:                        "",
	Title:                          "API documentation with Redoc",
	Template:                       defaultRedocTemplate,
	MinCharacterLengthToInitSearch: 0,
}

const defaultRedocTemplate = `<html lang="en">
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <title>{{ .Title }}</title>
</head>
<body>
  <redoc
    spec-url="{{ .SpecUrl }}"
    {{- if .DisableSearch }}
	disable-search="true"
	{{- end }}
	{{- if ne .MinCharacterLengthToInitSearch 0 }}
	min-character-length-to-init-search="{{ .MinCharacterLengthToInitSearch }}"
	{{- end }}
  ></redoc>
  <script src="https://cdn.redoc.ly/redoc/latest/bundles/redoc.standalone.js"> </script>
</body>
</html>
`

func RedocDocumentsHandler(config RedocConfig) echo.HandlerFunc {
	if config.Template == "" {
		config.Template = DefaultRedocConfig.Template
	}
	if config.Title == "" {
		config.Title = DefaultRedocConfig.Title
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

		params := redocTemplateParams{
			RedocConfig: config,
			BasePath:    basePath,
			SpecUrl:     specUrl,
		}

		buf := new(bytes.Buffer)
		if err := pageTmpl.Execute(buf, params); err != nil {
			panic(err)
		}

		return c.HTML(http.StatusOK, buf.String())
	}
}

func RedocDocuments(e *echo.Echo, pathPrefix string, config RedocConfig) {
	e.GET(pathPrefix+"*", RedocDocumentsHandler(config))
}

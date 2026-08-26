package openapidocs

import (
	"bytes"
	htmltemplate "html/template"
	"io/fs"
	"net/http"
	"path"
	"strings"

	"github.com/labstack/echo/v5"
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
	// If empty, CDNRedocTemplate is used.
	Template string
	// Assets is an optional filesystem that serves Redoc static files
	// (redoc.standalone.js) under "{pathPrefix}/_assets/*".
	// When set, RedocDocuments also registers AssetsHandler.
	// Prefer github.com/kohkimakimoto/echo-openapidocs/v5/embedded/redoc which sets
	// Assets and Template together.
	Assets fs.FS

	// DisableSearch is the Redoc `disableSearch` configuration.
	DisableSearch bool
	// MinCharacterLengthToInitSearch is the Redoc `minCharacterLengthToInitSearch` configuration.
	MinCharacterLengthToInitSearch int

	// TODO: Add more Redoc configuration options...
}

type redocTemplateParams struct {
	RedocConfig
	BasePath      string
	AssetsBaseURL string
	SpecUrl       string
}

var DefaultRedocConfig = RedocConfig{
	Spec:                           "",
	SpecUrl:                        "",
	Title:                          "API documentation with Redoc",
	Template:                       CDNRedocTemplate,
	MinCharacterLengthToInitSearch: 0,
}

// CDNRedocTemplate loads Redoc from the CDN.
const CDNRedocTemplate = `<html lang="en">
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
  <script src="https://cdn.redoc.ly/redoc/v2.5.3/bundles/redoc.standalone.js"> </script>
</body>
</html>
`

// RedocDocumentsHandler returns an echo.HandlerFunc to serve the OpenAPI documentation with Redoc.
func RedocDocumentsHandler(config RedocConfig) echo.HandlerFunc {
	if config.Template == "" {
		config.Template = CDNRedocTemplate
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

		params := redocTemplateParams{
			RedocConfig:   config,
			BasePath:      basePath,
			AssetsBaseURL: path.Join(basePath, AssetsPath),
			SpecUrl:       specUrl,
		}

		buf := new(bytes.Buffer)
		if err := pageTmpl.Execute(buf, params); err != nil {
			panic(err)
		}

		return c.HTML(http.StatusOK, buf.String())
	}
}

// RedocDocuments registers a handler to serve the OpenAPI documentation with Redoc.
// When config.Assets is set, it also registers AssetsHandler at "{pathPrefix}/_assets/*".
// If you use RedocDocumentsHandler directly, register AssetsHandler yourself when using Assets.
func RedocDocuments(e *echo.Echo, pathPrefix string, config RedocConfig) {
	registerAssets(e, pathPrefix, config.Assets)
	e.GET(pathPrefix+"*", RedocDocumentsHandler(config))
}

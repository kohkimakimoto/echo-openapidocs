package openapidocs

import (
	"bytes"
	"github.com/labstack/echo/v4"
	htmltemplate "html/template"
	"net/http"
	"path"
	"strings"
)

// ElementsConfig is the configuration for ElementsDocumentsHandler to generate the OpenAPI documentation with Stoplight Elements.
// Some fields are Elements configuration options.
// See https://github.com/stoplightio/elements/blob/main/docs/getting-started/elements/elements-options.md
type ElementsConfig struct {
	// Spec is the OpenAPI specification.
	Spec string
	// SpecUrl is the URL of the OpenAPI specification. If Spec is not empty, SpecUrl is ignored.
	SpecUrl string
	// Title is the title of the page.
	Title string
	// Template is a template string for rendering the page with html/template.
	Template string

	// Router is the Elements `router` configuration.
	Router ElementsRouter
	// Layout is the Elements `layout` configuration.
	Layout ElementsLayout
	// HideInternal is the Elements `hideInternal` configuration.
	HideInternal bool
	// HideTryIt is the Elements `hideTryIt` configuration.
	HideTryIt bool
	// HideSchemas is the Elements `hideSchemas` configuration.
	HideSchemas bool
	// HideExport is the Elements `hideExport` configuration.
	HideExport bool
	// TryItCorsProxy is the Elements `tryItCorsProxy` configuration.
	TryItCorsProxy string
	// TryItCredentialsPolicy is the Elements `tryItCredentialsPolicy` configuration.
	TryItCredentialsPolicy ElementsTryItCredentialsPolicy
	// Logo is the Elements `logo` configuration.
	Logo string
}

type elementsTemplateParams struct {
	ElementsConfig
	BasePath          string
	ApiDescriptionUrl string
}

type ElementsRouter string

const (
	ElementsRouterHash    ElementsRouter = "hash"
	ElementsRouterHistory ElementsRouter = "history"
	ElementsRouterMemory  ElementsRouter = "memory"
)

type ElementsLayout string

const (
	ElementsLayoutSidebar    ElementsLayout = "sidebar"
	ElementsLayoutResponsive ElementsLayout = "responsive"
	ElementsLayoutStacked    ElementsLayout = "stacked"
)

type ElementsTryItCredentialsPolicy string

const (
	ElementsTryItCredentialsPolicyOmit ElementsTryItCredentialsPolicy = "omit"
	ElementsTryItCredentialsInclude    ElementsTryItCredentialsPolicy = "include"
	ElementsTryItCredentialsSameOrigin ElementsTryItCredentialsPolicy = "same-origin"
)

var DefaultElementsConfig = ElementsConfig{
	Spec:                   "",
	SpecUrl:                "",
	Title:                  "API documentation with Stoplight Elements",
	Template:               defaultElementsTemplate,
	Router:                 ElementsRouterHistory,
	Layout:                 ElementsLayoutSidebar,
	HideInternal:           false,
	HideTryIt:              false,
	HideSchemas:            false,
	HideExport:             false,
	TryItCorsProxy:         "",
	TryItCredentialsPolicy: ElementsTryItCredentialsPolicyOmit,
	Logo:                   "",
}

const defaultElementsTemplate = `<html lang="en">
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <title>{{ .Title }}</title>
  <script src="https://unpkg.com/@stoplight/elements/web-components.min.js"></script>
  <link rel="stylesheet" href="https://unpkg.com/@stoplight/elements/styles.min.css">
</head>
<body>
  <elements-api
    apiDescriptionUrl="{{ .ApiDescriptionUrl }}"
    {{- if ne .BasePath "" }}
    basePath="{{ .BasePath }}"
    {{- end }}
    {{- if .HideInternal }}
    hideInternal="true"
    {{- end }}
    {{- if .HideTryIt }}
    hideTryIt="true"
    {{- end }}
    {{- if .HideSchemas }}
    hideSchemas="true"
    {{- end }}
    {{- if .HideExport }}
    hideExport="true"
    {{- end }}
    {{- if ne .TryItCorsProxy "" }}
    tryItCorsProxy="{{ .TryItCorsProxy }}"
    {{- end }}
    {{- if and (ne .TryItCredentialsPolicy "") (ne .TryItCredentialsPolicy "omit") }}
    tryItCredentialsPolicy="{{ .TryItCredentialsPolicy }}"
    {{- end }}
    layout="{{ .Layout }}"
    {{- if ne .Logo "" }}
    logo="{{ .Logo }}"
    {{- end }}
    router="{{ .Router }}"
  />
</body>
</html>
`

// ElementsDocumentsHandler returns an echo.HandlerFunc to serve the OpenAPI documentation with Stoplight Elements.
func ElementsDocumentsHandler(config ElementsConfig) echo.HandlerFunc {
	if config.Template == "" {
		config.Template = DefaultElementsConfig.Template
	}
	if config.Router == "" {
		config.Router = DefaultElementsConfig.Router
	}
	if config.Layout == "" {
		config.Layout = DefaultElementsConfig.Layout
	}
	if config.Title == "" {
		config.Title = DefaultElementsConfig.Title
	}
	if config.TryItCredentialsPolicy == "" {
		config.TryItCredentialsPolicy = DefaultElementsConfig.TryItCredentialsPolicy
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

		if config.Router != ElementsRouterHistory && relPath != "" {
			// If the router is not history mode, the document site only works with the base path.
			return c.Redirect(http.StatusFound, basePath)
		}

		params := elementsTemplateParams{
			ElementsConfig:    config,
			BasePath:          basePath,
			ApiDescriptionUrl: specUrl,
		}

		buf := new(bytes.Buffer)
		if err := pageTmpl.Execute(buf, params); err != nil {
			panic(err)
		}

		return c.HTML(http.StatusOK, buf.String())
	}
}

// ElementsDocuments registers a handler to serve the OpenAPI documentation with Stoplight Elements.
func ElementsDocuments(e *echo.Echo, pathPrefix string, config ElementsConfig) {
	e.GET(pathPrefix+"*", ElementsDocumentsHandler(config))
}

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

// ScalarConfig is the configuration for ScalarDocumentsHandler to generate the OpenAPI documentation with Scalar.
// Some fields are Scalar configuration options.
// See https://github.com/scalar/scalar?tab=readme-ov-file#configuration
type ScalarConfig struct {
	// Spec is the OpenAPI specification.
	Spec string
	// SpecUrl is the URL of the OpenAPI specification. If Spec is not empty, SpecUrl is ignored.
	SpecUrl string
	// Title is the title of the page.
	Title string
	// Template is a template string for rendering the page with html/template.
	Template string

	// IsEditable is the Scalar `isEditable` configuration.
	IsEditable bool
	// ProxyUrl is the Scalar `proxyUrl` configuration.
	ProxyUrl string
	// DarkMode is the Scalar `darkMode` configuration.
	DarkMode bool
	// Layout is the Scalar `layout` configuration.
	Layout ScalarLayout
	// Theme is the Scalar `theme` configuration.
	Theme ScalarTheme
	// HideSidebar is the inverse of the Scalar `showSidebar` configuration.
	// If HideSidebar is true, the `showSidebar` configuration is false.
	// Scalar has a default value of `showSidebar` as true, so if you want to hide the sidebar, set this value to true.
	HideSidebar bool
	// SearchHotKey is the Scalar `searchHotKey` configuration.
	SearchHotKey string
}

type scalarTemplateParams struct {
	ScalarConfig
	BasePath                  string
	ApiReferenceConfiguration htmltemplate.JS
}

type apiReferenceConfiguration struct {
	IsEditable   bool                          `json:"isEditable,omitempty"`
	Spec         apiReferenceConfigurationSpec `json:"spec"`
	ProxyUrl     string                        `json:"proxyUrl,omitempty"`
	DarkMode     bool                          `json:"darkMode,omitempty"`
	Layout       ScalarLayout                  `json:"layout,omitempty"`
	Theme        ScalarTheme                   `json:"theme,omitempty"`
	ShowSidebar  bool                          `json:"showSidebar"` // the default value is true, that is the reason why it is not omitted
	SearchHotKey string                        `json:"searchHotKey,omitempty"`
}

type ScalarLayout string

const (
	ScalarLayoutModern  ScalarLayout = "modern"
	ScalarLayoutClassic ScalarLayout = "classic"
)

// ScalarTheme is the Scalar theme configuration.
// https://github.com/scalar/scalar?tab=readme-ov-file#themes
type ScalarTheme string

const (
	ScalarThemeAlternate  ScalarTheme = "alternate"
	ScalarThemeDefault    ScalarTheme = "default"
	ScalarThemeMoon       ScalarTheme = "moon"
	ScalarThemePurple     ScalarTheme = "purple"
	ScalarThemeSolarized  ScalarTheme = "solarized"
	ScalarThemeBluePlanet ScalarTheme = "bluePlanet"
	ScalarThemeSaturn     ScalarTheme = "saturn"
	ScalarThemeMars       ScalarTheme = "mars"
	ScalarThemeDeepSpace  ScalarTheme = "deepSpace"
	ScalarThemeNone       ScalarTheme = "none"
)

type apiReferenceConfigurationSpec struct {
	URL string `json:"url"`
}

var DefaultScalarConfig = ScalarConfig{
	Spec:         "",
	SpecUrl:      "",
	Title:        "API documentation with Scalar",
	Template:     defaultScalarTemplate,
	IsEditable:   false,
	ProxyUrl:     "",
	DarkMode:     false,
	Layout:       ScalarLayoutModern,
	Theme:        ScalarThemeDefault,
	HideSidebar:  false,
	SearchHotKey: "",
}

const defaultScalarTemplate = `<html lang="en">
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <title>{{ .Title }}</title>
</head>
<body>
  <script id="api-reference" type="application/json"></script>
  <script>
    var configuration = {{ .ApiReferenceConfiguration }};
    var apiReference = document.getElementById('api-reference');
    apiReference.dataset.configuration = JSON.stringify(configuration);
  </script>
  <script src="https://cdn.jsdelivr.net/npm/@scalar/api-reference"></script>
</body>
</html>
`

// ScalarDocumentsHandler returns an echo.HandlerFunc to serve the OpenAPI documentation with Scalar.
func ScalarDocumentsHandler(config ScalarConfig) echo.HandlerFunc {
	if config.Template == "" {
		config.Template = DefaultScalarConfig.Template
	}
	if config.Title == "" {
		config.Title = DefaultScalarConfig.Title
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

		apiReferenceConfiguration := apiReferenceConfiguration{
			IsEditable: config.IsEditable,
			Spec: apiReferenceConfigurationSpec{
				URL: specUrl,
			},
			ProxyUrl:     config.ProxyUrl,
			DarkMode:     config.DarkMode,
			Layout:       config.Layout,
			Theme:        config.Theme,
			ShowSidebar:  !config.HideSidebar,
			SearchHotKey: config.SearchHotKey,
		}

		jsonDate, err := json.Marshal(apiReferenceConfiguration)
		if err != nil {
			return err
		}

		params := scalarTemplateParams{
			ScalarConfig:              config,
			BasePath:                  basePath,
			ApiReferenceConfiguration: htmltemplate.JS(jsonDate),
		}

		buf := new(bytes.Buffer)
		if err := pageTmpl.Execute(buf, params); err != nil {
			panic(err)
		}
		return c.HTML(http.StatusOK, buf.String())
	}
}

// ScalarDocuments registers a handler to serve the OpenAPI documentation with Scalar.
func ScalarDocuments(e *echo.Echo, pathPrefix string, config ScalarConfig) {
	e.GET(pathPrefix+"*", ScalarDocumentsHandler(config))
}

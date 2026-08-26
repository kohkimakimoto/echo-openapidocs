// Package swaggerui provides optionally importable bundled Swagger UI static assets.
//
// Importing this package increases binary size because the assets are embedded
// with go:embed. Applications that are fine with CDN-hosted assets should use
// the root openapidocs package only.
package swaggerui

import (
	"embed"

	openapidocs "github.com/kohkimakimoto/echo-openapidocs/v4"
	"github.com/labstack/echo/v4"
)

//go:embed swagger-ui.css swagger-ui-bundle.js
var FS embed.FS

// Template loads Swagger UI from Assets served by openapidocs.AssetsHandler under
// "{pathPrefix}/_assets/*".
const Template = `<html lang="en">
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <title>{{ .Title }}</title>
  <link rel="stylesheet" href="{{ .AssetsBaseURL }}/swagger-ui.css" />
</head>
<body>
  <div id="swagger-ui"></div>
  <script src="{{ .AssetsBaseURL }}/swagger-ui-bundle.js" crossorigin></script>
  <script>
	var configuration = {{ .SwaggerUIConfiguration }};
    window.onload = () => {
	  window.ui = SwaggerUIBundle(configuration);
    };
  </script>
</body>
</html>
`

// Documents registers Swagger UI documentation using the bundled static assets
// instead of the CDN.
func Documents(e *echo.Echo, pathPrefix string, config openapidocs.SwaggerUIConfig) {
	config.Assets = FS
	if config.Template == "" {
		config.Template = Template
	}
	openapidocs.SwaggerUIDocuments(e, pathPrefix, config)
}

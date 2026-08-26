// Package scalar provides optionally importable bundled Scalar static assets.
//
// Importing this package increases binary size because the assets are embedded
// with go:embed. Applications that are fine with CDN-hosted assets should use
// the root openapidocs package only.
package scalar

import (
	"embed"

	openapidocs "github.com/kohkimakimoto/echo-openapidocs/v5"
	"github.com/labstack/echo/v5"
)

//go:embed standalone.js
var FS embed.FS

// Template loads Scalar from Assets served by openapidocs.AssetsHandler under
// "{pathPrefix}/_assets/*".
const Template = `<html lang="en">
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
  <script src="{{ .AssetsBaseURL }}/standalone.js"></script>
</body>
</html>
`

// Documents registers Scalar documentation using the bundled static assets
// instead of the CDN.
func Documents(e *echo.Echo, pathPrefix string, config openapidocs.ScalarConfig) {
	config.Assets = FS
	if config.Template == "" {
		config.Template = Template
	}
	openapidocs.ScalarDocuments(e, pathPrefix, config)
}

// Package redoc provides optionally importable bundled Redoc static assets.
//
// Importing this package increases binary size because the assets are embedded
// with go:embed. Applications that are fine with CDN-hosted assets should use
// the root openapidocs package only.
package redoc

import (
	"embed"

	openapidocs "github.com/kohkimakimoto/echo-openapidocs/v5"
	"github.com/labstack/echo/v5"
)

//go:embed redoc.standalone.js
var FS embed.FS

// Template loads Redoc from Assets served by openapidocs.AssetsHandler under
// "{pathPrefix}/_assets/*".
const Template = `<html lang="en">
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
  <script src="{{ .AssetsBaseURL }}/redoc.standalone.js"> </script>
</body>
</html>
`

// Documents registers Redoc documentation using the bundled static assets
// instead of the CDN.
func Documents(e *echo.Echo, pathPrefix string, config openapidocs.RedocConfig) {
	config.Assets = FS
	if config.Template == "" {
		config.Template = Template
	}
	openapidocs.RedocDocuments(e, pathPrefix, config)
}

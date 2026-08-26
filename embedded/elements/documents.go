// Package elements provides optionally importable bundled Stoplight Elements
// static assets.
//
// Importing this package increases binary size because the assets are embedded
// with go:embed. Applications that are fine with CDN-hosted assets should use
// the root openapidocs package only.
package elements

import (
	"embed"

	openapidocs "github.com/kohkimakimoto/echo-openapidocs/v5"
	"github.com/labstack/echo/v5"
)

//go:embed web-components.min.js styles.min.css
var FS embed.FS

// Template loads Stoplight Elements from Assets served by
// openapidocs.AssetsHandler under "{pathPrefix}/_assets/*".
const Template = `<html lang="en">
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <title>{{ .Title }}</title>
  <script src="{{ .AssetsBaseURL }}/web-components.min.js"></script>
  <link rel="stylesheet" href="{{ .AssetsBaseURL }}/styles.min.css">
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

// Documents registers Stoplight Elements documentation using the bundled
// static assets instead of the CDN.
func Documents(e *echo.Echo, pathPrefix string, config openapidocs.ElementsConfig) {
	config.Assets = FS
	if config.Template == "" {
		config.Template = Template
	}
	openapidocs.ElementsDocuments(e, pathPrefix, config)
}

// Package embedded provides optionally importable bundled static assets for
// OpenAPI documentation generators.
//
// Importing this package increases binary size because the assets are embedded
// with go:embed. Applications that are fine with CDN-hosted assets should use
// the root openapidocs package only.
package embedded

import (
	"embed"
	"io/fs"

	openapidocs "github.com/kohkimakimoto/echo-openapidocs/v4"
	"github.com/labstack/echo/v4"
)

//go:embed elements/web-components.min.js elements/styles.min.css
var elementsFS embed.FS

// Elements is the filesystem that contains Stoplight Elements static files
// (web-components.min.js and styles.min.css) at the FS root.
var Elements = mustSubFS(elementsFS, "elements")

// ElementsTemplate loads Stoplight Elements from Assets served by
// openapidocs.ElementsAssetsHandler under "{pathPrefix}/_assets/*".
const ElementsTemplate = `<html lang="en">
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

func mustSubFS(filesystem fs.FS, dir string) fs.FS {
	sub, err := fs.Sub(filesystem, dir)
	if err != nil {
		panic(err)
	}
	return sub
}

// ElementsDocuments registers Stoplight Elements documentation using the
// bundled static assets instead of the CDN.
func ElementsDocuments(e *echo.Echo, pathPrefix string, config openapidocs.ElementsConfig) {
	config.Assets = Elements
	if config.Template == "" {
		config.Template = ElementsTemplate
	}
	openapidocs.ElementsDocuments(e, pathPrefix, config)
}

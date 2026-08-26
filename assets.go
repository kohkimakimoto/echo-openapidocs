package openapidocs

import (
	"io"
	"io/fs"
	"net/http"
	"path"
	"strings"

	"github.com/labstack/echo/v5"
)

// AssetsPath is the URL path segment used to serve Config.Assets.
// The leading underscore reduces the chance of colliding with OpenAPI document paths.
const AssetsPath = "_assets"

// AssetsHandler returns an echo.HandlerFunc that serves static files from filesystem
// under the request path after pathPrefix (typically "{docsPrefix}/_assets").
func AssetsHandler(filesystem fs.FS) echo.HandlerFunc {
	return func(c *echo.Context) error {
		name := c.Param("*")
		name = strings.TrimPrefix(name, "/")
		if name == "" || strings.Contains(name, "..") {
			return echo.ErrNotFound
		}

		f, err := filesystem.Open(name)
		if err != nil {
			return echo.ErrNotFound
		}
		defer f.Close()

		stat, err := f.Stat()
		if err != nil {
			return err
		}
		if stat.IsDir() {
			return echo.ErrNotFound
		}

		rs, ok := f.(io.ReadSeeker)
		if !ok {
			b, err := io.ReadAll(f)
			if err != nil {
				return err
			}
			return c.Blob(http.StatusOK, contentTypeByExt(name), b)
		}

		http.ServeContent(c.Response(), c.Request(), stat.Name(), stat.ModTime(), rs)
		return nil
	}
}

func contentTypeByExt(name string) string {
	switch strings.ToLower(path.Ext(name)) {
	case ".js":
		return "application/javascript; charset=utf-8"
	case ".css":
		return "text/css; charset=utf-8"
	case ".map":
		return "application/json; charset=utf-8"
	default:
		return "application/octet-stream"
	}
}

func registerAssets(e *echo.Echo, pathPrefix string, filesystem fs.FS) {
	if filesystem == nil {
		return
	}
	assetsPrefix := strings.TrimSuffix(pathPrefix, "/") + "/" + AssetsPath + "/"
	e.GET(assetsPrefix+"*", AssetsHandler(filesystem))
}

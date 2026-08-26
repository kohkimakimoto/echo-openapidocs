package openapidocs_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v5"

	openapidocs "github.com/kohkimakimoto/echo-openapidocs/v5"
	"github.com/kohkimakimoto/echo-openapidocs/v5/embedded/swaggerui"
)

func TestSwaggerUIDocuments_CDN(t *testing.T) {
	e := echo.New()
	openapidocs.SwaggerUIDocuments(e, "/docs", openapidocs.SwaggerUIConfig{
		SpecUrl: "https://example.com/openapi.yaml",
		Title:   "Test API",
	})

	req := httptest.NewRequest(http.MethodGet, "/docs", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}
	body := rec.Body.String()
	if !strings.Contains(body, "unpkg.com/swagger-ui-dist@5.32.14/swagger-ui-bundle.js") {
		t.Fatalf("expected CDN script URL in body, got:\n%s", body)
	}
	if !strings.Contains(body, "unpkg.com/swagger-ui-dist@5.32.14/swagger-ui.css") {
		t.Fatalf("expected CDN stylesheet URL in body, got:\n%s", body)
	}
	if strings.Contains(body, "/docs/_assets/") {
		t.Fatalf("CDN template should not reference local assets")
	}
}

func TestSwaggerUIDocuments_Embedded(t *testing.T) {
	e := echo.New()
	swaggerui.Documents(e, "/docs", openapidocs.SwaggerUIConfig{
		SpecUrl: "https://example.com/openapi.yaml",
		Title:   "Test API",
	})

	req := httptest.NewRequest(http.MethodGet, "/docs", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}
	body := rec.Body.String()
	if !strings.Contains(body, `/docs/_assets/swagger-ui-bundle.js`) {
		t.Fatalf("expected embedded script URL in body, got:\n%s", body)
	}
	if !strings.Contains(body, `/docs/_assets/swagger-ui.css`) {
		t.Fatalf("expected embedded stylesheet URL in body, got:\n%s", body)
	}

	for _, name := range []string{"swagger-ui-bundle.js", "swagger-ui.css"} {
		req = httptest.NewRequest(http.MethodGet, "/docs/_assets/"+name, nil)
		rec = httptest.NewRecorder()
		e.ServeHTTP(rec, req)
		if rec.Code != http.StatusOK {
			t.Fatalf("%s: status = %d, want %d", name, rec.Code, http.StatusOK)
		}
		b, err := io.ReadAll(rec.Body)
		if err != nil {
			t.Fatal(err)
		}
		if len(b) == 0 {
			t.Fatalf("%s: empty body", name)
		}
	}
}

func TestSwaggerUIAssetsHandler_NotFound(t *testing.T) {
	e := echo.New()
	openapidocs.SwaggerUIDocuments(e, "/docs", openapidocs.SwaggerUIConfig{
		SpecUrl:  "https://example.com/openapi.yaml",
		Assets:   swaggerui.FS,
		Template: swaggerui.Template,
	})

	req := httptest.NewRequest(http.MethodGet, "/docs/_assets/missing.js", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	if rec.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusNotFound)
	}
}

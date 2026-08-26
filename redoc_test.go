package openapidocs_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"

	openapidocs "github.com/kohkimakimoto/echo-openapidocs/v4"
	"github.com/kohkimakimoto/echo-openapidocs/v4/embedded/redoc"
)

func TestRedocDocuments_CDN(t *testing.T) {
	e := echo.New()
	openapidocs.RedocDocuments(e, "/docs", openapidocs.RedocConfig{
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
	if !strings.Contains(body, "cdn.redoc.ly/redoc/v2.5.3/bundles/redoc.standalone.js") {
		t.Fatalf("expected CDN script URL in body, got:\n%s", body)
	}
	if strings.Contains(body, "/docs/_assets/") {
		t.Fatalf("CDN template should not reference local assets")
	}
}

func TestRedocDocuments_Embedded(t *testing.T) {
	e := echo.New()
	redoc.Documents(e, "/docs", openapidocs.RedocConfig{
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
	if !strings.Contains(body, `/docs/_assets/redoc.standalone.js`) {
		t.Fatalf("expected embedded script URL in body, got:\n%s", body)
	}

	req = httptest.NewRequest(http.MethodGet, "/docs/_assets/redoc.standalone.js", nil)
	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}
	b, err := io.ReadAll(rec.Body)
	if err != nil {
		t.Fatal(err)
	}
	if len(b) == 0 {
		t.Fatal("empty body")
	}
}

func TestRedocAssetsHandler_NotFound(t *testing.T) {
	e := echo.New()
	openapidocs.RedocDocuments(e, "/docs", openapidocs.RedocConfig{
		SpecUrl:  "https://example.com/openapi.yaml",
		Assets:   redoc.FS,
		Template: redoc.Template,
	})

	req := httptest.NewRequest(http.MethodGet, "/docs/_assets/missing.js", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	if rec.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusNotFound)
	}
}

package openapidocs_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v5"

	openapidocs "github.com/kohkimakimoto/echo-openapidocs/v5"
	"github.com/kohkimakimoto/echo-openapidocs/v5/embedded/scalar"
)

func TestScalarDocuments_CDN(t *testing.T) {
	e := echo.New()
	openapidocs.ScalarDocuments(e, "/docs", openapidocs.ScalarConfig{
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
	if !strings.Contains(body, "cdn.jsdelivr.net/npm/@scalar/api-reference@1.66.1") {
		t.Fatalf("expected CDN script URL in body, got:\n%s", body)
	}
	if strings.Contains(body, "/docs/_assets/") {
		t.Fatalf("CDN template should not reference local assets")
	}
}

func TestScalarDocuments_Embedded(t *testing.T) {
	e := echo.New()
	scalar.Documents(e, "/docs", openapidocs.ScalarConfig{
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
	if !strings.Contains(body, `/docs/_assets/standalone.js`) {
		t.Fatalf("expected embedded script URL in body, got:\n%s", body)
	}

	req = httptest.NewRequest(http.MethodGet, "/docs/_assets/standalone.js", nil)
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

func TestScalarAssetsHandler_NotFound(t *testing.T) {
	e := echo.New()
	openapidocs.ScalarDocuments(e, "/docs", openapidocs.ScalarConfig{
		SpecUrl:  "https://example.com/openapi.yaml",
		Assets:   scalar.FS,
		Template: scalar.Template,
	})

	req := httptest.NewRequest(http.MethodGet, "/docs/_assets/missing.js", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	if rec.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusNotFound)
	}
}

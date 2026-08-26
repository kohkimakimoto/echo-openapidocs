package openapidocs_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"

	openapidocs "github.com/kohkimakimoto/echo-openapidocs/v4"
	"github.com/kohkimakimoto/echo-openapidocs/v4/embedded"
)

func TestElementsDocuments_CDN(t *testing.T) {
	e := echo.New()
	openapidocs.ElementsDocuments(e, "/docs", openapidocs.ElementsConfig{
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
	if !strings.Contains(body, "unpkg.com/@stoplight/elements@9.0.24") {
		t.Fatalf("expected CDN script URL in body, got:\n%s", body)
	}
	if strings.Contains(body, "/docs/_assets/") {
		t.Fatalf("CDN template should not reference local assets")
	}
}

func TestElementsDocuments_Embedded(t *testing.T) {
	e := echo.New()
	embedded.ElementsDocuments(e, "/docs", openapidocs.ElementsConfig{
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
	if !strings.Contains(body, `/docs/_assets/web-components.min.js`) {
		t.Fatalf("expected embedded script URL in body, got:\n%s", body)
	}
	if !strings.Contains(body, `/docs/_assets/styles.min.css`) {
		t.Fatalf("expected embedded stylesheet URL in body, got:\n%s", body)
	}

	for _, name := range []string{"web-components.min.js", "styles.min.css"} {
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

func TestElementsAssetsHandler_NotFound(t *testing.T) {
	e := echo.New()
	openapidocs.ElementsDocuments(e, "/docs", openapidocs.ElementsConfig{
		SpecUrl:  "https://example.com/openapi.yaml",
		Assets:   embedded.Elements,
		Template: embedded.ElementsTemplate,
	})

	req := httptest.NewRequest(http.MethodGet, "/docs/_assets/missing.js", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	if rec.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusNotFound)
	}
}

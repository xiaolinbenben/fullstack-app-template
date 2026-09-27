package httpapi

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"testing/fstest"

	"fullstack-app-template/server/internal/config"
)

func testHandler() http.Handler {
	webFS := fstest.MapFS{
		"public/index.html":    &fstest.MapFile{Data: []byte("<html>public</html>")},
		"public/assets/app.js": &fstest.MapFile{Data: []byte("console.log('app')")},
		"admin/index.html":     &fstest.MapFile{Data: []byte("<html>admin</html>")},
	}
	return New(config.Config{Addr: ":8000"}, webFS)
}

func TestHealthz(t *testing.T) {
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	testHandler().ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("healthz status = %d, want %d", recorder.Code, http.StatusOK)
	}
	if got := recorder.Header().Get("Content-Type"); got != "application/json; charset=utf-8" {
		t.Fatalf("content type = %q", got)
	}
	if body := recorder.Body.String(); body != `{"success":true,"message":"","data":{"status":"ok"}}` {
		t.Fatalf("unexpected body: %s", body)
	}
}

func TestFrontendFallbacks(t *testing.T) {
	tests := []struct {
		path string
		want string
	}{
		{path: "/", want: "public"},
		{path: "/app/settings", want: "public"},
		{path: "/admin/", want: "admin"},
		{path: "/admin/users", want: "admin"},
		{path: "/assets/app.js", want: "console.log('app')"},
	}

	for _, test := range tests {
		t.Run(test.path, func(t *testing.T) {
			recorder := httptest.NewRecorder()
			request := httptest.NewRequest(http.MethodGet, test.path, nil)
			testHandler().ServeHTTP(recorder, request)
			if recorder.Code != http.StatusOK {
				t.Fatalf("status = %d", recorder.Code)
			}
			if body := recorder.Body.String(); body != test.want && !contains(body, test.want) {
				t.Fatalf("body %q does not contain %q", body, test.want)
			}
		})
	}
}

func TestUnknownAPI(t *testing.T) {
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/api/unknown", nil)
	testHandler().ServeHTTP(recorder, request)

	if recorder.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusNotFound)
	}
	if body := recorder.Body.String(); body != `{"success":false,"message":"接口不存在","data":{}}` {
		t.Fatalf("unexpected body: %s", body)
	}
}

func TestHealthzRejectsUnsupportedMethod(t *testing.T) {
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/healthz", nil)
	testHandler().ServeHTTP(recorder, request)

	if recorder.Code != http.StatusMethodNotAllowed {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusMethodNotAllowed)
	}
	if body := recorder.Body.String(); body != `{"success":false,"message":"请求方法不支持","data":{}}` {
		t.Fatalf("unexpected body: %s", body)
	}
}

func contains(value, expected string) bool {
	for i := 0; i+len(expected) <= len(value); i++ {
		if value[i:i+len(expected)] == expected {
			return true
		}
	}
	return false
}

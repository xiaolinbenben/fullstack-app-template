package httpapi

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestWoodfishStartsAtZeroAndIncrements(t *testing.T) {
	handler := testHandler(t)
	if got := woodfishCount(t, handler); got != 0 {
		t.Fatalf("initial count = %d, want 0", got)
	}
	if got := knockWoodfish(t, handler); got != 1 {
		t.Fatalf("first knock = %d, want 1", got)
	}
	if got := knockWoodfish(t, handler); got != 2 {
		t.Fatalf("second knock = %d, want 2", got)
	}
	if got := woodfishCount(t, handler); got != 2 {
		t.Fatalf("saved count = %d, want 2", got)
	}
}

func TestWoodfishRejectsOtherMethods(t *testing.T) {
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPut, "/api/woodfish", nil)
	testHandler(t).ServeHTTP(recorder, request)
	if recorder.Code != http.StatusMethodNotAllowed {
		t.Fatalf("status = %d, body %s", recorder.Code, recorder.Body.String())
	}
}

func TestAdminSession(t *testing.T) {
	handler := testHandler(t)

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/api/admin/login", strings.NewReader(`{"username":"admin","password":"admin"}`))
	request.Header.Set("Content-Type", "application/json")
	handler.ServeHTTP(recorder, request)
	if recorder.Code != http.StatusUnauthorized || !strings.Contains(recorder.Body.String(), "账号或密码错误") {
		t.Fatalf("wrong password status = %d, body %s", recorder.Code, recorder.Body.String())
	}

	recorder = httptest.NewRecorder()
	request = httptest.NewRequest(http.MethodGet, "/api/admin/session", nil)
	handler.ServeHTTP(recorder, request)
	if recorder.Code != http.StatusUnauthorized {
		t.Fatalf("anonymous session status = %d", recorder.Code)
	}

	recorder = httptest.NewRecorder()
	request = httptest.NewRequest(http.MethodPost, "/api/admin/login", strings.NewReader(`{"username":"admin","password":"admin123"}`))
	request.Header.Set("Content-Type", "application/json")
	handler.ServeHTTP(recorder, request)
	if recorder.Code != http.StatusOK {
		t.Fatalf("login status = %d, body %s", recorder.Code, recorder.Body.String())
	}
	cookies := recorder.Result().Cookies()
	if len(cookies) == 0 {
		t.Fatal("login did not set a cookie")
	}

	recorder = httptest.NewRecorder()
	request = httptest.NewRequest(http.MethodGet, "/api/admin/session", nil)
	for _, cookie := range cookies {
		request.AddCookie(cookie)
	}
	handler.ServeHTTP(recorder, request)
	if recorder.Code != http.StatusOK || !strings.Contains(recorder.Body.String(), `"username":"admin"`) {
		t.Fatalf("session status = %d, body %s", recorder.Code, recorder.Body.String())
	}

	recorder = httptest.NewRecorder()
	request = httptest.NewRequest(http.MethodPost, "/api/admin/logout", nil)
	for _, cookie := range cookies {
		request.AddCookie(cookie)
	}
	handler.ServeHTTP(recorder, request)
	if recorder.Code != http.StatusOK {
		t.Fatalf("logout status = %d", recorder.Code)
	}
	cleared := recorder.Result().Cookies()

	recorder = httptest.NewRecorder()
	request = httptest.NewRequest(http.MethodGet, "/api/admin/session", nil)
	for _, cookie := range cleared {
		request.AddCookie(cookie)
	}
	handler.ServeHTTP(recorder, request)
	if recorder.Code != http.StatusUnauthorized || !strings.Contains(recorder.Body.String(), "未登录") {
		t.Fatalf("session after logout status = %d, body %s", recorder.Code, recorder.Body.String())
	}
}

func woodfishCount(t *testing.T, handler http.Handler) int64 {
	t.Helper()
	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/api/woodfish", nil))
	return decodeCount(t, recorder)
}

func knockWoodfish(t *testing.T, handler http.Handler) int64 {
	t.Helper()
	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, httptest.NewRequest(http.MethodPost, "/api/woodfish", nil))
	return decodeCount(t, recorder)
}

func decodeCount(t *testing.T, recorder *httptest.ResponseRecorder) int64 {
	t.Helper()
	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d, body %s", recorder.Code, recorder.Body.String())
	}
	var payload struct {
		Success bool   `json:"success"`
		Message string `json:"message"`
		Data    struct {
			Count int64 `json:"count"`
		} `json:"data"`
	}
	if err := json.Unmarshal(recorder.Body.Bytes(), &payload); err != nil {
		t.Fatal(err)
	}
	if !payload.Success || payload.Message != "" {
		t.Fatalf("body = %s", recorder.Body.String())
	}
	return payload.Data.Count
}

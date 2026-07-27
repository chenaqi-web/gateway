package controller

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"backend/gateway/internal/config"
)

func newTestRefreshCookieManager(t *testing.T, secure bool) *RefreshCookieManager {
	t.Helper()
	manager, err := NewRefreshCookieManager(config.AuthConfig{RefreshExpire: "168h", CookieSecure: secure})
	if err != nil {
		t.Fatalf("NewRefreshCookieManager() error = %v", err)
	}
	manager.now = func() time.Time { return time.Date(2026, 7, 27, 12, 0, 0, 0, time.UTC) }
	return manager
}

func TestRefreshCookieManagerSetGetAndClear(t *testing.T) {
	manager := newTestRefreshCookieManager(t, true)
	recorder := httptest.NewRecorder()
	manager.Set(recorder, "refresh-value")

	cookies := recorder.Result().Cookies()
	if len(cookies) != 1 {
		t.Fatalf("Set() cookies = %d, want 1", len(cookies))
	}
	cookie := cookies[0]
	if cookie.Name != refreshCookieName || cookie.Value != "refresh-value" || cookie.Path != refreshCookiePath {
		t.Fatalf("Set() cookie identity = %+v", cookie)
	}
	if !cookie.HttpOnly || !cookie.Secure || cookie.SameSite != http.SameSiteLaxMode || cookie.Domain != "" {
		t.Fatalf("Set() cookie security = %+v", cookie)
	}
	if cookie.MaxAge != 7*24*60*60 || !cookie.Expires.Equal(manager.now().Add(7*24*time.Hour)) {
		t.Fatalf("Set() cookie lifetime = %+v", cookie)
	}

	request := httptest.NewRequest(http.MethodPost, "/api/v1/auth/refresh", nil)
	request.AddCookie(cookie)
	value, err := manager.Get(request)
	if err != nil || value != "refresh-value" {
		t.Fatalf("Get() = %q, %v", value, err)
	}

	clearRecorder := httptest.NewRecorder()
	manager.Clear(clearRecorder)
	cleared := clearRecorder.Result().Cookies()[0]
	if cleared.Name != cookie.Name || cleared.Path != cookie.Path || cleared.Domain != cookie.Domain ||
		cleared.Secure != cookie.Secure || cleared.SameSite != cookie.SameSite || !cleared.HttpOnly {
		t.Fatalf("Clear() attributes differ: set=%+v clear=%+v", cookie, cleared)
	}
	if cleared.Value != "" || cleared.MaxAge != -1 || !cleared.Expires.Before(manager.now()) {
		t.Fatalf("Clear() cookie = %+v", cleared)
	}
}

func TestRefreshCookieManagerGetRejectsMissingOrEmptyCookie(t *testing.T) {
	manager := newTestRefreshCookieManager(t, false)
	request := httptest.NewRequest(http.MethodPost, "/api/v1/auth/refresh", nil)
	if _, err := manager.Get(request); err == nil {
		t.Fatal("Get() expected missing-cookie error")
	}
	request.AddCookie(&http.Cookie{Name: refreshCookieName, Value: ""})
	if _, err := manager.Get(request); err == nil {
		t.Fatal("Get() expected empty-cookie error")
	}
}

package utils

import (
	"backend/gateway/internal/config"
	"net/http"
	"time"
)

// 新建token

// 解析token

const (
	refreshCookieName = "refresh_token"
	refreshCookiePath = "/api/v1/auth"
)

type RefreshCookieManager struct {
	secure bool
	ttl    time.Duration
	now    func() time.Time
}

func NewRefreshCookieManager(cfg config.AuthConfig) (*RefreshCookieManager, error) {
	ttl, err := cfg.RefreshDuration()
	if err != nil {
		return nil, err
	}
	return &RefreshCookieManager{
		secure: cfg.CookieSecure,
		ttl:    ttl,
		now:    time.Now,
	}, nil
}

func (m *RefreshCookieManager) Set(writer http.ResponseWriter, refreshToken string) {
	http.SetCookie(writer, &http.Cookie{
		Name:     refreshCookieName,
		Value:    refreshToken,
		Path:     refreshCookiePath,
		Expires:  m.now().UTC().Add(m.ttl),
		MaxAge:   int(m.ttl / time.Second),
		HttpOnly: true,
		Secure:   m.secure,
		SameSite: http.SameSiteLaxMode,
	})
}

func (m *RefreshCookieManager) Get(request *http.Request) (string, error) {
	cookie, err := request.Cookie(refreshCookieName)
	if err != nil {
		return "", err
	}
	if cookie.Value == "" {
		return "", http.ErrNoCookie
	}
	return cookie.Value, nil
}

func (m *RefreshCookieManager) Clear(writer http.ResponseWriter) {
	http.SetCookie(writer, &http.Cookie{
		Name:     refreshCookieName,
		Value:    "",
		Path:     refreshCookiePath,
		Expires:  time.Unix(1, 0).UTC(),
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   m.secure,
		SameSite: http.SameSiteLaxMode,
	})
}

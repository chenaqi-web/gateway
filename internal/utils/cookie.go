package utils

import (
	"net/http"
	"time"

	"backend/gateway/internal/config"
)

const (
	refreshCookieName = "refresh_token"
	refreshCookiePath = "/api/v1"
)

func SetRefreshCookie(writer http.ResponseWriter, refreshToken string, cfg config.AuthConfig) error {
	http.SetCookie(writer, &http.Cookie{
		Name:  refreshCookieName,
		Value: refreshToken,
		Path:  refreshCookiePath,
		//Expires:  time.Now().UTC().Add(ttl), 直接用max就行
		MaxAge:   int(cfg.RefreshDuration()),
		HttpOnly: true,
		Secure:   cfg.CookieSecure,
		SameSite: http.SameSiteLaxMode,
	})
	return nil
}

func RefreshTokenFromCookie(request *http.Request) (string, error) {
	cookie, err := request.Cookie(refreshCookieName)
	if err != nil {
		return "", err
	}
	if cookie.Value == "" {
		return "", http.ErrNoCookie
	}
	return cookie.Value, nil
}

func ClearRefreshCookie(writer http.ResponseWriter, cfg config.AuthConfig) {
	http.SetCookie(writer, &http.Cookie{
		Name:     refreshCookieName,
		Value:    "",
		Path:     refreshCookiePath,
		Expires:  time.Unix(1, 0).UTC(),
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   cfg.CookieSecure,
		SameSite: http.SameSiteLaxMode,
	})
}

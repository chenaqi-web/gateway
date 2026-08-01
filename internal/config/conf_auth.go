package config

import (
	"fmt"
	"time"
)

const (
	defaultAccessExpire  = 15 * time.Minute
	defaultRefreshExpire = 7 * 24 * time.Hour
	minJWTSecretBytes    = 32
)

type AuthConfig struct {
	JWTSecret     string `yaml:"jwt_secret"`
	AccessExpire  string `yaml:"access_expire"`
	RefreshExpire string `yaml:"refresh_expire"`
	CookieSecure  bool   `yaml:"cookie_secure"`
}

func (c AuthConfig) JWTSigningKey() ([]byte, error) {
	key := []byte(c.JWTSecret)
	if len(key) < minJWTSecretBytes {
		return nil, fmt.Errorf("auth jwt_secret must be at least %d bytes", minJWTSecretBytes)
	}
	return key, nil
}

func (c AuthConfig) AccessDuration() (time.Duration, error) {
	if c.AccessExpire == "" {
		return defaultAccessExpire, nil
	}
	duration, err := time.ParseDuration(c.AccessExpire)
	if err != nil || duration <= 0 {
		return 0, fmt.Errorf("auth access_expire must be a positive duration")
	}
	return duration, nil
}

func (c AuthConfig) RefreshDuration() (time.Duration, error) {
	if c.RefreshExpire == "" {
		return defaultRefreshExpire, nil
	}
	duration, err := time.ParseDuration(c.RefreshExpire)
	if err != nil || duration <= 0 {
		return 0, fmt.Errorf("auth refresh_expire must be a positive duration")
	}
	return duration, nil
}

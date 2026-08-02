package config

import (
	"time"
)

type AuthConfig struct {
	JWTSecret     string `yaml:"jwt_secret"`
	AccessExpire  int    `yaml:"access_expire"`
	RefreshExpire int    `yaml:"refresh_expire"`
	CookieSecure  bool   `yaml:"cookie_secure"`
}

func (c AuthConfig) JWTSigningKey() []byte {
	key := []byte(c.JWTSecret)
	return key
}

func (c AuthConfig) AccessDuration() time.Duration {
	return time.Duration(c.AccessExpire)
}

func (c AuthConfig) RefreshDuration() time.Duration {
	return time.Duration(c.RefreshExpire)
}

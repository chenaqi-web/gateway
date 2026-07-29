package config

import (
	"fmt"
	"time"
)

const defaultRefreshExpire = 7 * 24 * time.Hour

type AuthConfig struct {
	RefreshExpire string `yaml:"refresh_expire"`
	CookieSecure  bool   `yaml:"cookie_secure"`
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

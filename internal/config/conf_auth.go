package config

type AuthConfig struct {
	JWTSecret     string `yaml:"jwt_secret"`
	AccessExpire  int    `yaml:"access_expire"`
	RefreshExpire int    `yaml:"refresh_expire"`
	CookieSecure  bool   `yaml:"cookie_secure"`
}

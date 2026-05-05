package config

type JWTConfig struct {
	Secret        string
	ExpirationHrs int
}

func NewJWTConfig(cfg *Config) *JWTConfig {
	return &JWTConfig{
		Secret:        cfg.JWTSecret,
		ExpirationHrs: 24,
	}
}

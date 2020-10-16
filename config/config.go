package config

import "github.com/spf13/viper"

// Config holds app-wide config
type Config struct {
	Prod               bool
	Host               string
	Port               string
	CorsAllowedOrigins []string `mapstructure:"cors_allowed_origins"`
}

// Load loads and creates the config object
func Load() (*Config, error) {
	v := viper.New()
	v.SetConfigType("env")
	v.SetConfigName(".env")
	v.AddConfigPath(".")

	if err := v.ReadInConfig(); err != nil {
		return nil, err
	}

	var c Config
	if err := v.Unmarshal(&c); err != nil {
		return nil, err
	}

	return &c, nil
}

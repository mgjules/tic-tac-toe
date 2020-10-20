package config

import "github.com/spf13/viper"

// Config holds app-wide config
type Config struct {
	Prod                          bool
	Host                          string
	Port                          string
	CorsAllowedOrigins            []string `mapstructure:"cors_allowed_origins"`
	FirebaseDBURL                 string   `mapstructure:"firebase_db_url"`
	FirebaseServiceAccountKeyPath string   `mapstructure:"firebase_service_account_key_path"`
}

// Load loads and creates the config object using config path
func Load(path string) (*Config, error) {
	v := viper.New()
	v.SetConfigType("env")
	v.SetConfigName(path)
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

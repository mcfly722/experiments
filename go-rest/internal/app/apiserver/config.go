package apiserver

import "github.com/mcfly722/experiments/go-rest/internal/app/store"

// Config ...
type Config struct {
	BindAddr      string `toml:"bind_addr"`
	LogLevel      string `toml:"log_level"`
	JWTPrivateKey string `toml:"jwt_private_key"`
	Store         *store.Config
}

//NewConfig ...
func NewConfig() *Config {
	return &Config{
		BindAddr:      ":8080",
		LogLevel:      "debug",
		Store:         store.NewConfig(),
		JWTPrivateKey: "",
	}

}

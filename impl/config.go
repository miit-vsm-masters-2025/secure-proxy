package impl

import (
	"fmt"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type AppConfig struct {
	CookieName string     `yaml:"cookieName"`
	Valkey     Valkey     `yaml:"valkey"`
	AuthDomain string     `yaml:"authDomain"`
	Upstreams  []Upstream `yaml:"upstreams"`
	Users      []User     `yaml:"users"`
}

type Valkey struct {
	Address    string        `yaml:"address"`
	SessionTtl time.Duration `yaml:"sessionTtl"`
}

type Upstream struct {
	Host        string `yaml:"host"`
	Destination string `yaml:"destination"`
}

type User struct {
	Username         string   `yaml:"username"`
	TotpSecret       string   `yaml:"totpSecret"`
	AvailableDomains []string `yaml:"availableDomains"`
}

func createConfig() *AppConfig {
	duration, err := time.ParseDuration("10m")
	if err != nil {
		panic(err)
	}
	config := AppConfig{
		CookieName: "SECURE_PROXY_SESSION",
		Valkey: Valkey{
			Address:    "127.0.0.1:6379",
			SessionTtl: duration,
		},
	}
	configFile, configPathOverridden := os.LookupEnv("APP_CONFIG")
	if !configPathOverridden {
		configFile = "config.yaml"
	}
	err = cleanenv.ReadConfig(configFile, &config)

	if err != nil {
		panic(fmt.Errorf("failed to load config: %w", err))
	}
	return &config
}

var config = createConfig()

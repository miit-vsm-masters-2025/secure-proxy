package impl

import (
	"fmt"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type AppConfig struct {
	AuthDomain string     `yaml:"authDomain"`
	Upstreams  []Upstream `yaml:"upstreams"`
}

type Upstream struct {
	Host        string `yaml:"host"`
	Destination string `yaml:"destination"`
}

func createConfig() *AppConfig {
	config := AppConfig{}
	configFile, configPathOverridden := os.LookupEnv("APP_CONFIG")
	if !configPathOverridden {
		configFile = "config.yaml"
	}
	err := cleanenv.ReadConfig(configFile, &config)

	if err != nil {
		panic(fmt.Errorf("failed to load config: %w", err))
	}
	return &config
}

var config = createConfig()

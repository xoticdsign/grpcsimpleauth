package config

import (
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env             string        `yaml:"env"`
	StoragePath     string        `yaml:"storage_path"`
	MigrationsPath  string        `yaml:"migrations_path"`
	MigrationsTable string        `yaml:"migrations_table"`
	TokenTTL        time.Duration `yaml:"token_ttl"`
	GRPC            GRPCConfig    `yaml:"grpc"`
}

type GRPCConfig struct {
	Port    int           `yaml:"port"`
	Timeout time.Duration `yaml:"timeout"`
}

func MustLoad() *Config {
	var cfg Config

	env := os.Getenv("env")
	defer os.Unsetenv("env")

	if env != "" {
		switch {
		case env == "dev":
			// TODO

		case env == "prod":
			// TODO
		}
	} else {
		err := cleanenv.ReadConfig("./sso/config/local.yaml", &cfg)
		if err != nil {
			panic("couldn't read config: " + err.Error())
		}
	}

	return &cfg
}

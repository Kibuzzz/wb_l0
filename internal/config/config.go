package config

import (
	"log"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Database struct {
		DSN string `yaml:"dsn"`
	} `yaml:"database"`
	Server struct {
		Host string `yaml:"host"  env-default:"0.0.0.0"`
		Port string `yaml:"port"  env-default:"8081"`
	} `yaml:"server"`
	NATS struct {
		ClusterID string `yaml:"cluster_id"`
		ClientID  string `yaml:"client_id"`
	} `yaml:"nats"`
}

func Get() Config {
	var cfg Config
	if err := cleanenv.ReadConfig("config.yaml", &cfg); err != nil {
		log.Fatalln(err)
	}
	return cfg
}

package config

import (
	"io"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server        server        `yaml:"server"`
	Microservices microservices `yaml:"microservices"`
}

type server struct {
	Addr string `yaml:"addr"`
	Port string `yaml:"port"`
}

type microservices struct {
	Static server `yaml:"static"`
}

func ReadConfig(path string) (*Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	cfg := &Config{}
	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

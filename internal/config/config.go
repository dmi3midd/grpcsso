package config

import (
	"time"
)

type Config struct {
	Server   ServerConfig
	Postgres PostgresConfig
}

type ServerConfig struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

type PostgresConfig struct {
	Name         string        `yaml:"name"`
	Host         string        `yaml:"host"`
	Port         int           `yaml:"port"`
	User         string        `yaml:"user"`
	Password     string        `yaml:"password"`
	SSLMode      string        `yaml:"sslmode"`
	MaxOpenConns int           `yaml:"maxOpenConns"`
	MaxIdleConns int           `yaml:"maxIdleConns"`
	MaxIdleTime  time.Duration `yaml:"maxIdleTime"`
}

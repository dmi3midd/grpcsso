package config

import (
	"crypto/rsa"
	"time"
)

type Config struct {
	Server   ServerConfig
	Postgres PostgresConfig
	Redis    RedisConfig
	PEM      PEMConfig
	Keys     KeysPair
}

type ServerConfig struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

type PostgresConfig struct {
	DBName       string        `yaml:"dbname"`
	Host         string        `yaml:"host"`
	Port         int           `yaml:"port"`
	User         string        `yaml:"user"`
	Password     string        `yaml:"password"`
	SSLMode      string        `yaml:"sslmode"`
	MaxOpenConns int           `yaml:"maxOpenConns"`
	MaxIdleConns int           `yaml:"maxIdleConns"`
	MaxIdleTime  time.Duration `yaml:"maxIdleTime"`
}

type RedisConfig struct {
	URI string `yaml:"uri"`
}

type PEMConfig struct {
	PrivPath string `yaml:"priv_path"`
	PubPath  string `yaml:"pub_path"`
}

type KeysPair struct {
	PrivateKey *rsa.PrivateKey
	PublicKey  *rsa.PublicKey
}

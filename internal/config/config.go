package config

import (
	"crypto/rsa"
	"time"
)

type Config struct {
	Server   ServerConfig   `yaml:"server"`
	Postgres PostgresConfig `yaml:"postgres"`
	JWT      JWTConfig      `yaml:"jwt"`
	PEM      PEMConfig      `yaml:"pem"`
	Keys     KeysPair
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

type JWTConfig struct {
	AccessTokenTTL  time.Duration `yaml:"accessTokenTTL"`
	RefreshTokenTTL time.Duration `yaml:"refreshTokenTTL"`
	Audience        string        `yaml:"audience"`
}

type PEMConfig struct {
	PrivPath string `yaml:"privPath"`
	PubPath  string `yaml:"pubPath"`
}

type KeysPair struct {
	PrivateKey *rsa.PrivateKey
	PublicKey  *rsa.PublicKey
}

package config

import (
	"fmt"
	"os"

	"github.com/golang-jwt/jwt/v5"
	"gopkg.in/yaml.v3"
)

func LoadConfig() (*Config, error) {
	data, err := os.ReadFile("./config.yaml")
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	cfg := &Config{}
	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config file: %w", err)
	}

	keys, err := LoadKeys(cfg.PEM.PrivPath, cfg.PEM.PubPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load keys: %w", err)
	}
	cfg.Keys = *keys

	return cfg, nil
}

func LoadKeys(privPath, pubPath string) (*KeysPair, error) {
	privKeyData, err := os.ReadFile(privPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read private key from %s: %w", privPath, err)
	}

	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(privKeyData)
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %w", err)
	}

	pubKeyData, err := os.ReadFile(pubPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read public key from %s: %w", pubPath, err)
	}

	publicKey, err := jwt.ParseRSAPublicKeyFromPEM(pubKeyData)
	if err != nil {
		return nil, fmt.Errorf("failed to parse public key: %w", err)
	}

	return &KeysPair{
		PrivateKey: privateKey,
		PublicKey:  publicKey,
	}, nil
}

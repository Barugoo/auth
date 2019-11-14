package config

import (
	"github.com/kelseyhightower/envconfig"
)

type ServiceConfig struct {
	Issuer2FA   string `envconfig:"issuer_2fa"`
	AppSecret   string `envconfig:"app_secret"`
	AddressGRPC string `envconfig:"address_grpc"`
	MongoURI    string `envconfig:"mongo_uri"`
}

func NewConfig() (*ServiceConfig, error) {
	cfg := &ServiceConfig{}
	err := envconfig.Process("AUTH", cfg)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}

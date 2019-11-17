package config

import (
	"github.com/kelseyhightower/envconfig"
)

type ServiceConfig struct {
	ServiceName string `envconfig:"service_name"`
	Issuer2FA   string `envconfig:"issuer_2fa"`
	AppSecret   string `envconfig:"app_secret"`
	GRPCAddr    string `envconfig:"grpc_addr"`
	MongoAddr   string `envconfig:"mongo_addr"`
	RedisAddr   string `envconfig:"mongo_addr"`
}

func NewConfig() (*ServiceConfig, error) {
	cfg := &ServiceConfig{}
	err := envconfig.Process("AUTH", cfg)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}

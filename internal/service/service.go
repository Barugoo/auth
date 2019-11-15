package service

import (
	opentracing "github.com/opentracing/opentracing-go"
)

type AuthService interface {
	SetKV(key, value string) (bool, error)
	GetKV(key string) (string, error)
	GetTracer() opentracing.Tracer
}

type authService struct {
	tracer opentracing.Tracer
	kv     map[string]string
}

func (a *authService) SetKV(k, v string) (bool, error) {
	a.kv[k] = v
	return true, nil
}

func (a *authService) GetKV(k string) (string, error) {
	return a.kv[k], nil
}

func (a *authService) GetTracer() opentracing.Tracer {
	return a.tracer
}

func NewAuthService(tracer opentracing.Tracer) AuthService {
	return &authService{
		tracer: tracer,
		kv:     make(map[string]string),
	}
}

package service

import (
	opentracing "github.com/opentracing/opentracing-go"
)

type AuthService interface {
	SetKV(key, value string) (bool, error)
	GetKV(key string) (string, error)
	StartSpan(ctx context.Context, name string) opentracing.Span
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

func (a *authService) StartSpan(ctx context.Context, name string) opentracing.Span {
	span := opentracing.SpanFromContext(ctx)
	if span != nil {
		span = a.tracer.StartSpan(name, opentracing.ChildOf(span.Context()))
	} else {
		span = a.tracer.StartSpan(name)
	}

	return span
}

func NewAuthService(tracer opentracing.Tracer) AuthService {
	return &authService{
		tracer: tracer,
		kv:     make(map[string]string),
	}
}

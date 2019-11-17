package service

import (
	"context"

	"github.com/go-redis/redis/v7"
	opentracing "github.com/opentracing/opentracing-go"
)

type AuthService interface {
	SetKV(ctx context.Context, key, value string) (bool, error)
	GetKV(ctx context.Context, key string) (string, error)

	StartSpan(ctx context.Context, name string) opentracing.Span
	ContextWithSpan(ctx context.Context, span opentracing.Span) context.Context
}

type authService struct {
	redisClient *redis.Client
	tracer      opentracing.Tracer
	kv          map[string]string
}

func (a *authService) SetKV(ctx context.Context, key, value string) (bool, error) {
	span := a.StartSpan(ctx, "SetKV")
	defer span.Finish()

	return a.setKV(key, value)
}

func (a *authService) setKV(key, value string) (bool, error) {
	err := a.redisClient.Set(key, value, 0).Err()
	if err != nil {
		return false, err
	}
	val, err := a.redisClient.Get(key).Result()
	if err != nil {
		return false, err
	}
	if val != value {
		return false, err
	}
	return true, nil
}

func (a *authService) GetKV(ctx context.Context, key string) (string, error) {
	span := a.StartSpan(ctx, "GetKV")
	defer span.Finish()

	return a.getKV(key)
}

func (a *authService) getKV(key string) (string, error) {
	val, err := a.redisClient.Get(key).Result()
	if err != nil {
		return "", err
	}
	return val, nil
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

func (a *authService) ContextWithSpan(ctx context.Context, span opentracing.Span) context.Context {
	return opentracing.ContextWithSpan(ctx, span)
}

func NewAuthService(redisClient *redis.Client, tracer opentracing.Tracer) AuthService {
	return &authService{
		redisClient: redisClient,
		tracer:      tracer,
		kv:          make(map[string]string),
	}
}

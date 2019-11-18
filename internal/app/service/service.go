package service

import (
	"context"

	"github.com/go-redis/redis/v7"
	opentracing "github.com/opentracing/opentracing-go"

	"github.com/barugoo/oscillo-auth/internal/app/errors"
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
}

func (a *authService) SetKV(ctx context.Context, key, value string) (bool, error) {
	methodName := "SetKV/redis"

	span := a.StartSpan(ctx, methodName)
	defer span.Finish()

	ok, err := a.setKV(key, value)
	if err != nil {
		err = a.wrapError(err, methodName)
	}
	return ok, err
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
	methodName := "GetKV/redis"

	span := a.StartSpan(ctx, methodName)
	defer span.Finish()

	value, err := a.getKV(key)
	if err != nil {
		err = a.wrapError(err, methodName)
	}
	return value, err
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

func (a *authService) wrapError(err error, method string) error {
	return errors.ServiceError{
		Method: method,
		Err:    err,
	}
}

func NewAuthService(redisClient *redis.Client, tracer opentracing.Tracer) AuthService {
	return &authService{
		redisClient: redisClient,
		tracer:      tracer,
	}
}

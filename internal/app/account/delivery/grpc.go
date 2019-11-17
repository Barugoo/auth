package delivery

import (
	"context"
	"errors"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/barugoo/oscillo-auth/api/grpc"

	models "github.com/barugoo/oscillo-auth/internal/app/account"

	"github.com/barugoo/oscillo-auth/internal/app/account/usecase"
	errs "github.com/barugoo/oscillo-auth/internal/app/errors"
	"github.com/barugoo/oscillo-auth/internal/app/service"
)

const (
	deliveryMethodTemplate = "%s/delivery"
)

type authGRPCServer struct {
	service     service.AuthService
	accountCase usecase.AccountUsecase
}

func NewAuthGRPCServer(service service.AuthService, accountUsecase usecase.AccountUsecase) pb.AuthServer {
	return &authGRPCServer{
		service:     service,
		accountCase: accountUsecase,
	}
}

func (auth *authGRPCServer) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	methodName, err := auth.getMethodFromContext(ctx)
	if err != nil {
		return nil, err
	}

	span := auth.service.StartSpan(ctx, methodName)
	defer span.Finish()

	spanCtx := auth.service.ContextWithSpan(context.Background(), span)
	methodCtx := auth.contextWithMethod(spanCtx, methodName)

	resp, err := auth.register(methodCtx, req)
	if err != nil {
		err = auth.grpcError(auth.wrapError(err, req.Email))
	}
	return resp, err
}

func (auth *authGRPCServer) register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	r := &models.Credentials{
		Email:    req.Email,
		Password: req.Password,
	}
	ok, err := auth.accountCase.RegisterWithCredentials(ctx, r)
	if err != nil {
		return nil, err
	}
	return &pb.RegisterResponse{
		Ok: ok,
	}, err
}

func (auth *authGRPCServer) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	methodName, err := auth.getMethodFromContext(ctx)
	if err != nil {
		return nil, err
	}

	span := auth.service.StartSpan(ctx, methodName)
	defer span.Finish()

	spanCtx := auth.service.ContextWithSpan(context.Background(), span)
	methodCtx := auth.contextWithMethod(spanCtx, methodName)

	resp, err := auth.login(methodCtx, req)
	if err != nil {
		err = auth.grpcError(auth.wrapError(err, req.Email))
	}
	return resp, err
}

func (auth *authGRPCServer) login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	r := &models.Credentials{
		Email:    req.Email,
		Password: req.Password,
	}
	token, err := auth.accountCase.AuthByCredentials(ctx, r)
	if err != nil {
		return nil, err
	}
	return &pb.LoginResponse{
		Token: token,
	}, err
}

func (auth *authGRPCServer) UpdateCredentials(ctx context.Context, req *pb.UpdateCredentialsRequest) (*pb.UpdateCredentialsResponse, error) {
	methodName, err := auth.getMethodFromContext(ctx)
	if err != nil {
		return nil, err
	}

	span := auth.service.StartSpan(ctx, methodName)
	defer span.Finish()

	spanCtx := auth.service.ContextWithSpan(context.Background(), span)
	methodCtx := auth.contextWithMethod(spanCtx, methodName)

	resp, err := auth.updateCredentials(methodCtx, req)
	if err != nil {
		err = auth.grpcError(auth.wrapError(err, req.Email))
	}
	return resp, err
}

func (auth *authGRPCServer) updateCredentials(ctx context.Context, req *pb.UpdateCredentialsRequest) (*pb.UpdateCredentialsResponse, error) {
	r := &models.Credentials{
		Email:    req.Email,
		Password: req.Password,
	}
	ok, err := auth.accountCase.UpdateCredentials(ctx, r)
	if err != nil {
		return nil, err
	}
	return &pb.UpdateCredentialsResponse{
		Ok: ok,
	}, err
}

func (auth *authGRPCServer) ActivateAccount(ctx context.Context, req *pb.ActivateAccountRequest) (*pb.ActivateAccountResponse, error) {
	methodName, err := auth.getMethodFromContext(ctx)
	if err != nil {
		return nil, err
	}

	span := auth.service.StartSpan(ctx, methodName)
	defer span.Finish()

	spanCtx := auth.service.ContextWithSpan(context.Background(), span)
	methodCtx := auth.contextWithMethod(spanCtx, methodName)

	resp, err := auth.activateAccount(methodCtx, req)
	if err != nil {
		err = auth.grpcError(auth.wrapError(err, req.Email))
	}
	return resp, err
}

func (auth *authGRPCServer) activateAccount(ctx context.Context, req *pb.ActivateAccountRequest) (*pb.ActivateAccountResponse, error) {
	ok, err := auth.accountCase.ActivateAccount(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	return &pb.ActivateAccountResponse{
		Ok: ok,
	}, err
}

func (auth *authGRPCServer) Generate2FA(ctx context.Context, req *pb.Generate2FARequest) (*pb.Generate2FAResponse, error) {
	methodName, err := auth.getMethodFromContext(ctx)
	if err != nil {
		return nil, err
	}

	span := auth.service.StartSpan(ctx, methodName)
	defer span.Finish()

	spanCtx := auth.service.ContextWithSpan(context.Background(), span)
	methodCtx := auth.contextWithMethod(spanCtx, methodName)

	resp, err := auth.generate2FA(methodCtx, req)
	if err != nil {
		err = auth.grpcError(auth.wrapError(err, req.Email))
	}
	return resp, err
}

func (auth *authGRPCServer) generate2FA(ctx context.Context, req *pb.Generate2FARequest) (*pb.Generate2FAResponse, error) {
	img, err := auth.accountCase.Generate2FA(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	return &pb.Generate2FAResponse{
		QrImage: img,
	}, err
}

func (auth *authGRPCServer) Setup2FA(ctx context.Context, req *pb.Setup2FARequest) (*pb.Setup2FAResponse, error) {
	methodName, err := auth.getMethodFromContext(ctx)
	if err != nil {
		return nil, err
	}

	span := auth.service.StartSpan(ctx, methodName)
	defer span.Finish()

	spanCtx := auth.service.ContextWithSpan(context.Background(), span)
	methodCtx := auth.contextWithMethod(spanCtx, methodName)

	resp, err := auth.setup2FA(methodCtx, req)
	if err != nil {
		err = auth.grpcError(auth.wrapError(err, req.Email))
	}
	return resp, err
}

func (auth *authGRPCServer) setup2FA(ctx context.Context, req *pb.Setup2FARequest) (*pb.Setup2FAResponse, error) {
	ok, err := auth.accountCase.Setup2FA(ctx, req.Email, req.Code)
	if err != nil {
		return nil, err
	}
	return &pb.Setup2FAResponse{
		Ok: ok,
	}, err
}

func (auth *authGRPCServer) Disable2FA(ctx context.Context, req *pb.Disable2FARequest) (*pb.Disable2FAResponse, error) {
	methodName, err := auth.getMethodFromContext(ctx)
	if err != nil {
		return nil, err
	}

	span := auth.service.StartSpan(ctx, methodName)
	defer span.Finish()

	spanCtx := auth.service.ContextWithSpan(context.Background(), span)
	methodCtx := auth.contextWithMethod(spanCtx, methodName)

	resp, err := auth.disable2FA(methodCtx, req)
	if err != nil {
		err = auth.grpcError(auth.wrapError(err, req.Email))
	}
	return resp, err
}

func (auth *authGRPCServer) disable2FA(ctx context.Context, req *pb.Disable2FARequest) (*pb.Disable2FAResponse, error) {
	ok, err := auth.accountCase.Remove2FA(ctx, req.Email, req.Code)
	if err != nil {
		return nil, err
	}
	return &pb.Disable2FAResponse{
		Ok: ok,
	}, err
}

func (auth *authGRPCServer) Verify2FA(ctx context.Context, req *pb.Verify2FARequest) (*pb.Verify2FAResponse, error) {
	methodName, err := auth.getMethodFromContext(ctx)
	if err != nil {
		return nil, err
	}

	span := auth.service.StartSpan(ctx, methodName)
	defer span.Finish()

	spanCtx := auth.service.ContextWithSpan(context.Background(), span)
	methodCtx := auth.contextWithMethod(spanCtx, methodName)

	resp, err := auth.verify2FA(methodCtx, req)
	if err != nil {
		err = auth.grpcError(auth.wrapError(err, req.Email))
	}
	return resp, err
}

func (auth *authGRPCServer) verify2FA(ctx context.Context, req *pb.Verify2FARequest) (*pb.Verify2FAResponse, error) {
	ok, err := auth.accountCase.Verify2FA(ctx, req.Email, req.Code)
	if err != nil {
		return nil, err
	}
	return &pb.Verify2FAResponse{
		Ok: ok,
	}, err
}

func (auth *authGRPCServer) contextWithMethod(ctx context.Context, method string) context.Context {
	return context.WithValue(ctx, "method", method)
}

func (auth *authGRPCServer) getMethodFromContext(ctx context.Context) (string, error) {
	methodName, ok := grpc.Method(ctx)
	if !ok {
		return "", errs.ErrBrokenContext
	}
	return fmt.Sprintf(deliveryMethodTemplate, methodName), nil
}

func (auth *authGRPCServer) grpcError(err error) error {
	return status.Error(auth.mapStatusCode(err), err.Error())
}

func (auth *authGRPCServer) wrapError(err error, email string) error {
	return &errs.DeliveryError{
		Email: email,
		Err:   err,
	}
}

func (auth *authGRPCServer) mapStatusCode(err error) codes.Code {
	var repErr *errs.RepositoryError
	if errors.As(err, &repErr) {
		switch repErr.Err {
		case errs.ErrNotFound:
			return codes.NotFound
		default:
			return codes.Internal
		}
	}

	var serviceErr *errs.ServiceError
	if errors.As(err, &serviceErr) {
		switch serviceErr.Err {
		default:
			return codes.Internal
		}
	}

	var caseErr *errs.UsecaseError
	if errors.As(err, &caseErr) {
		switch caseErr.Err {
		case errs.ErrWrongPassword, errs.ErrInvalid2FACode, errs.ErrInactiveAccount:
			return codes.Unauthenticated
		case errs.Err2FADisabled:
			return codes.InvalidArgument
		case errs.ErrUnableToStoreKey:
			fallthrough
		default:
			return codes.Internal
		}
	}
	return codes.Unknown
}

package delivery

import (
	"context"

	pb "github.com/barugoo/oscillo-auth/api/grpc"

	"github.com/barugoo/oscillo-auth/internal/app/models"

	"github.com/barugoo/oscillo-auth/internal/app/account/usecase"
	"github.com/barugoo/oscillo-auth/internal/app/service"
)

type authGRPCServer struct {
	service     service.AuthService
	accountCase usecase.AccountUsecase
}

func NewAuthGRPCServer(service service.AuthService, accountUsecase usecase.AccountUsecase) pb.AuthServer {
	return &authGRPCServer{
		service:        service,
		accountUsecase: accountUsecase,
	}
}

func (auth *authGRPCServer) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterReply, error) {
	span := auth.service.StartSpan(ctx, "Register")
	defer span.Finish()

	r := &models.Credentials{
		Email:    req.Email,
		Password: req.Password,
	}
	ok, err := auth.accountCase.RegisterWithCredentials(r)
	if err != nil {
		return nil, err
	}
	return &pb.RegisterReply{
		Ok: ok,
	}, err
}

func (auth *authGRPCServer) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginReply, error) {
	span := auth.service.StartSpan(ctx, "Login")
	defer span.Finish()

	r := &models.Credentials{
		Email:    req.Email,
		Password: req.Password,
	}
	token, err := auth.accountCase.AuthByCredentials(r)
	if err != nil {
		return nil, err
	}
	return &pb.LoginReply{
		Token: token,
	}, err
}

func (auth *authGRPCServer) UpdateCredentials(ctx context.Context, req *pb.UpdateCredentialsRequest) (*pb.UpdateCredentialsReply, error) {
	span := auth.service.StartSpan(ctx, "UpdateCredentials")
	defer span.Finish()

	r := &models.Credentials{
		Email:    req.Email,
		Password: req.Password,
	}
	ok, err := auth.accountCase.UpdateCredentials(r)
	if err != nil {
		return nil, err
	}
	return &pb.UpdateCredentialsReply{
		Ok: ok,
	}, err
}

func (auth *authGRPCServer) ActivateAccount(ctx context.Context, req *pb.ActivateAccountRequest) (*pb.ActivateAccountReply, error) {
	span := auth.service.StartSpan(ctx, "ActivateAccount")
	defer span.Finish()

	ok, err := auth.accountCase.ActivateAccount(req.Email)
	if err != nil {
		return nil, err
	}
	return &pb.ActivateAccountReply{
		Ok: ok,
	}, err
}

func (auth *authGRPCServer) Generate2FA(ctx context.Context, req *pb.Generate2FARequest) (*pb.Generate2FAReply, error) {
	span := auth.service.StartSpan(ctx, "Generate2FA")
	defer span.Finish()

	img, err := auth.accountCase.Generate2FA(req.Email)
	if err != nil {
		return nil, err
	}
	return &pb.Generate2FAReply{
		QrImage: img,
	}, err
}

func (auth *authGRPCServer) Setup2FA(ctx context.Context, req *pb.Setup2FARequest) (*pb.Setup2FAReply, error) {
	span := auth.service.StartSpan(ctx, "Setup2FA")
	defer span.Finish()

	ok, err := auth.accountCase.Setup2FA(req.Email, req.Code)
	if err != nil {
		return nil, err
	}
	return &pb.Setup2FAReply{
		Ok: ok,
	}, err
}

func (auth *authGRPCServer) Disable2FA(ctx context.Context, req *pb.Disable2FARequest) (*pb.Disable2FAReply, error) {
	span := auth.service.StartSpan(ctx, "Disable2FA")
	defer span.Finish()

	ok, err := auth.accountCase.Remove2FA(req.Email, req.Code)
	if err != nil {
		return nil, err
	}
	return &pb.Disable2FAReply{
		Ok: ok,
	}, err
}

func (auth *authGRPCServer) Verify2FA(ctx context.Context, req *pb.Verify2FARequest) (*pb.Verify2FAReply, error) {
	span := auth.service.StartSpan(ctx, "Verify2FA")
	defer span.Finish()

	ok, err := auth.accountCase.Verify2FA(req.Email, req.Code)
	if err != nil {
		return nil, err
	}
	return &pb.Verify2FAReply{
		Ok: ok,
	}, err
}

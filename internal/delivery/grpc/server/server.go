package server

import (
	"context"

	pb "github.com/barugoo/oscillo-auth/api/grpc"

	"github.com/barugoo/oscillo-auth/internal/models"
	"github.com/barugoo/oscillo-auth/internal/usecase"
)

type authGRPCServer struct {
	uc usecase.AccountUsecase
}

func (auth *authGRPCServer) Register(ctx *context.Context, req *pb.RegisterRequest) (*pb.RegisterReply, error) {
	r := &models.Credentials{
		Email:    req.Email,
		Password: req.Password,
	}
	ok, err := auth.uc.RegisterWithCredentials(r)
	if err != nil {
		return nil, err
	}
	return &pb.RegisterReply{
		Ok: ok,
	}, err
}

func (auth *authGRPCServer) Login(ctx *context.Context, req *pb.LoginRequest) (*pb.LoginReply, error) {
	r := &models.Credentials{
		Email:    req.Email,
		Password: req.Password,
	}
	token, err := auth.uc.AuthByCredentials(r)
	if err != nil {
		return nil, err
	}
	return &pb.LoginReply{
		Token: token,
	}, err
}

func (auth *authGRPCServer) UpdateCredentials(ctx *context.Context, req *pb.LoginRequest) (*pb.UpdateCredentialsReply, error) {
	r := &models.Credentials{
		Email:    req.Email,
		Password: req.Password,
	}
	ok, err := auth.uc.UpdateCredentials(r)
	if err != nil {
		return nil, err
	}
	return &pb.UpdateCredentialsReply{
		Ok: ok,
	}, err
}

func (auth *authGRPCServer) ActivateAccount(ctx *context.Context, req *pb.ActivateAccountRequest) (*pb.ActivateAccountReply, error) {
	ok, err := auth.uc.ActivateAccount(req.Email)
	if err != nil {
		return nil, err
	}
	return &pb.ActivateAccountReply{
		Ok: ok,
	}, err
}

func (auth *authGRPCServer) Generate2FA(ctx *context.Context, req *pb.Generate2FARequest) (*pb.Generate2FAReply, error) {
	img, err := auth.uc.Generate2FA(req.Email)
	if err != nil {
		return nil, err
	}
	return &pb.Generate2FAReply{
		QrImage: img,
	}, err
}

func (auth *authGRPCServer) Setup2FA(ctx *context.Context, req *pb.Setup2FARequest) (*pb.Setup2FAReply, error) {
	ok, err := auth.uc.Setup2FA(req.Email, req.Code)
	if err != nil {
		return nil, err
	}
	return &pb.Setup2FAReply{
		Ok: ok,
	}, err
}

func (auth *authGRPCServer) Disable2FA(ctx *context.Context, req *pb.Disable2FARequest) (*pb.Disable2FAReply, error) {
	ok, err := auth.uc.Remove2FA(req.Email, req.Code)
	if err != nil {
		return nil, err
	}
	return &pb.Disable2FAReply{
		Ok: ok,
	}, err
}

func (auth *authGRPCServer) Verify2FA(ctx *context.Context, req *pb.Verify2FARequest) (*pb.Verify2FAReply, error) {
	ok, err := auth.uc.Verify2FA(req.Email, req.Code)
	if err != nil {
		return nil, err
	}
	return &pb.Verify2FAReply{
		Ok: ok,
	}, err
}

func NewAuthGRPCServer(accountUsecase usecase.AccountUsecase) pb.AuthServer {
	return &authGRPCServer{accountUsecase}
}

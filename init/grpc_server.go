package init

import (
	"google.golang.org/grpc"

	api "github.com/barugoo/oscillo-auth/api/grpc"
	auth "github.com/barugoo/oscillo-auth/internal/delivery/grpc/server"
	"github.com/barugoo/oscillo-auth/internal/usecase"
)

func NewGRPCServer(accUC usecase.AccountUsecase) *grpc.Server {
	s := auth.NewAuthGRPCServer(accUC)

	grpcServer := grpc.NewServer()
	api.RegisterAuthServer(grpcServer, s)

	return grpcServer
}

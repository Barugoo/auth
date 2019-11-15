package app

import (
	"io"
	"net"

	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc"

	api "github.com/barugoo/oscillo-auth/api/grpc"
	"github.com/barugoo/oscillo-auth/config"
	"github.com/barugoo/oscillo-auth/init/tracer"

	"github.com/barugoo/oscillo-auth/internal/app/service"

	accountDelivery "github.com/barugoo/oscillo-auth/internal/app/account/delivery"
	accountRepository "github.com/barugoo/oscillo-auth/internal/app/account/repository"
	accountUsecase "github.com/barugoo/oscillo-auth/internal/app/account/usecase"
)

type App interface {
	Run() error
	Shutdown()
}

type authApp struct {
	server       *grpc.Server
	config       *config.ServiceConfig
	tracerCloser io.Closer
}

func NewAuthApp(config *config.ServiceConfig, db *mongo.Database) (App, error) {

	tracer, closer, err := tracer.NewTracer(config)
	if err != nil {
		return nil, err
	}

	service := service.NewAuthService(tracer)

	accountRep := accountRepository.NewAccountRepository(db.Collection("account"))
	accountCase := accountUsecase.NewAccountUsecase(config, service, accountRep)
	accountDelv := accountDelivery.NewAuthGRPCServer(accountCase)

	grpcServ := grpc.NewServer()
	api.RegisterAuthServer(grpcServ, accountDelv)

	return &authApp{
		tracerCloser: closer,
		server:       grpcServ,
		config:       config,
	}, nil
}

func (app *authApp) Run() error {
	lis, err := net.Listen("tcp", app.config.AddressGRPC)
	if err != nil {
		return err
	}
	err = app.server.Serve(lis)
	return err
}

func (app *authApp) Shutdown() {
	app.tracerCloser.Close()
}

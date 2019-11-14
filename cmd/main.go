package main

import (
	"log"
	"net"

	"github.com/barugoo/oscillo-auth/config"
	"github.com/barugoo/oscillo-auth/init"
	"github.com/barugoo/oscillo-auth/internal/repository"
	"github.com/barugoo/oscillo-auth/internal/service"
	"github.com/barugoo/oscillo-auth/internal/usecase"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatal(err)
	}

	mgoClient, err := init.NewMongoClient(cfg.MongoURI)
	if err != nil {
		log.Fatal(err)
	}

	db := repository.NewAccountRepository(mgoClient)
	service := service.NewAuthService()
	usecase := usecase.NewAccountUsecase(cfg, service, db)

	grpcServer, err := init.NewGRPCServer(usecase)
	if err != nil {
		log.Fatal(err)
	}
	lis, err := net.Listen("tcp", cfg.AddressGRPC)
	if err != nil {
		return nil, err
	}
	grpcServer.Serve(lis)
}

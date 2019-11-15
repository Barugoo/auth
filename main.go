package main

import (
	"log"

	"github.com/barugoo/oscillo-auth/config"
	mongo "github.com/barugoo/oscillo-auth/init/mongo/client"
	"github.com/barugoo/oscillo-auth/internal/app"
)

const (
	authDB = "auth_db"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatal(err)
	}

	mgoClient, err := mongo.NewMongoClient(cfg)
	if err != nil {
		log.Fatal(err)
	}

	authApp, err := app.NewAuthApp(cfg, mgoClient.Database(authDB))
	if err != nil {
		log.Fatal(err)
	}

	authApp.Run()
}

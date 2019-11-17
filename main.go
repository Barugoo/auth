package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/barugoo/oscillo-auth/config"
	"github.com/barugoo/oscillo-auth/init/mongo"
	"github.com/barugoo/oscillo-auth/init/redis"
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

	redisClient, err := redis.NewRedisClient(cfg)
	if err != nil {
		log.Fatal(err)
	}

	authApp, err := app.NewAuthApp(cfg, redisClient, mgoClient.Database(authDB))
	if err != nil {
		log.Fatal(err)
	}

	stop := make(chan os.Signal)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		defer signal.Stop(stop)
		<-stop
		authApp.Shutdown()
	}()

	err = authApp.Run()
	if err != nil {
		log.Fatal(err)
	}
}

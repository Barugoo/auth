package mongo

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/barugoo/oscillo-auth/config"
)

func NewMongoClient(config *config.ServiceConfig) (*mongo.Client, error) {
	client, err := mongo.NewClient(options.Client().ApplyURI(config.MongoAddr))
	if err != nil {
		return nil, err
	}
	err = client.Connect(context.TODO())
	if err != nil {
		return nil, err
	}
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		return nil, err
	}
	return client, nil
}

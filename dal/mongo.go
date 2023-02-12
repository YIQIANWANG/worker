package dal

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"worker/conf"
)

type MongoClient struct {
	client *mongo.Client
}

func NewMongoClient() (*MongoClient, error) {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(getURI()))
	if err != nil {
		return nil, err
	}

	return &MongoClient{client: client}, nil
}

func (mc *MongoClient) GetDataBase() *mongo.Database {
	return mc.client.Database(conf.DATABASE)
}

/*
func getURI() string {
	protocol := conf.PROTOCOL
	username := conf.USERNAME
	password := conf.PASSWORD
	address := conf.ADDRESS
	authentication := conf.AUTHENTICATION
	return fmt.Sprintf("%s://%s:%s@%s/%s", protocol, username, password, address, authentication)
}
*/

func getURI() string {
	return "mongodb://localhost:27017"
}

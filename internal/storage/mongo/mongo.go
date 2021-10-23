package mongo

import (
	"context"
	"github.com/ppal31/grpc-lab/internal/storage/mongo/codecs"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsoncodec"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"time"
)


func InitDb(uri, dbName string) (*mongo.Database, error) {
	var client *mongo.Client
	var err error
	if client, err = mongo.NewClient(buildClientOptions(uri)...); err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		return nil, err
	}

	if err = pingMongo(client); err != nil {
		return nil, err
	}
	return client.Database(dbName), nil
}

func buildClientOptions(uri string) []*options.ClientOptions {
	opts := make([]*options.ClientOptions, 0)
	opts = append(opts, options.Client().ApplyURI(uri))
	r := buildRegistry()
	opts = append(opts, options.Client().SetRegistry(r))
	return opts
}

func buildRegistry() *bsoncodec.Registry {
	rb := bson.NewRegistryBuilder()
	for k, v := range codecs.EncoderMap(){
		rb.RegisterTypeEncoder(k, v)
	}
	return rb.Build()
}

func pingMongo(client *mongo.Client) error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	return client.Ping(ctx, readpref.Primary())
}

package mongo

import (
	"context"
	"fmt"
	"github.com/go-slark/slark/logger"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/event"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

var (
	client  *mongo.Client
	mongoDB *mongo.Database
)

type MongoConf struct {
	Url         string `json:"url"`
	DateBase    string `json:"data_base"`
	Timeout     int    `json:"timeout"`
	MaxPoolSize uint64 `json:"max_pool_size"`
	MinPoolSize uint64 `json:"min_pool_size"`
	Monitor     bool   `json:"monitor"`
	logger.Logger
}

func createMongoClient(c *MongoConf) (*mongo.Client, error) {
	opts := options.Client()
	if c.MaxPoolSize != 0 {
		opts.SetMaxPoolSize(c.MaxPoolSize)
	}
	if c.MinPoolSize != 0 {
		opts.SetMinPoolSize(c.MinPoolSize)
	}
	if c.Monitor {
		opts.SetMonitor(&event.CommandMonitor{
			Started: func(ctx context.Context, startedEvent *event.CommandStartedEvent) {
				c.Log(ctx, logger.InfoLevel, map[string]interface{}{}, fmt.Sprintf("%v", startedEvent.Command))
			},
			Succeeded: func(ctx context.Context, succeededEvent *event.CommandSucceededEvent) {
				c.Log(ctx, logger.InfoLevel, map[string]interface{}{}, fmt.Sprintf("%v", succeededEvent.Reply))
			},
			Failed: func(ctx context.Context, failedEvent *event.CommandFailedEvent) {
				c.Log(ctx, logger.InfoLevel, map[string]interface{}{}, fmt.Sprintf("%v", failedEvent.Failure))
			},
		})
	}
	opts.ApplyURI(c.Url)
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(c.Timeout)*time.Second)
	defer cancel()
	cli, err := mongo.Connect(ctx, opts)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	err = cli.Ping(context.TODO(), nil)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return cli, nil
}

func InitMongoDB(c *MongoConf) {
	cli, err := createMongoClient(c)
	if err != nil {
		panic(errors.New(fmt.Sprintf("use %+v create mongo client error %+v", c, err)))
	}
	client = cli
	mongoDB = cli.Database(c.DateBase)
}

func NewMongoCollection(coll string) *mongo.Collection {
	return mongoDB.Collection(coll)
}

func GetMongoDB() *mongo.Database {
	return mongoDB
}

func CloseMongo() error {
	if client == nil {
		return nil
	}
	return client.Disconnect(context.TODO())
}

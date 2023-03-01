package mongodb

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Connection configuration
type Option struct {
	Uri      string // example => mongodb://{user}:{password}@{host}:{port}/{dbname}?connect=direct
	Timeout  time.Duration
	PoolSize uint64
}

// Storage structure
type Storage struct {
	mgo *mongo.Client
}

// Create MongoDB storage instance
func NewStorage(option *Option) (*Storage, error) {
	ctx, cancelFunc := context.WithTimeout(context.Background(), option.Timeout)
	defer cancelFunc()

	clientOpts := options.Client().ApplyURI(option.Uri)
	clientOpts.SetMaxPoolSize(option.PoolSize)

	client, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		return nil, err
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, errors.New("Ping failed")
	}

	return &Storage{mgo: client}, nil
}

// Get Database
func (s *Storage) Db(name string) *mongo.Database {
	return s.mgo.Database(name)
}

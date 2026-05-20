package database

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DB struct {
	Client   *mongo.Client
	Database *mongo.Database
}

func Connect(ctx context.Context, uri, dbName string) (*DB, error) {
	opts := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		return nil, err
	}
	if err := client.Ping(ctx, nil); err != nil {
		return nil, err
	}
	db := client.Database(dbName)
	if err := ensureIndexes(ctx, db); err != nil {
		return nil, err
	}
	return &DB{Client: client, Database: db}, nil
}

func (d *DB) Collection(name string) *mongo.Collection {
	return d.Database.Collection(name)
}

func (d *DB) Disconnect(ctx context.Context) error {
	return d.Client.Disconnect(ctx)
}

func ensureIndexes(ctx context.Context, db *mongo.Database) error {
	indexes := []struct {
		coll string
		keys bson.D
		unique bool
	}{
		{"users", bson.D{{Key: "email", Value: 1}}, true},
		{"firs", bson.D{{Key: "firId", Value: 1}}, true},
		{"firs", bson.D{{Key: "status", Value: 1}}, false},
		{"firs", bson.D{{Key: "createdAt", Value: -1}}, false},
		{"firs", bson.D{{Key: "citizenId", Value: 1}}, false},
		{"audit_logs", bson.D{{Key: "createdAt", Value: -1}}, false},
		{"notifications", bson.D{{Key: "userId", Value: 1}, {Key: "read", Value: 1}}, false},
	}
	for _, idx := range indexes {
		opts := options.Index()
		if idx.unique {
			opts.SetUnique(true)
		}
		_, err := db.Collection(idx.coll).Indexes().CreateOne(ctx, mongo.IndexModel{
			Keys:    idx.keys,
			Options: opts,
		})
		if err != nil {
			// ignore duplicate index errors on restart
			continue
		}
	}
	return nil
}

func DefaultTimeout() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), 10*time.Second)
}

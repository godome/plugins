package mongo

import (
	"context"
	"time"

	"github.com/godome/godome/pkg/adapter"
	"github.com/godome/godome/pkg/config"
	"github.com/godome/godome/pkg/logger"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

const AdapterType adapter.AdapterType = "MongoAdapter"

type MongoAdapter interface {
	adapter.Adapter
	Collection(name string, opts ...*options.CollectionOptions) MongoCollection
	Disconnect()
}
type MongoCollection interface {
	Drop(ctx context.Context) error
	Aggregate(ctx context.Context, pipeline interface{}, opts ...*options.AggregateOptions) (*mongo.Cursor, error)
	Find(ctx context.Context, filter interface{}, opts ...*options.FindOptions) (*mongo.Cursor, error)
	FindOne(ctx context.Context, filter interface{}, opts ...*options.FindOneOptions) *mongo.SingleResult
	BulkWrite(ctx context.Context, models []mongo.WriteModel, opts ...*options.BulkWriteOptions) (*mongo.BulkWriteResult, error)
	CountDocuments(ctx context.Context, filter interface{}, opts ...*options.CountOptions) (int64, error)
	DeleteOne(ctx context.Context, filter interface{}, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error)
	DeleteMany(ctx context.Context, filter interface{}, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error)
	UpdateMany(ctx context.Context, filter interface{}, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error)
	UpdateOne(ctx context.Context, filter interface{}, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error)
	InsertMany(ctx context.Context, documents []interface{}, opts ...*options.InsertManyOptions) (*mongo.InsertManyResult, error)
	InsertOne(ctx context.Context, document interface{}, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error)
	FindOneAndUpdate(ctx context.Context, filter interface{}, update interface{}, opts ...*options.FindOneAndUpdateOptions) *mongo.SingleResult
	FindOneAndDelete(ctx context.Context, filter interface{}, opts ...*options.FindOneAndDeleteOptions) *mongo.SingleResult
}

type mongoAdapter struct {
	adapterType adapter.AdapterType
	config      config.Config
	client      *mongo.Client
	database    *mongo.Database
}

func NewMongoAdapter(uri string, name string, retrywrites bool) MongoAdapter {
	a := new(mongoAdapter)

	if uri == "" {
		logger.Fatal("uri is required")
	}

	if name == "" {
		logger.Fatal("database name is required")
	}

	logger.Info("connecting " + name + " db...")
	connectionURI := uri + "/" + name

	if retrywrites {
		connectionURI = connectionURI + "?retryWrites=true"
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(connectionURI))
	if err != nil {
		logger.Debug(err.Error())
		logger.Fatal("mongo connection error!")
	}

	// Check the connection
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		logger.Fatal(err.Error())
	}

	database := client.Database(name)
	logger.Info(name + " db is connected successfully!")

	a.database = database
	a.client = client

	return a
}

func (d *mongoAdapter) Disconnect() {
	logger.Info("disconnection " + d.database.Name() + " db...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	d.database.Client().Disconnect(ctx)
	logger.Info(d.database.Name() + " is disconnected successfully")
}

func (d *mongoAdapter) Collection(name string, opts ...*options.CollectionOptions) MongoCollection {
	return d.database.Collection(name, opts...)
}

func (r *mongoAdapter) GetType() adapter.AdapterType {
	return r.adapterType
}

func (r *mongoAdapter) Logger() logger.Logger {
	return logger.GetLogger()
}

func (r *mongoAdapter) Config() config.Config {
	return r.config
}

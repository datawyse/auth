package mongodb

import (
	"context"
	"time"

	"auth/internal"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

// MongoDb contains the Mongo client and database objects
type MongoDb struct {
	Client *mongo.Client
	Db     *mongo.Database
	logger *zap.Logger
}

// NewMongoDb configures the MongoDB client and initializes the database connection.
func NewMongoDb(lifecycle fx.Lifecycle, ctx context.Context, log *zap.Logger, config *internal.AppConfig) (*MongoDb, error) {
	log.Info("Connecting to MongoDB", zap.String("url", config.DatabaseURI), zap.String("database", config.DatabaseName))

	client, err := mongo.NewClient(options.Client().ApplyURI(config.DatabaseURI).SetRegistry(UUIDRegistry))
	db := client.Database(config.DatabaseName)

	if err != nil {
		log.Fatal("Error creating MongoDB client", zap.Error(err))
		return nil, err
	}

	lifecycle.Append(fx.Hook{
		OnStart: func(context.Context) error {
			mongoCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
			defer cancel()

			err := client.Connect(mongoCtx)
			if err != nil {
				return err
			}

			// Ping the primary
			if err := client.Ping(mongoCtx, readpref.Primary()); err != nil {
				cancel()
				return err
			}
			log.Info("Successfully connected and pinged.")

			return nil
		},
		OnStop: func(context.Context) error {
			log.Debug("Disconnecting from MongoDB", zap.String("uri", config.DatabaseURI), zap.String("database", config.DatabaseName))
			if err := client.Disconnect(ctx); err != nil {
				log.Error(err.Error())
				return err
			}

			log.Info("Successfully disconnected from MongoDB")
			return nil
		},
	})

	mongoDb := MongoDb{Client: client, Db: db, logger: log}
	return &mongoDb, nil
}

// GetCollection returns a MongoDB collection
func (database *MongoDb) GetCollection(name string) *mongo.Collection {
	return database.Db.Collection(name)
}

// GetClient returns the MongoDB client
func (database *MongoDb) GetClient() *mongo.Client {
	return database.Client
}

// GetDatabase returns the MongoDB database
func (database *MongoDb) GetDatabase() *mongo.Database {
	return database.Db
}

var MongoModule = NewMongoDb

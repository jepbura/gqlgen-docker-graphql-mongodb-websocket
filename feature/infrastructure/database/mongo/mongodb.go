package mongo

import (
	"context"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/jepbura/go-server/config"
	"github.com/jepbura/go-server/constant"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

// Target is parameters to get all MongoDBConnection's dependencies
type Target struct {
	fx.In
	MongoURL string `name:"mongo_url" optional:"true"`
	DBHost   string `name:"db_host"`
	DBPort   string `name:"db_port"`
	DBUser   string `name:"db_user"`
	DBPass   string `name:"db_pass"`
	Lc       fx.Lifecycle
	Logger   *zap.Logger
}

// Connection is connection provider to access to global mongodb client
type Connection struct {
	client *mongo.Client
}

func NewMongoDatabase(target Target) (*Connection, error) {
	if target.MongoURL == "" {
		return nil, nil
	}

	dbHost := target.DBHost
	dbPort := target.DBPort
	dbUser := target.DBUser
	dbPass := target.DBPass

	dbHost = config.DefaultIfEmpty(dbHost, string(constant.DBHost))
	dbPort = config.DefaultIfEmpty(dbPort, string(constant.DBPort))

	mongodbURI := fmt.Sprintf("mongodb://%s:%s@%s:%s", dbUser, dbPass, dbHost, dbPort)

	if dbUser == "" || dbPass == "" {
		mongodbURI = fmt.Sprintf("mongodb://%s:%s", dbHost, dbPort)
	}

	clientOptions := options.Client().ApplyURI(mongodbURI)
	client, err := mongo.Connect(context.TODO(), clientOptions)

	// Error check
	if err != nil {
		log.Fatal(err)
	}

	// Connect check
	err = client.Ping(context.TODO(), nil)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB!")

	target.Lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			target.Logger.Info("Disconnect MongoDB server at " + target.MongoURL + ".")
			return client.Disconnect(ctx)
		},
	})
	return &Connection{
		client,
	}, err
}

// Client is a getter for the client field.
func (c *Connection) Client() *mongo.Client {
	return c.client
}

type key string

const (
	// mongoClient key for mongo session in each request
	mongoClient key = "mongo_client"
)

// Connect is method return adpater for http request that
// inject the database client in context
func (m *Connection) Connect() gin.HandlerFunc {
	return func(c *gin.Context) {
		if m != nil {
			// save it in the mux context
			ctx := context.WithValue(c.Request.Context(), mongoClient, m.client)
			c.Request = c.Request.WithContext(ctx)
		} else {
			// TODO: Warn
			log.Println("Warning: Connection is nil")
		}
		// pass execution to the original handler
		c.Next()
	}
}

// WithContext is method apply mongoClient into context
func (m *Connection) WithContext(ctx context.Context) context.Context {
	if m != nil {
		// save it in the mux context
		return context.WithValue(ctx, mongoClient, m.client)
	} else {
		// TODO: Warn
		log.Println("Warning: Connection is nil")
	}
	return ctx
}

// ForContext is method to get mongodb client from context
func ForContext(ctx context.Context) *mongo.Client {
	client, ok := ctx.Value(mongoClient).(*mongo.Client)
	if !ok {
		panic("ctx passing is not contain mongodb client")
	}
	return client
}

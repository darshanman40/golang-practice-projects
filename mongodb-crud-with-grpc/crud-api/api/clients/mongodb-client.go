package client

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//MongoClient ...
type MongoClient struct {
	db       *mongo.Client
	mongoCtx context.Context
	blogdb   *mongo.Collection
}

//InitClient ...
func (m *MongoClient) InitClient() {
	fmt.Println("Connecting to MongoDB...")
	m.mongoCtx = context.Background()
	db, err := mongo.Connect(m.mongoCtx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}
	m.db = db
	if err := db.Ping(m.mongoCtx, nil); err != nil {
		log.Fatalf("Could not connect to MongoDB: %v\n", err)
	} else {
		fmt.Println("Connected to Mongodb")
	}

	m.blogdb = db.Database("mydb").Collection("blog")
}

//StopClient ...
func (m *MongoClient) StopClient() {
	fmt.Println("Closing MongoDB connection")
	m.db.Disconnect(m.mongoCtx)
}

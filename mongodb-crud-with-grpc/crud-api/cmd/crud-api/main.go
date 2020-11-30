package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"

	"github.com/darshanman40/golang-practice-projects/mongodb-crud-with-grpc/crud-api/api/services"
	"github.com/darshanman40/golang-practice-projects/mongodb-crud-with-grpc/crud-api/configs"
	blogpb "github.com/darshanman40/golang-practice-projects/mongodb-crud-with-grpc/crud-api/internal/proto"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Func main should be as small as possible and do as little as possible by convention
func main() {

	// Generate our config based on the config supplied
	// by the user in the flags
	cfgPath, err := configs.ParseFlags()
	if err != nil {
		log.Fatal(err)
	}
	cfg, err := configs.NewConfig(cfgPath)
	if err != nil {
		log.Fatal(err)
	}
	serverPortString := cfg.GetServerPortString()
	mongodbURL := cfg.GetMongoDBURLString()

	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Println("Starting server on port ", serverPortString)

	listener, err := net.Listen(cfg.Server.Network, serverPortString)
	if err != nil {
		log.Fatal("Unable to listen on port ", serverPortString, err)
	}

	crudServer := services.CrudServer{}
	crudServer.InitServer(listener)

	// Initialize MongoDb client
	fmt.Println("Connecting to MongoDB...")
	mongoCtx := context.Background()
	db, err := mongo.Connect(mongoCtx, options.Client().ApplyURI(mongodbURL))
	if err != nil {
		log.Fatal(err)
	}
	err = db.Ping(mongoCtx, nil)
	if err != nil {

		log.Fatalf("Could not connect to MongoDB: %v\n", err)
	} else {
		fmt.Println("Connected to Mongodb ")
	}

	srv := &services.BlogServiceServer{}
	srv.Init(db, "blog")
	blogpb.RegisterBlogServiceServer(crudServer.GetGRPCServer(), srv)

	go func() {
		crudServer.StartServer()
	}()
	fmt.Println("Server succesfully started on port", serverPortString)

	// Create a channel to receive OS signals
	c := make(chan os.Signal)

	// Relay os.Interrupt to our channel (os.Interrupt = CTRL+C)
	// Ignore other incoming signals
	signal.Notify(c, os.Interrupt)

	// Block main routine until a signal is received
	// As long as user doesn't press CTRL+C a message is not passed
	// And our main routine keeps running
	// If the main routine were to shutdown so would the child routine that is Serving the server
	<-c

	fmt.Println("\nStopping the server...")
	crudServer.StopServer()
	fmt.Println("Closing MongoDB connection")
	db.Disconnect(mongoCtx)
	fmt.Println("Done.")
}

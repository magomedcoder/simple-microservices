package main

import (
	"context"
	"fmt"
	"github.com/magomedcoder/simple-microservice/logger-service/internal"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"net/http"
	"net/rpc"
	"os"
	"time"
)

var client *mongo.Client

func main() {
	mongoClient, err := internal.ConnectToMongo()
	if err != nil {
		log.Panic(err)
	}

	client = mongoClient
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()
	app := internal.Config{
		Models: internal.New(client),
	}
	err = rpc.Register(new(internal.RPCServer))
	go app.RpcListen()
	go app.GRPCListen()
	port := os.Getenv("PORT")
	log.Printf("начало работы службы логирования на порту %s\n", port)
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: app.Routes(),
	}

	if err = srv.ListenAndServe(); err != nil {
		log.Panic()
	}
}

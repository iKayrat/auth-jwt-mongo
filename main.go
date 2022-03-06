package main

import (
	"context"
	"log"
	"time"

	"github.com/iKayrat/auth-jwt-mongo/api"
	"github.com/iKayrat/auth-jwt-mongo/util"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {

	// load config parametres from env
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("Cannot load config:", err)
	}

	// parent context
	parent := context.Background()
	ctx, cancel := context.WithTimeout(parent, time.Second*15)
	defer cancel()

	// connection to db
	client, err := mongo.NewClient(options.Client().ApplyURI(config.MONGODB_ADDR))
	if err != nil {
		log.Fatal(err)
	}

	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		if err := client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	// new server
	server := api.NewServer(client)

	if err := server.Router.Run("127.0.0.1:8080"); err != nil {
		log.Fatal(err)
	}

}

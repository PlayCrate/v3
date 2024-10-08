package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/kattah7/v3/api"
	"github.com/kattah7/v3/models"
	"github.com/kattah7/v3/storage"

	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

type Person struct {
	Name  string
	Age   int
	Email string
}

func main() {
	configFile := flag.String("config", "config.json", "json config file")
	flag.Parse()

	cfg := models.NewConfig(*configFile)

	var DBNumber int = 0
	if cfg.Prod {
		log.Println("Running in production mode")
	} else {
		DBNumber = 1
		log.Println("Running in development mode")
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",       // no password set
		DB:       DBNumber, // use default DB
	})

	store, err := storage.NewPostgresStore(ctx, cfg, rdb)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		log.Fatal(err)
	}
	defer store.Close()

	if err := store.Init(); err != nil {
		log.Fatal(err)
	}

	server := api.NewAPIServer(context.Background(), cfg, store, rdb)
	server.Run()
}

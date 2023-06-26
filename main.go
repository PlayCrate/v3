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
)

func main() {
	configFile := flag.String("config", "config.json", "json config file")
	flag.Parse()

	cfg := models.NewConfig(*configFile)

	store, err := storage.NewPostgresStore(context.Background(), cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		log.Fatal(err)
	}

	defer store.Close()

	if err := store.Init(); err != nil {
		log.Fatal(err)
	}

	server := api.NewAPIServer(cfg, store)
	server.Run()
}

package main

import (
	"github.com/trenchesdeveloper/social-blue/config"
	"github.com/trenchesdeveloper/social-blue/internal/store"
	"log"
)

func main() {
	cfg, err := config.LoadConfig(".")

	if err != nil {
		panic(err)
	}
	storage := store.NewStorage(nil)
	app := &server{
		config: cfg,
		store:  storage,
	}
	mux := app.mount()
	if err := app.start(mux); err != nil {
		log.Fatal(err)
	}

}

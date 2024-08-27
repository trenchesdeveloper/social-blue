package main

import (
	"github.com/trenchesdeveloper/social-blue/config"
	"log"
)

func main() {
	cfg, err := config.LoadConfig(".")

	if err != nil {
		panic(err)
	}
	app := &server{
		config: cfg,
	}
	mux := app.mount()
	if err := app.start(mux); err != nil {
		log.Fatal(err)
	}

}

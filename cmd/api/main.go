package main

import "log"

func main() {
	cfg := config{
		addr: ":8000",
	}
	app := &server{
		config: cfg,
	}
	mux := app.mount()
	if err := app.start(mux); err != nil {
		log.Fatal(err)
	}

}

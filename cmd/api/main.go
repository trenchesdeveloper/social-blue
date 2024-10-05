package main

import (
	"context"
	"database/sql"
	"github.com/trenchesdeveloper/social-blue/config"
	db "github.com/trenchesdeveloper/social-blue/internal/db/sqlc"
	"log"
	"time"
)

const version = "0.0.1"

//	@title			Social Blue API
//	@description	This is a social media API
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger.io

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

//	@BasePath					/v1
//	@securityDefinitions.apikey	ApiKeyAuth
//	@in							header
//	@name						Authorization
func main() {
	cfg, err := config.LoadConfig(".")

	if err != nil {
		panic(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	conn, err := sql.Open(cfg.DBdriver, cfg.DBSource)
	if err != nil {
		log.Fatal(err)
	}
	conn.SetMaxOpenConns(30)
	conn.SetMaxIdleConns(30)
	err = conn.PingContext(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	log.Println("Database connection established")
	storage := db.NewStore(conn)

	app := &server{
		config: cfg,
		store:  storage,
	}
	mux := app.mount()
	if err := app.start(mux); err != nil {
		log.Fatal(err)
	}

}

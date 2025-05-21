package main

import (
	"algo-arena-be/internals/config"
	"algo-arena-be/internals/dbconn"
	"algo-arena-be/internals/migrate"
	"algo-arena-be/internals/server"
	"algo-arena-be/internals/services"
	"context"
	"log"
	"net/http"
	"os"
)

func main() {
	if os.Getenv("ENVIRONMENT") != "PRODUCTION" {
		err := config.LoadConfigurationFile(".env")
		if err != nil {
			log.Fatal("Error loading .env : ", err)
		}
	}
	cfg, err := config.GetConfig()
	if err != nil {
		log.Fatal("Error loading config : ", err)
	}

	db, err := dbconn.GetPostgresDb(cfg.DB_URI)
	if err != nil {
		log.Fatal("Error connecting to postgres db: ", err)
	}

	defer func() {
		if err = dbconn.ClosePostgresDbConn(db); err != nil {
			log.Println("Error closing postgres db: ", err)
		}
	}()

	log.Println("Postgres connected successfully! ")

	if err := migrate.MigratePostgres(context.Background(), db); err != nil {
		log.Fatal("Error migrating database : ", err)
	}

	problemService, submissionService, err := services.GetServices(db)
	if err != nil {
		log.Fatal("Error initializing services : ", err)
	}
	log.Println("[+] Services Initialized")

	r := server.NewChiRouter()
	server.RegisterRoutes(r, problemService, submissionService)

	log.Println("[+] Routes registered")

	s := http.Server{
		Addr:    cfg.SERVER_PORT,
		Handler: r,
	}

	server.RunWithGracefulShutdown(&s)
}

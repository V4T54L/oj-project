package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"sync"
)

func main() {
	redisService := NewRedisService(redisUrl)

	postgres, err := NewPostgreSQLDB(dbUrl, maxIdleConns, maxOpenConns)
	if err != nil {
		log.Fatalf("Error initializing PostgreSQL: %v", err)
	}

	defer func() {
		if err := postgres.Close(); err != nil {
			log.Printf("Error closing the PostgreSQL connection: %v", err)
		}
	}()

	db := postgres.conn

	srv := NewService(db, redisService.ExecuteCode)

	h := NewHandler(srv)
	// Set up the routes
	r := h.Routes()

	ctx, cancel := context.WithCancel(context.TODO())
	defer cancel()

	var wg sync.WaitGroup

	redisService.StartResultWorker(ctx, func(er *ExecutionResponse) {
		log.Printf("\n\n\n### Results received:\n ID: %d \n Type: %s \n\n\n", er.SubmissionID, er.ExecutionType)
	}, &wg)

	// Set up the server
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", serverPort),
		Handler: r,
	}

	// Start the server
	log.Println("Server started on :8080")
	if err := server.ListenAndServe(); err != nil {
		log.Fatal("Server failed to start: ", err)
	}

	wg.Wait()

}

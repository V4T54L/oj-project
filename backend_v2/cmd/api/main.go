package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	authmodule "online-judge/internal/auth_module"
	"online-judge/internal/config"
	"online-judge/internal/handlers"
	"online-judge/internal/models"
	"online-judge/internal/repo"
	"online-judge/internal/router"
	"online-judge/internal/server"
	"online-judge/internal/services"
	"os"
	"sync"
	"time"
)

func main() {
	if os.Getenv("PORT") != "PRODUCTION" {
		_ = config.LoadEnv(".env")
	}
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	auth := authmodule.NewJWTAuth([]byte(cfg.JWT_SECRET), time.Minute*time.Duration(cfg.TOKEN_EXPIRY_MINUTES))

	redisClient := services.NewRedisService(cfg.REDIS_ADDR)
	defer redisClient.Close()

	ctx, cancel := context.WithCancel(context.TODO())
	defer cancel()

	var wg sync.WaitGroup

	problemRepo := repo.NewProblemRepo()
	submissionRepo := repo.NewSubmissionRepo()

	redisClient.StartResultWorker(ctx, func(ecr *models.ExecuteCodeResponse) {
		status, runtime, memory := "Accepted", 0, 0
		for i, v := range ecr.TestCaseResults {
			if v.RuntimeMS > runtime {
				runtime = v.RuntimeMS
			}
			if v.MemoryKB > memory {
				memory = v.MemoryKB
			}
			if v.Status != "Accepted" {
				status = fmt.Sprintf("%s on Test Case : %d", v.Status, i+1)
			}
		}
		log.Println("Result received: ", status, runtime, memory)
		if ecr.ExecutionType == "submission" {
			submissionRepo.UpdateSubmission(ctx, ecr.ID, runtime, memory, status)
		} else if ecr.ExecutionType == "validation" {
			problemRepo.UpdateProblemStatusByID(ctx, ecr.ID, status)
		}
	}, &wg)

	handler, err := handlers.NewHandler(submissionRepo, problemRepo, redisClient)
	if err != nil {
		log.Fatalf("Failed to load handler: %v", err)
	}

	r := router.NewChiRouter(nil, auth, *handler)

	s := http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.SERVER_PORT),
		Handler: r,
	}

	server.RunWithGracefulShutdown(&s)
	wg.Wait()
}

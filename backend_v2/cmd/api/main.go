package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	authmodule "online-judge/internal/auth_module"
	"online-judge/internal/config"
	"online-judge/internal/models"
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

	redisClient.StartResultWorker(ctx, func(ecr *models.ExecuteCodeResponse) {
		// TODO
		fmt.Println("TODO: Handle Execution Result based on ExecutionType Result : ", ecr)
	}, &wg)

	r := router.NewChiRouter(nil, auth, redisClient.ExecuteCode)

	s := http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.SERVER_PORT),
		Handler: r,
	}

	server.RunWithGracefulShutdown(&s)
	wg.Wait()
}

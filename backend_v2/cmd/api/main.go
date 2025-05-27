package main

import (
	"log"
	"net/http"
	authmodule "online-judge/internal/auth_module"
	"online-judge/internal/config"
	"online-judge/internal/router"
	"online-judge/internal/server"
	"os"
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

	r := router.NewChiRouter(nil, auth)

	s := http.Server{
		Addr:    ":8000",
		Handler: r,
	}

	server.RunWithGracefulShutdown(&s)
}

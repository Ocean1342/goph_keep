package main

import (
	"context"
	_ "github.com/jackc/pgx/v5/stdlib"
	log "github.com/sirupsen/logrus"
	"goph_keeper/config"
	"goph_keeper/internal/api"
	"goph_keeper/internal/migrator"
	"goph_keeper/internal/server"
	"goph_keeper/internal/storage"
	"goph_keeper/pkg/auth"
	"os"
)

func init() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)
}

// @title Goph_Keeper Server
// @version 1.0
// @description Это минимальный пример API с Swagger
// @host localhost:8080
func main() {
	cfg := config.New()
	ctx := context.Background()
	repo := storage.New(ctx, cfg.DatabaseURL)
	err := migrator.Migrate(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("could not migrate database: %v", err)
	}
	handler := api.New(repo, auth.New(cfg.SecretKey, cfg.TokenTTL))
	server.Init(cfg, handler)
}

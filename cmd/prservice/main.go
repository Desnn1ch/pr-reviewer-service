package main

import (
	"context"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/Desnn1ch/pr-reviewer-service/internal/app/service"
	"github.com/Desnn1ch/pr-reviewer-service/internal/config"
	"github.com/Desnn1ch/pr-reviewer-service/internal/domain/common"
	dbinfra "github.com/Desnn1ch/pr-reviewer-service/internal/infrastructure/persistence/db"
	httpserver "github.com/Desnn1ch/pr-reviewer-service/internal/interface/httpserver"
	"github.com/Desnn1ch/pr-reviewer-service/internal/interface/httpserver/handler"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	configPath := config.Getenv("CONFIG_PATH", "internal/config/config.yaml")

	cfg, err := config.Load(configPath)
	if err != nil {
		log.Printf("failed to load config: %v", err)
		return
	}

	db, err := dbinfra.New(ctx, dbinfra.Config{
		DSN:             cfg.Database.DSN(),
		MigrationsDir:   "./migrations",
		MaxOpenConns:    cfg.Database.Pool.MaxOpenConns,
		MaxIdleConns:    cfg.Database.Pool.MaxIdleConns,
		ConnMaxLifetime: cfg.Database.Pool.ConnMaxLifetime.Duration,
	})
	if err != nil {
		log.Printf("failed to init db: %v", err)
		return
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Printf("db close error: %v", err)
		}
	}()

	clock := common.StandardClock{}

	repos := dbinfra.NewRepositories(db)

	teamSvc := service.NewTeamService(repos.Teams, repos.Users, repos.Tx)
	userSvc := service.NewUserService(repos.Users, repos.PRs)
	prSvc := service.NewPRService(repos.PRs, repos.Users, repos.Tx, clock)

	teamHandler := handler.NewTeamHandler(teamSvc)
	userHandler := handler.NewUserHandler(userSvc)
	prHandler := handler.NewPRHandler(prSvc)

	router := httpserver.NewRouter(teamHandler, userHandler, prHandler)

	srv := &http.Server{
		Addr:         cfg.Server.Address,
		Handler:      router,
		ReadTimeout:  cfg.Server.ReadTimeout.Duration,
		WriteTimeout: cfg.Server.WriteTimeout.Duration,
		IdleTimeout:  cfg.Server.IdleTimeout.Duration,
	}

	go func() {
		log.Printf("listening on %s", cfg.Server.Address)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	<-ctx.Done()
	log.Println("shutdown signal received")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("server shutdown error: %v", err)
	} else {
		log.Println("server shutdown complete")
	}

	if err := db.Close(); err != nil {
		log.Printf("db close error: %v", err)
	} else {
		log.Println("db connection closed")
	}
}

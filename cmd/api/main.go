package main

import (
	"airops/internal/infrastructure/postgres"
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "airops/docs"
	"airops/internal/app"
)

// @title           airops API
// @version         1.0
// @description     Demo API for flights/bookings/stats
// @BasePath        /api/v1
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name X-API-Key
func main() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL is required")
	}

	ctx := context.Background()

	pool, err := postgres.NewPool(ctx, dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	a := app.New(pool, ":8080")

	go func() {
		if err := a.Run(); err != nil {
			log.Fatal(err)
		}
	}()

	// graceful shutdown
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	<-sigCh

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_ = a.Shutdown(shutdownCtx)
}

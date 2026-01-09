// cmd/api/main.go
package main

import (
	"airops/internal/infrastructure/postgres"
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"airops/internal/app"
)

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

	// старт
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

package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"airops/internal/app"
	"airops/internal/infra/db/pg"
)

func main() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL is required")
	}

	ctx := context.Background()

	pool, err := pg.NewPool(ctx, dsn) // или как у тебя создаётся pool
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	a := app.New(pool)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: a.Handler,
	}

	log.Println("listening on :8080")

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	<-sigCh

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_ = srv.Shutdown(shutdownCtx)
}

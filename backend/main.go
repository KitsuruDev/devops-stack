// backend/main.go
package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	initDB()
	ensureSchema(context.Background())

	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", health)
	mux.HandleFunc("GET /api/notes", listNotes)
	mux.HandleFunc("POST /api/notes", addNote)
	mux.HandleFunc("GET /api/notes/{id}", getNote)

	srv := &http.Server{
		Addr: ":8000",
		Handler: mux,
		ReadTimeout: 5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	// Graceful shutdown: SIGTERM/SIGINT → http.Server.Shutdown
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer stop()

	go func() {
		log.Println("notes-api listening on :8000")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %v", err)
		}
	}()

	<-ctx.Done()
	log.Println("SIGTERM received, shutting down gracefully...")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("graceful shutdown failed: %v", err)
		os.Exit(1)
	}
	log.Println("server stopped")
}

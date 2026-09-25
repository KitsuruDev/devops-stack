// backend/database.go
package main

import (
	"context"
	"log"
	"os"
	"github.com/jackc/pgx/v5/pgxpool"
)

var pool *pgxpool.Pool

func initDB() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "postgres://notes:notes@db:5432/notes"
	}
	ctx := context.Background()
	cfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		log.Fatalf("parse DSN: %v", err)
	}
	pool, err = pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		log.Fatalf("connect: %v", err)
	}
	if err = pool.Ping(ctx); err != nil {
		log.Printf("DB not reachable at start: %v (will retry per request)", err)
	}
}
func ensureSchema(ctx context.Context) {
	_, err := pool.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS notes (
			id SERIAL PRIMARY KEY,
			title TEXT NOT NULL,
			body TEXT NOT NULL DEFAULT ''
		);`)
	if err != nil {
		log.Printf("ensure schema: %v", err)
	}
}

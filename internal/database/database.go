package database

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"

	"scan_to_score/internal/config"
)

type Service interface {
	Health(ctx context.Context) map[string]string
	Close()
	Pool() *pgxpool.Pool
}

type service struct {
	pool *pgxpool.Pool
}

func (s *service) Health(ctx context.Context) map[string]string {
	err := s.pool.Ping(ctx)
	if err != nil {
		return map[string]string{"status": "down"}
	}
	return map[string]string{"status": "up"}
}

func (s *service) Close() {
	s.pool.Close()
}

func (s *service) Pool() *pgxpool.Pool {
	return s.pool
}

var dbInstance *service

func New(ctx context.Context, cfg *config.Config) Service {
	if dbInstance != nil {
		return dbInstance
	}

	connStr := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable&search_path=%s",
		cfg.DB.User, cfg.DB.Password, cfg.DB.Host, cfg.DB.Port, cfg.DB.Name, cfg.DB.Schema,
	)

	pool, err := pgxpool.New(ctx, connStr)
	if err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}

	dbInstance = &service{pool: pool}
	return dbInstance
}

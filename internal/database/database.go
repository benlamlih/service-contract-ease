package database

import (
	"database/sql"
	"fmt"
	"log"
	"my_project/internal/config"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type Service interface {
	Health() map[string]string
	Close() error
}

type service struct {
	db *sql.DB
}

func (s *service) Health() map[string]string {
	return map[string]string{"status": "up"}
}

func (s *service) Close() error {
	return s.db.Close()
}

var dbInstance *service

func New(cfg *config.Config) Service {
	if dbInstance != nil {
		return dbInstance
	}

	connStr := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable&search_path=%s",
		cfg.DB.User, cfg.DB.Password, cfg.DB.Host, cfg.DB.Port, cfg.DB.Name, cfg.DB.Schema,
	)

	db, err := sql.Open("pgx", connStr)
	if err != nil {
		log.Fatal(err)
	}

	dbInstance = &service{db: db}
	return dbInstance
}

package database_test

import (
	"context"
	"log"
	"os"
	"testing"
	"time"

	"my_project/internal/config"
	"my_project/internal/database"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

var cfg *config.Config

func mustStartPostgresContainer() (func(context.Context, ...testcontainers.TerminateOption) error, error) {
	const (
		dbName = "test_db"
		dbUser = "test_user"
		dbPass = "test_pass"
	)

	container, err := postgres.Run(
		context.Background(),
		"postgres:15-alpine",
		postgres.WithDatabase(dbName),
		postgres.WithUsername(dbUser),
		postgres.WithPassword(dbPass),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithStartupTimeout(10*time.Second),
		),
	)
	if err != nil {
		return nil, err
	}

	host, err := container.Host(context.Background())
	if err != nil {
		return container.Terminate, err
	}

	port, err := container.MappedPort(context.Background(), "5432/tcp")
	if err != nil {
		return container.Terminate, err
	}

	cfg = &config.Config{}
	cfg.DB.Name = dbName
	cfg.DB.User = dbUser
	cfg.DB.Password = dbPass
	cfg.DB.Host = host
	cfg.DB.Port = port.Port()
	cfg.DB.Schema = "public"
	cfg.App.Port = "8080"

	return container.Terminate, nil
}

func TestMain(m *testing.M) {
	teardown, err := mustStartPostgresContainer()
	if err != nil {
		log.Fatalf("could not start postgres container: %v", err)
	}

	code := m.Run()

	if teardown != nil {
		_ = teardown(context.Background())
	}

	os.Exit(code)
}

func TestNew(t *testing.T) {
	srv := database.New(cfg)
	if srv == nil {
		t.Fatal("New() returned nil")
	}
}

func TestHealth(t *testing.T) {
	srv := database.New(cfg)

	stats := srv.Health()

	if stats["status"] != "up" {
		t.Fatalf("expected status to be up, got %s", stats["status"])
	}
	if _, ok := stats["error"]; ok {
		t.Fatalf("expected error not to be present")
	}
	if stats["message"] != "It's healthy" {
		t.Fatalf("expected message to be 'It's healthy', got %s", stats["message"])
	}
}

func TestClose(t *testing.T) {
	srv := database.New(cfg)

	if err := srv.Close(); err != nil {
		t.Fatalf("expected Close() to return nil, got %v", err)
	}
}

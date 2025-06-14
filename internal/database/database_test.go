package database

import (
	"context"
	"log"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"

	"contract_ease/internal/config"
)

var (
	ctx        context.Context
	container  testcontainers.Container
	testConfig *config.Config
	mu         sync.Mutex // Mutex to protect dbInstance access
)

func TestMain(m *testing.M) {
	ctx = context.Background()
	var err error

	container, err = postgres.Run(ctx,
		"postgres:16-alpine",
		postgres.WithDatabase("test_db"),
		postgres.WithUsername("test_user"),
		postgres.WithPassword("test_pass"),
		testcontainers.WithWaitStrategy(
			wait.ForAll(
				wait.ForLog("database system is ready to accept connections").WithOccurrence(2),
				wait.ForListeningPort("5432/tcp"),
			).WithDeadline(time.Second*30),
		),
	)
	if err != nil {
		log.Fatalf("Failed to start container: %v", err)
	}

	host, err := container.Host(ctx)
	if err != nil {
		log.Fatalf("Failed to get container host: %v", err)
	}

	mappedPort, err := container.MappedPort(ctx, "5432")
	if err != nil {
		log.Fatalf("Failed to get container port: %v", err)
	}

	// Setup test config
	testConfig = &config.Config{}
	testConfig.DB.Name = "test_db"
	testConfig.DB.User = "test_user"
	testConfig.DB.Password = "test_pass"
	testConfig.DB.Host = host
	testConfig.DB.Port = mappedPort.Port()
	testConfig.DB.Schema = "public"

	code := m.Run()

	if err := container.Terminate(ctx); err != nil {
		log.Printf("Failed to terminate container: %v", err)
	}

	os.Exit(code)
}

func TestHealth(t *testing.T) {
	t.Parallel()
	mu.Lock()
	dbInstance = nil
	mu.Unlock()

	srv := New(ctx, testConfig)
	defer func() {
		srv.Close()
		mu.Lock()
		dbInstance = nil
		mu.Unlock()
	}()

	health := srv.Health(ctx)
	assert.Equal(t, "up", health["status"], "Expected health status to be 'up'")
}

func TestClose(t *testing.T) {
	t.Parallel()
	mu.Lock()
	dbInstance = nil
	mu.Unlock()

	srv := New(ctx, testConfig)
	srv.Close()

	mu.Lock()
	dbInstance = nil
	mu.Unlock()
}

func TestPool(t *testing.T) {
	t.Parallel()
	mu.Lock()
	dbInstance = nil
	mu.Unlock()

	srv := New(ctx, testConfig)
	defer func() {
		srv.Close()
		mu.Lock()
		dbInstance = nil
		mu.Unlock()
	}()

	pool := srv.Pool()
	assert.NotNil(t, pool, "Expected pool to not be nil")

	err := pool.Ping(ctx)
	assert.NoError(t, err, "Expected pool to be able to ping database")
}

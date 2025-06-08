package server

import (
	"context"
	"fmt"
	"go.opentelemetry.io/otel/trace"
	"net/http"
	"strconv"
	"time"

	"scan_to_score/internal/config"
	"scan_to_score/internal/database"
)

type Server struct {
	port int
	db   database.Service
}

func NewServer(ctx context.Context, tp trace.TracerProvider) *http.Server {
	cfg := config.LoadConfig()

	port, _ := strconv.Atoi(cfg.App.Port)

	s := &Server{
		port: port,
		db:   database.New(ctx, cfg, tp),
	}

	return &http.Server{
		Addr:         fmt.Sprintf(":%d", s.port),
		Handler:      s.RegisterRoutes(tp),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}
}

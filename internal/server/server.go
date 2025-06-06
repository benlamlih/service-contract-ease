package server

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"my_project/internal/config"
	"my_project/internal/database"
)

type Server struct {
	port int
	db   database.Service
}

func NewServer() *http.Server {
	cfg := config.LoadConfig()

	port, _ := strconv.Atoi(cfg.App.Port)

	s := &Server{
		port: port,
		db:   database.New(cfg),
	}

	return &http.Server{
		Addr:         fmt.Sprintf(":%d", s.port),
		Handler:      s.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}
}

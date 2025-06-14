package server

import (
	"contract_ease/internal/database"
)

type Server struct {
	port int
	db   database.Service
}

func NewServer(db database.Service, port int) *Server {
	return &Server{db: db, port: port}
}

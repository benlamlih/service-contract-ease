package server

import (
	"contract_ease/internal/database"
	"contract_ease/internal/repository"
)

type Server struct {
	port          int
	db            database.Service
	store         *repository.Store
	ZitadelClient ZitadelClient
}

func NewServer(db database.Service, port int, client ZitadelClient) *Server {
	return &Server{
		db:            db,
		store:         repository.NewStore(db.Pool()),
		port:          port,
		ZitadelClient: client,
	}
}

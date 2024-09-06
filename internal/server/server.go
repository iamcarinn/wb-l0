package server

import (
	"wb-l0/internal/handler"
	"net/http"
	"log"
)

type Server struct {
	httpServer http.Server
}

// Создание нового сервера
func NewServer() *Server {
	return &Server{
		httpServer: http.Server{
			Addr:    "127.0.0.1:3330",
			Handler: &handler.Handler{},
		},
	}
}

// Запуск сервера
func (s *Server) Run() error {
	log.Println("server is running")
	return s.httpServer.ListenAndServe()
}

// Остановка сервера
func (s *Server) Stop() error {
	return s.httpServer.Close()
}
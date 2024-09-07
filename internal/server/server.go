package server

import (
	"net/http"
	"log"
	_ "github.com/lib/pq"
)

type Server struct {
	httpServer http.Server
}

// Создание нового сервера
func New(h http.Handler) *Server {

	// Создаем и возвращаем сервер с подключением к БД
	return &Server{
		httpServer: http.Server{
			Addr:    "127.0.0.1:3330",
			Handler: h, // Передаем подключение к базе данных в обработчик
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
package server

import (
	"wb-l0/internal/handler"
	"net/http"
	"log"
	"database/sql"
	_ "github.com/lib/pq"
)

type Server struct {
	httpServer http.Server
	db         *sql.DB
}

// Создание нового сервера
func NewServer() *Server {

	connStr := "user=iamcarinn dbname=iamcarinn password=1307 host=localhost sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
	}

	// Создаем и возвращаем сервер с подключением к БД
	return &Server{
		httpServer: http.Server{
			Addr:    "127.0.0.1:3330",
			Handler: handler.NewHandler(db), // Передаем подключение к базе данных в обработчик
		},
		db: db,
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
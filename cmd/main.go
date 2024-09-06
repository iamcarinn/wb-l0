package main

import (
	"wb-l0/internal/server"
	"log"
)

func main() {
	s := server.NewServer() // Создаем сервер
	if err := s.Run(); err != nil {	// Запускаем сервер
		log.Fatalln(err)
	}
}
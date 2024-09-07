package main

import (
	"log"
	"wb-l0/internal/handler"
	"wb-l0/internal/repo"
	"wb-l0/internal/server"
)

func main() {
	r := repo.New()					// Создаем новое подключение к базе данных
	h := handler.New(r)				// Создаем новый HTTP-обработчик, передавая в него подключение к базе данных
	s := server.New(h.Route()) 		// Создаем новый сервер, передавая ему маршруты обработчика
	if err := s.Run(); err != nil {	// Запускаем сервер
		log.Fatalln(err)
	}
}


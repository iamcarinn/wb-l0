package main

import (
	"log"
	"wb-l0/internal/cache"
	"wb-l0/internal/handler"
	"wb-l0/internal/repo"
	"wb-l0/internal/server"
)

func main() {
	r := repo.New()                  // Создаем новое подключение к базе данных
	c := cache.New()                 // Создаем новый экземпляр кэша

	// Получаем все заказы из базы данных
	orders, err := r.GetAllOrders()
	if err != nil {
		log.Fatalf("Error retrieving orders from database: %v", err)
	}

	// Наполняем кэш данными
	for _, order := range orders {
		c.AddToCache(order)
	}

	h := handler.New(r, c)              // Создаем новый HTTP-обработчик, передавая в него подключение к базе данных и кэш
	s := server.New(h.Route())         // Создаем новый сервер, передавая ему маршруты обработчика
	if err := s.Run(); err != nil {    // Запускаем сервер
		log.Fatalln(err)
	}
}

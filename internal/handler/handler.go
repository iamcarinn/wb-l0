package handler

import (
	"wb-l0/internal/handler/order"
	"log"
	"net/http"
)

type Handler struct{}

// вместо маршрутизатора с URL "/"???
// Вызывается автоматически, когда поступает HTTP-запрос
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Printf("METHOD: [%v] URI: [%v]", r.Method, r.RequestURI)

	order.HandlePost(w, r)	// Обработчик
}
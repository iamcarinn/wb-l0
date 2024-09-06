package handler

import (
	"wb-l0/internal/handler/order"
	"log"
	"net/http"
	"database/sql"
)

type Handler struct {
	db *sql.DB
}

// Вместо маршрутизатора с URL "/"??? 
// Вызывается автоматически, когда поступает HTTP-запрос
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Printf("METHOD: [%v] URI: [%v]", r.Method, r.RequestURI)

	order.HandlePost(w, r, h.db) // Передаем подключение к базе данных в обработчик
}

// Создаем новый обработчик с подключением к базе данных
func NewHandler(db *sql.DB) *Handler {
	return &Handler{db: db}
}


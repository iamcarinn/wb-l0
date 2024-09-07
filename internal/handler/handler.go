package handler

import (
	"encoding/json"
	"net/http"
	"wb-l0/internal/handler/order"
	"wb-l0/internal/repo"
)

type Handler struct {
	repo repo.Repo
}

// Создаем роутер (маршрутизатор)
func (h *Handler) Route() http.Handler {
	router := http.NewServeMux() 	// мультиплексер
	router.HandleFunc("/", h.HandlePost)
	return router
}

// Создаем новый обработчик с подключением к базе данных
func New(repo repo.Repo) *Handler {
	return &Handler{repo: repo}
}

// Функция для обработки POST-запросов
func (h *Handler) HandlePost(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		var req order.OrderRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid input data", http.StatusBadRequest)
			return
		}
		orderData, err := h.repo.GetOrder(req.OrderUID)

		if err != nil {
			http.Error(w, "invalid input data", http.StatusBadRequest)
			return
		}

		// Если заказ найден, возвращаем информацию о заказе
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(orderData)


	default:
		http.Error(w, "http method not allowed", http.StatusMethodNotAllowed)
	}
}


package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"wb-l0/internal/cache"
	"wb-l0/internal/handler/order"
	"wb-l0/internal/repo"
)

type Handler struct {
	repo  repo.Repo
	cache *cache.Cache // Добавляем кэш
}

// // Создаем роутер (маршрутизатор)
// func (h *Handler) Route() http.Handler {
// 	router := http.NewServeMux() // мультиплексер
// 	router.HandleFunc("/", h.HandlePost)
// 	return router
// }

// Создаем роутер (маршрутизатор)
func (h *Handler) Route() http.Handler {
	router := http.NewServeMux() // мультиплексер

	// Маршрут для статических файлов (HTML-страница)
	router.Handle("/", http.FileServer(http.Dir(".")))

	// Маршрут для обработки API-запросов
	router.HandleFunc("/orders", h.HandlePost)

	return router
}


// Создаем новый обработчик с подключением к базе данных и кэшем
func New(repo repo.Repo, cache *cache.Cache) *Handler {
	return &Handler{
		repo:  repo,
		cache: cache,
	}
}

func (h *Handler) HandlePost(w http.ResponseWriter, r *http.Request) {
    log.Printf("Received %s request for %s", r.Method, r.URL.Path)

    switch r.Method {
    case http.MethodPost:
        var req order.OrderRequest
        if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
            http.Error(w, "invalid input data", http.StatusBadRequest)
            return
        }

        // Сначала проверяем кэш
        orderData, found := h.cache.GetFromCache(req.OrderUID)
        if found {
            w.WriteHeader(http.StatusOK)
            json.NewEncoder(w).Encode(orderData)
            return
        }

        // Если в кэше нет, ищем в базе данных
        orderDataPtr, err := h.repo.GetOrder(req.OrderUID)
        if err != nil {
            http.Error(w, "order not found", http.StatusNotFound)
            return
        }

        // Разыменовываем указатель перед добавлением в кэш
        orderData = *orderDataPtr

        // Сохраняем найденный заказ в кэш
        h.cache.AddToCache(orderData)

        // Возвращаем информацию о заказе
        w.WriteHeader(http.StatusOK)
        json.NewEncoder(w).Encode(orderData)

    default:
        http.Error(w, "http method not allowed", http.StatusMethodNotAllowed)
    }
}

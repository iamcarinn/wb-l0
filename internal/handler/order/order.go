package order

import (
	"encoding/json"
	"net/http"
	"os"
	"wb-l0/config"
)

// Структура для запроса, содержащего `order_uid`
type OrderRequest struct {
	OrderUID string `json:"order_uid"`
}

// Структура для ошибки, которую будем возвращать при неправильных запросах
type ErrorResponse struct {
	Error string `json:"error"`
}

// Функция для обработки POST-запросов
func HandlePost(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		var req OrderRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid input data", http.StatusBadRequest)
			return
		}

		// Читаем файл config.json
		data, err := os.ReadFile("config/config.json")
		if err != nil {
			http.Error(w, "unable to read config file", http.StatusInternalServerError)
			return
		}

		// Парсим JSON файл в структуру Order
		var orderData config.Order
		if err := json.Unmarshal(data, &orderData); err != nil {
			http.Error(w, "error parsing config file", http.StatusInternalServerError)
			return
		}

		// Проверяем, существует ли заказ с данным order_uid
		if req.OrderUID != orderData.OrderUID {
			resp := ErrorResponse{Error: "order not found"}
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(resp)
			return
		}

		// Если заказ найден, возвращаем информацию о заказе
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(orderData)

	default:
		http.Error(w, "http method not allowed", http.StatusMethodNotAllowed)
	}
}

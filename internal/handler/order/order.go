package order

import (
	"database/sql"
	"encoding/json"
	"net/http"
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
func HandlePost(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	switch r.Method {
	case http.MethodPost:
		var req OrderRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid input data", http.StatusBadRequest)
			return
		}

		// Выполняем запрос к базе данных
		query := `SELECT order_uid FROM orders WHERE order_uid = $1`
		var orderData config.Order
		err := db.QueryRow(query, req.OrderUID).Scan(&orderData.OrderUID)

		// Если заказ не найден
		if err == sql.ErrNoRows {
			resp := ErrorResponse{Error: "order not found"}
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(resp)
			return
		} else if err != nil {
			http.Error(w, "error querying database", http.StatusInternalServerError)
			return
		}

		// Если заказ найден, возвращаем информацию о заказе
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(orderData)

	default:
		http.Error(w, "http method not allowed", http.StatusMethodNotAllowed)
	}
}

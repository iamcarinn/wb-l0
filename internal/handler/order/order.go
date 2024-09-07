package order

// Структура для запроса, содержащего `order_uid`
type OrderRequest struct {
	OrderUID string `json:"order_uid"`
}

// Структура для ошибки, которую будем возвращать при неправильных запросах
type ErrorResponse struct {
	Error string `json:"error"`
}
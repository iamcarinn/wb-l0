package config

import (
	"encoding/json"
	"os"
	"time"
)

// Order представляет основной объект заказа
type Order struct {
	OrderUID          string    `json:"order_uid"`
	TrackNumber       string    `json:"track_number"`
	Entry             string    `json:"entry"`
	Delivery          Delivery  `json:"delivery"`
	Payment           Payment   `json:"payment"`
	Items             []Item    `json:"items"`
	Locale            string    `json:"locale"`
	InternalSignature string    `json:"internal_signature"`
	CustomerID        string    `json:"customer_id"`
	DeliveryService   string    `json:"delivery_service"`
	Shardkey          string    `json:"shardkey"`
	SmID              int       `json:"sm_id"`
	DateCreated       time.Time `json:"date_created"`
	OofShard          string    `json:"oof_shard"`
}

// Delivery представляет объект доставки
type Delivery struct {
	Name    string `json:"name"`
	Phone   string `json:"phone"`
	Zip     string `json:"zip"`
	City    string `json:"city"`
	Address string `json:"address"`
	Region  string `json:"region"`
	Email   string `json:"email"`
}

// Payment представляет объект платежа
type Payment struct {
	Transaction  string `json:"transaction"`
	RequestID    string `json:"request_id"`
	Currency     string `json:"currency"`
	Provider     string `json:"provider"`
	Amount       int    `json:"amount"`
	PaymentDT    int64  `json:"payment_dt"` // timestamp in seconds
	Bank         string `json:"bank"`
	DeliveryCost int    `json:"delivery_cost"`
	GoodsTotal   int    `json:"goods_total"`
	CustomFee    int    `json:"custom_fee"`
}

// Item представляет объект товара
type Item struct {
	ChrtID      int    `json:"chrt_id"`
	TrackNumber string `json:"track_number"`
	Price       int    `json:"price"`
	RID         string `json:"rid"`
	Name        string `json:"name"`
	Sale        int    `json:"sale"`
	Size        string `json:"size"`
	TotalPrice  int    `json:"total_price"`
	NMID        int    `json:"nm_id"`
	Brand       string `json:"brand"`
	Status      int    `json:"status"`
}

// DatabaseConfig содержит параметры для подключения к базе данных
var DatabaseConfig = struct {
	Username string
	Password string
	Database string
}{
	Username: "iamcarinn",
	Password: "",
	Database: "wbL0",
}

// FillModels читает JSON-файл и заполняет структуры
func (o *Order) FillModels() error {
	data, err := os.ReadFile("config.json")
	if err != nil {
		return err
	}

	err = json.Unmarshal(data, o)
	if err != nil {
		return err
	}

	return nil
}

// GetMapModels преобразует массив Items в карту, где ключ — это название товара, а значение — его цена
func (o *Order) GetMapModels() map[string]int {
	itemMap := make(map[string]int)
	for _, item := range o.Items {
		itemMap[item.Name] = item.Price
	}
	return itemMap
}

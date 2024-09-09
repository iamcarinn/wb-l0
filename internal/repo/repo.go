package repo

import (
	"database/sql"
	"log"

	//"wb-l0/internal/cache"
	"wb-l0/internal/model"
)

type Repo struct {
	db *sql.DB
}

// Создание нового подключения к БД
func New() Repo {
	// Подключение к базе данных
	connStr := "user=iamcarinn dbname=iamcarinn password=1307 host=localhost sslmode=disable"
	// Открытие подключения к базе данных
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
	}

	// Проверка подключения к базе данных
	err = db.Ping()
	if err != nil {
		log.Fatalf("Error ping database: %v", err)
	}

	// Возвращаем новый экземпляр Repo
	return Repo{
		db: db,
	}
}

// Добавление заказа в БД
func (repo *Repo) AddToDB(order model.Order) error {
	// Начинаем транзакцию
	tx, err := repo.db.Begin()
	if err != nil {
		return err
	}

	// Вставка данных в таблицу orders
	orderQuery := `INSERT INTO orders (order_uid, track_number, entry, locale, internal_signature, customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard) 
				   VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`
	_, err = tx.Exec(orderQuery, order.OrderUID, order.TrackNumber, order.Entry, order.Locale, order.InternalSignature, order.CustomerID, order.DeliveryService, order.Shardkey, order.SmID, order.DateCreated, order.OofShard)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Вставка данных в таблицу delivery
	deliveryQuery := `INSERT INTO delivery (order_uid, name, phone, zip, city, address, region, email) 
					  VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
	_, err = tx.Exec(deliveryQuery, order.OrderUID, order.Delivery.Name, order.Delivery.Phone, order.Delivery.Zip, order.Delivery.City, order.Delivery.Address, order.Delivery.Region, order.Delivery.Email)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Вставка данных в таблицу payment
	paymentQuery := `INSERT INTO payment (order_uid, transaction, request_id, currency, provider, amount, payment_dt, bank, delivery_cost, goods_total, custom_fee) 
					 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`
	_, err = tx.Exec(paymentQuery, order.OrderUID, order.Payment.Transaction, order.Payment.RequestID, order.Payment.Currency, order.Payment.Provider, order.Payment.Amount, order.Payment.PaymentDT, order.Payment.Bank, order.Payment.DeliveryCost, order.Payment.GoodsTotal, order.Payment.CustomFee)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Вставка данных в таблицу items
	for _, item := range order.Items {
		itemsQuery := `INSERT INTO items (chrt_id, order_uid, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status) 
					   VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`
		_, err = tx.Exec(itemsQuery, item.ChrtID, order.OrderUID, item.TrackNumber, item.Price, item.RID, item.Name, item.Sale, item.Size, item.TotalPrice, item.NMID, item.Brand, item.Status)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	// Подтверждаем транзакцию
	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func (repo *Repo) GetOrder(id string) (*model.Order, error) {
	// Выполняем запрос к базе данных
	query := `SELECT order_uid FROM orders WHERE order_uid = $1`
	var orderData model.Order
	err := repo.db.QueryRow(query, id).Scan(&orderData.OrderUID)
	// Если заказ не найден
	if err != nil {
		return nil, err
	}
	return &orderData, nil
}

// Возвращает все заказы из базы данных
func (repo *Repo) GetAllOrders() ([]model.Order, error) {
	query := `SELECT order_uid, track_number, entry, locale, internal_signature, customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard 
			  FROM orders`
	rows, err := repo.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []model.Order
	for rows.Next() {
		var order model.Order
		if err := rows.Scan(&order.OrderUID, &order.TrackNumber, &order.Entry, &order.Locale, &order.InternalSignature, &order.CustomerID, &order.DeliveryService, &order.Shardkey, &order.SmID, &order.DateCreated, &order.OofShard); err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return orders, nil
}

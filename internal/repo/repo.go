package repo

import (
	"database/sql"
	"log"

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
	orderQuery := `
		SELECT order_uid, track_number, entry, locale, internal_signature, customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard 
		FROM orders 
		WHERE order_uid = $1`
	var orderData model.Order
	err := repo.db.QueryRow(orderQuery, id).Scan(&orderData.OrderUID, &orderData.TrackNumber, &orderData.Entry, &orderData.Locale, &orderData.InternalSignature, &orderData.CustomerID, &orderData.DeliveryService, &orderData.Shardkey, &orderData.SmID, &orderData.DateCreated, &orderData.OofShard)
	if err != nil {
		return nil, err
	}

	// Получаем данные из таблицы delivery
	deliveryQuery := `SELECT name, phone, zip, city, address, region, email FROM delivery WHERE order_uid = $1`
	err = repo.db.QueryRow(deliveryQuery, id).Scan(&orderData.Delivery.Name, &orderData.Delivery.Phone, &orderData.Delivery.Zip, &orderData.Delivery.City, &orderData.Delivery.Address, &orderData.Delivery.Region, &orderData.Delivery.Email)
	if err != nil {
		return nil, err
	}

	// Получаем данные из таблицы payment
	paymentQuery := `SELECT transaction, request_id, currency, provider, amount, payment_dt, bank, delivery_cost, goods_total, custom_fee FROM payment WHERE order_uid = $1`
	err = repo.db.QueryRow(paymentQuery, id).Scan(&orderData.Payment.Transaction, &orderData.Payment.RequestID, &orderData.Payment.Currency, &orderData.Payment.Provider, &orderData.Payment.Amount, &orderData.Payment.PaymentDT, &orderData.Payment.Bank, &orderData.Payment.DeliveryCost, &orderData.Payment.GoodsTotal, &orderData.Payment.CustomFee)
	if err != nil {
		return nil, err
	}

	// Получаем данные из таблицы items
	itemsQuery := `SELECT chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status FROM items WHERE order_uid = $1`
	rows, err := repo.db.Query(itemsQuery, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []model.Item
	for rows.Next() {
		var item model.Item
		if err := rows.Scan(&item.ChrtID, &item.TrackNumber, &item.Price, &item.RID, &item.Name, &item.Sale, &item.Size, &item.TotalPrice, &item.NMID, &item.Brand, &item.Status); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	orderData.Items = items

	return &orderData, nil
}

func (repo *Repo) GetAllOrders() ([]model.Order, error) {
	// Получаем все заказы
	orderQuery := `
		SELECT order_uid 
		FROM orders`
	rows, err := repo.db.Query(orderQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []model.Order
	for rows.Next() {
		var orderUID string
		if err := rows.Scan(&orderUID); err != nil {
			return nil, err
		}

		// Используем GetOrder для получения полного заказа
		order, err := repo.GetOrder(orderUID)
		if err != nil {
			return nil, err
		}
		orders = append(orders, *order)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return orders, nil
}

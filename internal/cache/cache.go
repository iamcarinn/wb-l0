package cache

import (
	"wb-l0/internal/model"
	"github.com/patrickmn/go-cache"
)

// Структура содержащая кэш
type Cache struct {
	store *cache.Cache
}

// Создает новый экземпляр кэша
func New() *Cache {
	return &Cache{
		store: cache.New(cache.NoExpiration, cache.NoExpiration),
	}
}

// Добавляет заказ в кэш
func (c *Cache) AddToCache(order model.Order) {
	c.store.SetDefault(order.OrderUID, order)
}

// Получает заказ из кэша
func (c *Cache) GetFromCache(orderUID string) (model.Order, bool) {
	val, found := c.store.Get(orderUID)
	if !found {
		return model.Order{}, false
	}
	order, ok := val.(model.Order)
	return order, ok
}

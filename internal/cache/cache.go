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



// TODO: 3)
// OK Когда заново запускается мой сервис, заполняется кэш из бд
// OK это происходит в мейн
// OK добавить метод в репо который будет доставать все и возвращать (слайс моделей)
// OK и в кэш добавить метод который будет весь слайс моделей добавлять в кэш

// липо не в слайс а в мапу
// func NewFrom(defaultExpiration, cleanupInterval time.Duration, items map[string]Item) *Cache


// TODO:
// потом поменять хэндлер чтобы выдавать из кэша а не из базы
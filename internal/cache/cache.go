package cache

import (
	"wb-l0/internal/model"

	"github.com/patrickmn/go-cache"
)

type Cache struct {
	store *cache.Cache
}


func New() *cache.Cache{
	c := cache.New(cache.NoExpiration, cache.NoExpiration)
	return c
}

func (cache *Cache) AddToCache(order model.Order) {
	cache.store.SetDefault(order.OrderUID, order)
}

// TODO: 3)
// Когда заново запускается мой сервис, заполняется кэш из бд
// это происходит в мейн
// добавить метод в репо который будет доставать все и возвращать (слайс моделей)
// и в кэш добавить метод который будет весь слайс моделей добавлять в кэш

// липо не в слайс а в мапу
// func NewFrom(defaultExpiration, cleanupInterval time.Duration, items map[string]Item) *Cache


// TODO:
// потом поменять хэндлер чтобы выдавать из кэша а не из базы
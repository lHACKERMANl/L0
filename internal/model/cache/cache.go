package cache

import (
	"mvcModule/internal/model"
	"sync"
)

type ICache interface {
	Get(key string) (interface{}, bool)
	Set(key string, value interface{})
	Delete(key string)
}

type Cache struct {
	data map[string]interface{}
	mu   sync.RWMutex
}

func NewCache() *Cache {
	return &Cache{
		data: make(map[string]interface{}),
	}
}

func (c *Cache) Get(key string) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	value, ok := c.data[key]
	return value, ok
}

func (c *Cache) Set(key string, value interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.data[key] = value
}

func (c *Cache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.data, key)
}

func (c *Cache) SaveOrder(orderUID string, order model.OrderRepository) {
	c.Set(orderUID, order)
}

func (c *Cache) GetOrder(orderUID string) (model.OrderRepository, bool) {
	order, ok := c.Get(orderUID)
	if !ok {
		return model.OrderRepository{}, false
	}
	return order.(model.OrderRepository), true
}

func (c *Cache) SavePayment(paymentID string, payment model.PaymentRepository) {
	c.Set(paymentID, payment)
}

func (c *Cache) GetPayment(paymentID string) (model.PaymentRepository, bool) {
	payment, ok := c.Get(paymentID)
	if !ok {
		return model.PaymentRepository{}, false
	}
	return payment.(model.PaymentRepository), true
}

func (c *Cache) SaveItemsForOrder(orderUID string, items []model.ItemsRepository) {
	c.Set(orderUID+"_items", items)
}

func (c *Cache) GetItemsForOrder(orderUID string) ([]model.ItemsRepository, bool) {
	items, ok := c.Get(orderUID + "_items")
	if !ok {
		return nil, false
	}
	return items.([]model.ItemsRepository), true
}

func (c *Cache) SaveDelivery(deliveryID string, delivery model.DeliveryRepository) {
	c.Set(deliveryID, delivery)
}

func (c *Cache) GetDelivery(deliveryID string) (model.DeliveryRepository, bool) {
	delivery, ok := c.Get(deliveryID)
	if !ok {
		return model.DeliveryRepository{}, false
	}
	return delivery.(model.DeliveryRepository), true
}

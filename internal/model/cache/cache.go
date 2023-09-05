package cache

import (
	"database/sql"
	"log"
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
	db   *sql.DB
}

func NewCache() (*Cache, error) {
	db, err := sql.Open("postgres", "postgres://root:123@postgres/postgreDB?sslmode=disable")
	if err != nil {
		return nil, err
	}
	log.Printf("Ping establashed %v", db.Ping())

	return &Cache{
		data: make(map[string]interface{}),
		db:   db,
	}, nil
}

func (c *Cache) loadCacheFromDatabase() (map[string]interface{}, error) {
	query := "SELECT key, value FROM cache_table"
	rows, err := c.db.Query(query)
	if err == sql.ErrNoRows {
		log.Printf("Empty data")
	} else if err != nil {
		return nil, err
	}
	defer rows.Close()

	cacheData := make(map[string]interface{})
	for rows.Next() {
		var key string
		var value []byte
		if err := rows.Scan(&key, &value); err != nil {
			return nil, err
		}
		cacheData[key] = value
	}

	return cacheData, nil
}

func (c *Cache) saveCacheToDatabase(data map[string]interface{}) error {
	tx, err := c.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	for key, value := range data {
		_, err := tx.Exec("INSERT INTO cache_table (key, value) VALUES ($1, $2)", key, value)
		if err != nil {
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
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

func (c *Cache) ConstructOrderDetails(orderID string) model.OrderDetails {
	cOrder, ok := c.GetOrder(orderID)
	if ok != true {
		log.Fatalf("GetOrder chache exception")
	}
	cPayment, ok := c.GetPayment(orderID)
	if ok != true {
		log.Fatalf("GetPayment chache exception")
	}
	cItems, ok := c.GetItemsForOrder(orderID)
	if ok != true {
		log.Fatalf("GetItemsForOrder chache exception")
	}
	cDelivery, ok := c.GetDelivery(orderID)
	if ok != true {
		log.Fatalf("GetDelivery chache exception")
	}

	return model.OrderDetails{
		cOrder,
		cPayment,
		cItems,
		cDelivery,
	}
}

func (c *Cache) SaveOrderDetails(orderID string, details model.OrderDetails) {
	c.Set(orderID, details)
}

func (c *Cache) GetOrderDetails(orderID string) (model.OrderDetails, bool) {
	delivery, ok := c.Get(orderID)
	if !ok {
		return model.OrderDetails{}, false
	}
	return delivery.(model.OrderDetails), true
}

func (c *Cache) InitializeCache() error {
	cacheData, err := c.loadCacheFromDatabase()
	if err != nil {
		log.Fatalf("Error select data from db: %v", err)
		return err
	}

	if cacheData != nil {
		c.mu.Lock()
		defer c.mu.Unlock()
		c.data = cacheData
		log.Println("Cache loaded from database.")
	} else {
		if err := c.saveCacheToDatabase(c.data); err != nil {
			log.Fatalf("Error saving cache to database: %v", err)
		}
		log.Println("Cache created and saved to database.")
	}

	return nil
}

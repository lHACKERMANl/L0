package messaging

import (
	"fmt"
	_ "github.com/nats-io/nats.go"
	"github.com/nats-io/stan.go"
	"log"
	"mvcModule/internal/config"
	"mvcModule/internal/model"
	"mvcModule/internal/model/database"
	"reflect"
	"sync"
	"time"
)

type IMessaging interface {
	ConnectToMessaging() error
	Subscribe() error
}

type OrderDetails struct {
	OrderID  string
	Order    model.OrderRepository
	Payment  model.PaymentRepository
	Items    []model.ItemsRepository
	Delivery model.DeliveryRepository
}

type Messaging struct{}

func (*Messaging) Subscribe() error {
	panic("unimplemented")
}

type NATSMessaging struct {
	natsConn stan.Conn
	database database.IDatabase
}

func (m *NATSMessaging) ConnectToMessaging() error {

	url := fmt.Sprintf("nats://172.30.0.40:4222")
	sc, err := stan.Connect(
		"test-cluster",
		"worker",
		stan.NatsURL(url))
	if err != nil {
		log.Fatalf("Error in NATS connection: %v", err)
		return err
	}
	log.Printf("Connection with nats: %v estableshed", url)
	m.natsConn = sc
	return nil
}

func (m *NATSMessaging) NATSub() (OrderDetails, error) {
	var orderDetails OrderDetails

	cache := make(map[string]string)
	var cacheMutex sync.Mutex
	var orderID string

	terminateChannel := make(chan struct{})

	_, err := m.natsConn.Subscribe("getFromDB", func(msg *stan.Msg) {
		orderID = string(msg.Data)
		log.Printf("orderID: %v type: %v", orderID, reflect.TypeOf(orderID))

		err := msg.Ack()
		if err != nil {
			log.Printf("Error acknowledging message: %v", err)
		}

		conf, err := config.GetDataFromDockerCompose("config.yaml")
		if err != nil {
			log.Fatalln("no data from config")
		}

		err = m.database.ConnectToDB(conf)
		if err != nil {
			return
		}

		err = m.database.Ping()
		if err != nil {
			log.Printf("No DB connection: %v", err)
		} else {
			log.Printf("DB Ping status: %v")
		}

		order, err := m.database.GetOrderDetails(orderID)
		if err != nil {
			log.Printf("Error getting order details from messaging")
		}

		payment, err := m.database.GetPaymentDetails(orderID)
		if err != nil {
			log.Printf("Error getting payment details from messaging")
		}

		items, err := m.database.GetItemsForOrder(order.TrackNumber)
		if err != nil {
			log.Printf("Error getting items details from messaging")
		}

		delivery, err := m.database.GetDeliveryDetails(orderID)
		if err != nil {
			log.Printf("Error getting delivery details from messaging")
		}

		orderDetails = OrderDetails{
			Order:    order,
			Payment:  payment,
			Items:    items,
			Delivery: delivery,
		}

		//cache
		cacheMutex.Lock()
		cache["getFromDB"] = orderID
		cacheMutex.Unlock()

		log.Printf("OrderID `%s` got and wrote:\n", orderID)

		go func() {
			TTLcache := 10 * time.Minute
			ticker := time.NewTicker(TTLcache)
			defer ticker.Stop()

			for {
				select {
				case <-ticker.C:
					cacheMutex.Lock()
					for key := range cache {
						delete(cache, key)
					}
					cacheMutex.Unlock()
				case <-terminateChannel:
					return
				}
			}
		}()

	})
	if err != nil {
		log.Fatalf("Error with NATS connection: %v", err)
	}

	for {
		select {
		case <-terminateChannel:
			return orderDetails, nil
		}
	}
}

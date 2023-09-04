package messaging

import (
	"encoding/json"
	_ "github.com/nats-io/nats.go"
	"github.com/nats-io/stan.go"
	"log"
	"mvcModule/internal/config"
	"mvcModule/internal/model"
	"mvcModule/internal/model/cache"
	"mvcModule/internal/model/database"

	"github.com/mitchellh/mapstructure"
)

type IMessaging interface {
	ConnectToMessaging(MessagingConf config.LoginData) error
	Subscribe() error
}

type Messaging struct{}

// Subscribe implements IMessaging.
func (*Messaging) Subscribe() error {
	panic("unimplemented")
}

func (*Messaging) ConnectToMessaging(MessagingConf config.LoginData) error {
	panic("unimplemented")
}

type NATSMessaging struct {
	natsConn stan.Conn
	database database.IDatabase
}

func NatsConnect(db database.IDatabase, nats stan.Conn) *NATSMessaging {
	return &NATSMessaging{
		natsConn: nats,
		database: db,
	}
}

func (m *NATSMessaging) ConnectToMessaging(MessagingConf config.LoginData) error {
	var err error
	sc, err := stan.Connect(
		"test-cluster",
		"my-client",
		stan.NatsURL("nats://172.30.0.3:4222"))
	if err != nil {
		log.Fatalf("Error in NATS connection: %v", err)
		return err
	}
	msg := []byte("Hello")
	err = sc.Publish("subj", msg)
	if err != nil {
		log.Fatalf("Ошибка при отправке сообщения: %v", err)
	}
	return nil
}

func (m *NATSMessaging) Subscribe() error {
	sub, err := m.natsConn.Subscribe("saveToDB", func(msg *stan.Msg) {
		recivedMsg := msg.Data
		order, payment, items, delivery := handleMsg(recivedMsg)

		cache := cache.NewCache()

		cache.SaveOrder(order.OrderUID, order)
		cache.SavePayment(payment.PaymentDT, payment)
		cache.SaveItemsForOrder(order.OrderUID, items)
		cache.SaveDelivery(delivery.OrderID, delivery)

		m.database.SaveOrderRepositoryToDB(order)
		m.database.SavePaymentRepositoryToDB(payment)
		for _, item := range items {
			m.database.SaveItemsRepositoryToDB(item)
		}
		m.database.SaveDeliveryRepositoryToDB(delivery, order.OrderUID)
	}, stan.DurableName("saveData"))
	if err != nil {
		log.Fatalln("Subscription error: ", err)
		return err
	}
	defer sub.Unsubscribe()
	return nil
}

func handleMsg(msg []byte) (model.OrderRepository, model.PaymentRepository, []model.ItemsRepository, model.DeliveryRepository) {
	var data map[string]interface{}
	if err := json.Unmarshal(msg, &data); err != nil {
		log.Fatalln("Error decoding JSON:", err)
	}

	var order model.OrderRepository
	var payment model.PaymentRepository
	var items []model.ItemsRepository
	var delivery model.DeliveryRepository

	if err := mapstructure.Decode(data["order"], &order); err != nil {
		log.Fatalln("Error decoding order: ", err)
	}

	if err := mapstructure.Decode(data["payment"], &payment); err != nil {
		log.Fatalln("Error payment order: ", err)
	}

	if err := mapstructure.Decode(data["items"], &items); err != nil {
		log.Fatalln("Error items order: ", err)
	}

	if err := mapstructure.Decode(data["delivery"], &delivery); err != nil {
		log.Fatalln("Error delivery order: ", err)
	}

	return order, payment, items, delivery
}

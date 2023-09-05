package presener

import (
	"fmt"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"github.com/nats-io/stan.go"
	"log"
	"mvcModule/internal/config"
	"mvcModule/internal/model"
	"mvcModule/internal/model/cache"
	"mvcModule/internal/model/database"
	"mvcModule/internal/model/messaging"
	"mvcModule/view"
	"net/http"
)

type OrderDetails struct {
	OrderID  string
	Order    model.OrderRepository
	Payment  model.PaymentRepository
	Items    []model.ItemsRepository
	Delivery model.DeliveryRepository
}

type OrderPresenter interface {
	GetOrderDetails(orderID string) (OrderDetails, error)
}

func (p *Presenter) orderHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		orderID := r.URL.Query().Get("order_id")

		if orderID == "" {
			err := p.handleOrder(w)
			if err != nil {
				log.Printf("Warning in getting data from handleOrder() func: %v", err)
			}
		} else if orderID != "" {
			user := uuid.New().String()
			connect, err := stan.Connect(
				"test-cluster",
				user,
				stan.NatsURL("nats://172.30.0.40:4222"))
			if err != nil {
				log.Printf("no connection with nats: %v", err)
			}

			err = connect.Publish("getFromDB", []byte(orderID))
			if err != nil {
				return
			} else {
				log.Printf("orderID: %v was published", orderID)
			}
			defer connect.Close() //todo check defer

			err = p.handleOrderDetails(w, orderID)
			if err != nil {
				log.Printf("Warning in getting data from handleOrderDetails() func: %v", err)
			}
			//p.model.natsConn.NATSpub(orderID)
		}
	})
}

func (p *Presenter) handleOrder(w http.ResponseWriter) error {
	html := p.view.ShowOrder()

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprintf(w, html)
	return nil
}

func (p *Presenter) handleOrderDetails(w http.ResponseWriter, orderID string) error {
	orderDetails, err := p.model.GetOrderDetails(orderID)
	if err != nil {
		http.Error(w, "Error getting order details", http.StatusInternalServerError)
		return err
	}

	od := OrderDetails{
		OrderID:  orderID,
		Order:    orderDetails.Order,
		Payment:  orderDetails.Payment,
		Items:    orderDetails.Items,
		Delivery: orderDetails.Delivery,
	}

	html := p.view.ShowOrderDetails(View.OrderDetails(od))

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprintf(w, html)

	return nil
}

func (p *Presenter) SetupHandlers() {
	p.ViewHandler = p.orderHandler()
	http.Handle("/order", p.ViewHandler)
}

type Model struct {
	db       database.Database
	natsConn messaging.NATSMessaging
	view     View.View
}

func NewModel(db database.Database, natsConn messaging.NATSMessaging) *Model {
	return &Model{
		db:       db,
		natsConn: natsConn,
	}
}

type Presenter struct {
	model       Model
	view        View.View
	ViewHandler http.Handler
}

func NewPresenter(model Model, view View.View) *Presenter {
	presenter := &Presenter{
		model: model,
		view:  view,
	}
	return presenter
}

func (p *Model) GetOrderDetails(orderID string) (OrderDetails, error) {
	cache, err := cache.NewCache()
	if err != nil {
		log.Printf("NewCache() exception: %v", err)
	}
	cache.InitializeCache()
	var orderDetails OrderDetails
	var orderDetailsCached model.OrderDetails

	_, ok := cache.GetOrder(orderID)
	if ok {
		orderDetailsCached, ok = cache.GetOrderDetails(orderID)
		orderDetails = OrderDetails{
			orderID,
			orderDetailsCached.Order,
			orderDetailsCached.Payment,
			orderDetailsCached.Items,
			orderDetailsCached.Delivery,
		}
		if ok != true {
			log.Printf("cahce exception")
		}
	} else {
		order, err := p.db.GetOrderDetails(orderID)
		if err != nil {
			return OrderDetails{}, err
		}

		if err != nil && p.view != nil {
			p.view.ShowError("Error getting payment details from messaging")
		}

		payment, err := p.db.GetPaymentDetails(orderID)
		if err != nil && p.view != nil {
			p.view.ShowError("Error getting payment details from messaging")
		}

		items, err := p.db.GetItemsForOrder(order.TrackNumber)
		if err != nil && p.view != nil {
			p.view.ShowError("Error getting items details from messaging")
		}

		delivery, err := p.db.GetDeliveryDetails(orderID)
		if err != nil && p.view != nil {
			p.view.ShowError("Error getting delivery details from messaging")
		}

		orderDetails = OrderDetails{
			Order:    order,
			Payment:  payment,
			Items:    items,
			Delivery: delivery,
		}

		cache.SaveOrderDetails(orderID, orderDetailsCached)
	}

	return orderDetails, nil
}

func Init() Presenter {
	conf, err := config.GetDataFromDockerCompose("config.yaml")
	if err != nil {
		log.Fatalln("no data from config")
	}
	db := database.NewDatabase(nil)
	err = db.ConnectToDB(conf)
	if err != nil {
		return Presenter{}
	}
	var nats messaging.NATSMessaging
	err = nats.ConnectToMessaging()
	if err != nil {
		return Presenter{}
	}
	newModel := NewModel(*db, nats)
	view := View.NewFakeView()
	presenter := NewPresenter(*newModel, view)
	presenter.SetupHandlers()
	presenter.orderHandler()
	return *presenter
}

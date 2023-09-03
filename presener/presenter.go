package presener

import (
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"mvcModule/internal/config"
	"mvcModule/internal/model"
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
		// Извлечь параметры запроса
		orderID := r.URL.Query().Get("order_id")

		// Вызвать метод презентера для получения деталей заказа
		orderDetails, err := p.model.GetOrderDetails(orderID)
		if err != nil {
			http.Error(w, "Error getting order details", http.StatusInternalServerError)
			return
		}

		od := OrderDetails{
			OrderID:  orderID,
			Order:    orderDetails.Order,
			Payment:  orderDetails.Payment,
			Items:    orderDetails.Items,
			Delivery: orderDetails.Delivery}

		// Использовать view для отображения деталей заказа
		html := p.view.ShowOrderDetails(View.OrderDetails(od))

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprintf(w, html)
	})
}

func (p *Presenter) SetupHandlers() {
	// Создать обработчик для пути "/order" и связать его
	p.ViewHandler = p.orderHandler()
	http.Handle("/order", p.ViewHandler)
}

type Model struct {
	db       database.Database
	natsConn messaging.NATSMessaging
	view     View.View
	// database  database.IDatabase
	// messaging messaging.IMessaging
}

func NewModel(db database.Database, natsConn messaging.NATSMessaging) *Model {
	return &Model{
		db:       db,
		natsConn: natsConn,
	}
}

type Presenter struct {
	model       OrderPresenter
	view        View.View
	ViewHandler http.Handler
}

func NewPresenter(model OrderPresenter, view View.View) *Presenter {
	presenter := &Presenter{
		model: model,
		view:  view,
	}
	return presenter
}

func (p *Model) GetOrderDetails(orderID string) (OrderDetails, error) {
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

	orderDetails := OrderDetails{
		Order:    order,
		Payment:  payment,
		Items:    items,
		Delivery: delivery,
	}

	return orderDetails, nil
}

func Init() Presenter {
	// absFilePath, err := filepath.Abs("/app/config.yaml")
	// if err != nil {
	// 	log.Fatalln("no config.yaml in project")
	// }
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
	err = nats.ConnectToMessaging(conf)
	if err != nil {
		return Presenter{}
	}
	//natsMessaging := messaging.NatsConnect(db, nats)
	//err = natsMessaging.ConnectToMessaging(conf)
	if err != nil {
		return Presenter{}
	}
	newModel := NewModel(*db, nats)
	view := View.NewFakeView()
	presenter := NewPresenter(newModel, view)
	presenter.SetupHandlers()
	presenter.orderHandler()
	//view.ShowOrderDetails("b563feb7b2b84b6test")
	return *presenter
}

///order_id?order_id=b563feb7b2b84b6test

package View

import (
	"html/template"
	"log"
	"mvcModule/internal/model"
	"net/http"
)

type View interface {
	ShowOrderDetails(orderDetails OrderDetails) string
	ShowError(msg string)
}

type OrderDetails struct {
	OrderID  string
	Order    model.OrderRepository
	Payment  model.PaymentRepository
	Items    []model.ItemsRepository
	Delivery model.DeliveryRepository
}

type ViewHandler struct {
	view View
}

func NewViewHandler(view View) *ViewHandler {
	return &ViewHandler{
		view: view,
	}
}

func (h *ViewHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	orderID := r.URL.Query().Get("order_id")
	order := OrderDetails{
		OrderID:  orderID,
		Order:    model.OrderRepository{},
		Payment:  model.PaymentRepository{},
		Items:    []model.ItemsRepository{},
		Delivery: model.DeliveryRepository{},
	}
	orderDetails := h.view.ShowOrderDetails(order) // Получаем детали заказа

	tmpl, err := template.ParseFiles("static/index.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	err = tmpl.Execute(w, orderDetails)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

type FakeView struct{}

func NewFakeView() *FakeView {
	return &FakeView{}
}

func (v *FakeView) ShowOrderDetails(orderDetails OrderDetails) string {
	//tmpl, err := template.ParseFiles("static/index.html")
	//if err != nil {
	//	http.Error(w, err.Error(), http.StatusInternalServerError)
	//}
	//
	//err = tmpl.Execute(w, orderDetails)
	//if err != nil {
	//	http.Error(w, err.Error(), http.StatusInternalServerError)
	//}

	html := OrderHTMLView(orderDetails)
	return html
}

func (v *FakeView) ShowError(msg string) {
	log.Fatalln("Error in view: ", msg)
}

func initView() {
	view := NewFakeView()
	viewHandler := NewViewHandler(view)

	http.Handle("/order", viewHandler)
	http.ListenAndServe(":8080", nil)
}

package model

type OrderRepository struct {
	OrderUID          string `json:"order_uid"`
	TrackNumber       string `json:"track_number"`
	Entry             string `json:"entry"`
	Locale            string `json:"locale"`
	InternalSignature string `json:"internal_signature"`
	CustomerId        string `json:"customer_id"`
	DeliveryService   string `json:"delivery_service"`
	Shardkey          string `json:"shardkey"`
	SmID              string `json:"sm_id"`
	DataCreated       string `json:"data_created"`
	OofShard          string `json:"oof_shard"`
}

type PaymentRepository struct {
	Transaction  string `json:"transaction"`
	RequestID    string `json:"request_id"`
	Currency     string `json:"currency"`
	Provider     string `json:"provider"`
	Amount       string `json:"amount"`
	PaymentDT    string `json:"payment_dt"`
	Bank         string `json:"bank"`
	DeliveryCost string `json:"delivery_cost"`
	GoodsTotal   string `json:"goods_total"`
	CustomFee    string `json:"custom_fee"`
}

type ItemsRepository struct {
	ChrtID      string `json:"chrt_id"`
	TrackNumber string `json:"track_number"`
	Price       string `json:"price"`
	Rid         string `json:"rid"`
	Name        string `json:"name"`
	Sale        string `json:"sale"`
	Size        string `json:"size"`
	TotalPrice  string `json:"total_price"`
	NmID        string `json:"nm_id"`
	Brand       string `json:"brand"`
	Status      string `json:"status"`
	OrderID     string `json:"order_uid"`
}

type DeliveryRepository struct {
	Name    string `json:"name"`
	Phone   string `json:"phone"`
	Zip     string `json:"zip"`
	City    string `json:"city"`
	Address string `json:"address"`
	Region  string `json:"region"`
	Email   string `json:"email"`
	OrderID string `json:"order_uid"`
}

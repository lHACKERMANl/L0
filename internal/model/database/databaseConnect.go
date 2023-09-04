package database

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	"mvcModule/internal/config"
	"mvcModule/internal/model"

	_ "github.com/lib/pq"
)

type IDatabase interface {
	ConnectToDB(SQLconf config.LoginData) error
	SaveOrderRepositoryToDB(dataOrderRepository model.OrderRepository)
	SavePaymentRepositoryToDB(dataPaymentRepository model.PaymentRepository)
	SaveItemsRepositoryToDB(dataItemsRepository model.ItemsRepository)
	SaveDeliveryRepositoryToDB(dataDeliveryRepository model.DeliveryRepository, orderID string) // Добавлен параметр orderID
	GetOrderDetails(orderID string) (model.OrderRepository, error)
	GetPaymentDetails(transactionID string) (model.PaymentRepository, error)
	GetItemsForOrder(trackNumber string) ([]model.ItemsRepository, error)
	GetDeliveryDetails(orderUID string) (model.DeliveryRepository, error)
}

type Database struct {
	db *sql.DB
}

func NewDatabase(db *sql.DB) *Database {
	return &Database{
		db: db,
	}
}

func (d *Database) GetDB() *sql.DB {
	return d.db
}

func (d *Database) SetDB(newDB *sql.DB) {
	d.db = newDB
}

func (d *Database) ConnectToDB(PSQLconf config.LoginData) error {
	env := PSQLconf.Services.Postgres.Environment
	service := PSQLconf.Services.Postgres
	name := strings.Split(service.Image, ":")[0]
	var err error
	psqlInfo := fmt.Sprintf("postgres://%s:%s@postgres/%s?sslmode=disable",
		env.PostgresUser, env.PostgresPassword, env.PostgresDB)
	d.db, err = sql.Open(name, psqlInfo)
	if err != nil {
		log.Fatalln("Error opening DB: ", err)
	}
	err = d.db.Ping()
	if err != nil {
		log.Fatalf("Error to connect database: %v", err)
	}
	log.Printf("postgres://%s:%s@postgres/%s?sslmode=disable",
		env.PostgresUser, env.PostgresPassword, env.PostgresDB)
	//defer d.db.Close()

	//row := d.db.QueryRow(`SELECT * FROM orders WHERE "order_uid" = $1;`, "b563feb7b2b84b6test")
	//
	//var order model.OrderRepository
	//
	//err = row.Scan(
	//	&order.OrderUID,
	//	&order.TrackNumber,
	//	&order.Entry,
	//	&order.Locale,
	//	&order.InternalSignature,
	//	&order.CustomerId,
	//	&order.DeliveryService,
	//	&order.Shardkey,
	//	&order.SmID,
	//	&order.DataCreated,
	//	&order.OofShard,
	//)
	//
	//if err != nil {
	//	log.Printf("Error in CreateDB: %v", err)
	//}

	return nil
}

func (d *Database) SaveOrderRepositoryToDB(dataOrderRepository model.OrderRepository) {
	err := d.db.QueryRow(`INSERT INTO orders (
		order_uid, 
		track_number, 
		entry, 
		locale, 
		internal_signature,
		customer_id,
		delivery_service,
		shardkey,
		sm_id,
		date_created,
		oof_shard) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`,
		dataOrderRepository.OrderUID,
		dataOrderRepository.TrackNumber,
		dataOrderRepository.Entry,
		dataOrderRepository.Locale,
		dataOrderRepository.InternalSignature,
		dataOrderRepository.CustomerId,
		dataOrderRepository.DeliveryService,
		dataOrderRepository.Shardkey,
		dataOrderRepository.SmID,
		dataOrderRepository.DataCreated,
		dataOrderRepository.OofShard)

	if err != nil {
		log.Fatalln("Error with save Order Repository to DB: ", err)
	}
}

func (d *Database) SavePaymentRepositoryToDB(dataPaymentRepository model.PaymentRepository) {
	err := d.db.QueryRow(`INSERT INTO payment (
		transaction,
		request_id,
		currency,
		provider,
		amount,
		payment_dt,
		bank,
		delivery_cost,
		goods_total,
		custom_fee) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`,
		dataPaymentRepository.Transaction,
		dataPaymentRepository.RequestID,
		dataPaymentRepository.Currency,
		dataPaymentRepository.Provider,
		dataPaymentRepository.Amount,
		dataPaymentRepository.PaymentDT,
		dataPaymentRepository.Bank,
		dataPaymentRepository.DeliveryCost,
		dataPaymentRepository.GoodsTotal,
		dataPaymentRepository.CustomFee)

	if err != nil {
		log.Fatalln("Error with save Payment Repository to DB: ", err)
	}
}

func (d *Database) SaveItemsRepositoryToDB(dataItemsRepository model.ItemsRepository) {
	err := d.db.QueryRow(`INSERT INTO items (
		"chrt_id",
		"track_number",
		"price",
		"rid",
		"name",
		"sale",
		"size",
		"total_price",
		"nm_id",
		"brand",
		"status") VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`,
		dataItemsRepository.ChrtID,
		dataItemsRepository.TrackNumber,
		dataItemsRepository.Price,
		dataItemsRepository.Rid,
		dataItemsRepository.Name,
		dataItemsRepository.Sale,
		dataItemsRepository.Size,
		dataItemsRepository.TotalPrice,
		dataItemsRepository.NmID,
		dataItemsRepository.Brand,
		dataItemsRepository.Status)

	if err != nil {
		log.Fatalln("Error with save Items Repository to DB: ", err)
	}
}

func (d *Database) SaveDeliveryRepositoryToDB(dataDeliveryRepository model.DeliveryRepository, orderID string) {
	err := d.db.QueryRow(`INSERT INTO delivery (
		"name",
		"phone",
		"zip",
		"city",
		"address",
		"region",
		"email",
		"order_uid") VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		dataDeliveryRepository.Name,
		dataDeliveryRepository.Phone,
		dataDeliveryRepository.Zip,
		dataDeliveryRepository.City,
		dataDeliveryRepository.Address,
		dataDeliveryRepository.Region,
		dataDeliveryRepository.Email,
		orderID)

	if err != nil {
		log.Fatalln("Error with save Items Repository to DB: ", err)
	}
}

func (d *Database) GetOrderDetails(orderID string) (model.OrderRepository, error) {
	var order model.OrderRepository

	row := d.db.QueryRow(`SELECT * FROM orders WHERE "order_uid" = $1;`, orderID)

	var internalSignature sql.NullString
	err := row.Scan(
		&order.OrderUID,
		&order.TrackNumber,
		&order.Entry,
		&order.DataCreated,
		&order.CustomerId,
		&order.DeliveryService,
		&order.Shardkey,
		&order.SmID,
		&order.OofShard,
		&internalSignature,
		&order.Locale,
	)

	if err == sql.ErrNoRows {
		return model.OrderRepository{}, err
	}

	if internalSignature.Valid {
		order.InternalSignature = internalSignature.String
	} else {
		order.InternalSignature = ""
	}

	return order, nil
}

func (d *Database) GetPaymentDetails(transactionID string) (model.PaymentRepository, error) {
	var payment model.PaymentRepository

	row := d.db.QueryRow(`SELECT * FROM payment WHERE "transaction" = $1;`, transactionID)

	err := row.Scan(
		&payment.Transaction,
		&payment.RequestID,
		&payment.Currency,
		&payment.Provider,
		&payment.Amount,
		&payment.PaymentDT,
		&payment.Bank,
		&payment.DeliveryCost,
		&payment.GoodsTotal,
		&payment.CustomFee,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return model.PaymentRepository{}, err
		}
		return model.PaymentRepository{}, err
	}

	return payment, nil
}

func (d *Database) GetItemsForOrder(trackNumber string) ([]model.ItemsRepository, error) {
	var items []model.ItemsRepository

	rows, err := d.db.Query(` SELECT * FROM items WHERE "track_number" = $1;`, trackNumber)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var item model.ItemsRepository
		err := rows.Scan(
			&item.ChrtID,
			&item.OrderID,
			&item.TrackNumber,
			&item.Price,
			&item.Rid,
			&item.Name,
			&item.Sale,
			&item.Size,
			&item.TotalPrice,
			&item.NmID,
			&item.Brand,
			&item.Status,
		)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}

func (d *Database) GetDeliveryDetails(orderUID string) (model.DeliveryRepository, error) {
	var delivery model.DeliveryRepository

	err := d.db.QueryRow(`SELECT * FROM delivery WHERE "order_uid" = $1;`, orderUID).Scan(
		&delivery.OrderID,
		&delivery.Name,
		&delivery.Phone,
		&delivery.City,
		&delivery.Zip,
		&delivery.Address,
		&delivery.Region,
		&delivery.Email,
	)
	if err != nil {
		return model.DeliveryRepository{}, err
	}

	return delivery, nil
}

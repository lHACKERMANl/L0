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
	Ping() error
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
	db, err := sql.Open(name, psqlInfo)
	if err != nil {
		log.Fatalln("Error opening DB: ", err)
	}
	err = db.Ping()
	if err != nil {
		log.Fatalf("Error to connect database: %v", err)
	}
	log.Printf("postgres://%s:%s@postgres/%s?sslmode=disable",
		env.PostgresUser, env.PostgresPassword, env.PostgresDB)
	//defer d.db.Close()

	d.db = db
	return nil
}

func (d *Database) Ping() error {
	err := d.db.Ping()
	if err != nil {
		return err
	}
	return nil
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

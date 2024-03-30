package db

import (
	"L0/app/internal/models"
	"context"
	"log"

	"github.com/jackc/pgx/v5"
)

type Database interface {
	Create(ctx context.Context, m models.Model) error
	Get(ctx context.Context, id string) (models.Model, error)
	GetAll(ctx context.Context) ([]models.Model, error)
}

type DB interface {
	GetClient() (conn *pgx.Conn, err error)
}

type database struct {
	conn DB
}

func NewDB(conn DB) Database {
	return &database{conn: conn}
}

func (db *database) Create(ctx context.Context, m models.Model) error {
	conn, err := db.conn.GetClient()
	if err != nil {
		log.Println(err)
		return err
	}
	tx, err := conn.Begin(ctx)
	if err != nil {
		return err
	}
	q1 := `
		INSERT INTO model 
			(order_uid, track_number, entry,
			locale, internal_signature, customer_id,
			delivery_service, shardkey, sm_id,
			date_created, oof_shard)
		VALUES
			($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`
	q2 := `
		INSERT INTO delivery
			(name, phone, zip, city, address, region, email, fk_order_uid)
		VALUES
			($1, $2, $3, $4, $5, $6, $7, $8)
	`
	q3 := `
		INSERT INTO payment
			(transaction, request_id, currency,
			provider, amount, payment_dt,
			bank, delivery_cost, goods_total, custom_fee)
		VALUES
			($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`
	q4 := `
		INSERT INTO items
			(chrt_id, track_number, price,
			rid, name, sale, size,
			total_price, nm_id, brand, status)
		VALUES
			($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`
	_, err = tx.Exec(ctx, q1,
		m.OrderUid, m.Tracknumber, m.Entry, m.Locale,
		m.Internalsignature, m.Customerid, m.Deliveryservice,
		m.Shardkey, m.Smid, m.Datecreated, m.Oofshard)
	if err != nil {
		tx.Rollback(ctx)
		return err
	}
	_, err = tx.Exec(ctx, q2,
		m.Delivery.Name, m.Delivery.Phone, m.Delivery.Zip,
		m.Delivery.City, m.Delivery.Address, m.Delivery.Region,
		m.Delivery.Email, m.OrderUid)
	if err != nil {
		tx.Rollback(ctx)
		return err
	}
	_, err = tx.Exec(ctx, q3,
		m.Payment.Transaction, m.Payment.Requestid,
		m.Payment.Currency, m.Payment.Provider, m.Payment.Amount,
		m.Payment.Paymentdt, m.Payment.Bank, m.Payment.Deliverycost,
		m.Payment.Goodstotal, m.Payment.Customfee)
	if err != nil {
		tx.Rollback(ctx)
		return err
	}
	for _, item := range m.Items {
		_, err = tx.Exec(ctx, q4,
			item.Chrtid, item.Tracknumber, item.Price, item.Rid, item.Name, item.Sale,
			item.Size, item.Totalprice, item.Nmid, item.Brand, item.Status)
		if err != nil {
			tx.Rollback(ctx)
			return err
		}
	}
	err = tx.Commit(ctx)
	return err
}

func (db *database) Get(ctx context.Context, id string) (m models.Model, err error) {
	conn, err := db.conn.GetClient()
	if err != nil {
		log.Println(err)
		return
	}

	tx, err := conn.Begin(ctx)
	if err != nil {
		return models.Model{}, err
	}

	q1 := `
		SELECT 
			model.order_uid, model.track_number, model.entry, model.locale,
			model.internal_signature, model.customer_id, model.delivery_service,
			model.shardkey, model.sm_id, model.date_created, model.oof_shard,
			delivery.name, delivery.phone, delivery.zip, delivery.city,
			delivery.address, delivery.region, delivery.email,
			payment.transaction, payment.request_id, payment.currency,
			payment.provider, payment.amount, payment.payment_dt, payment.bank,
			payment.delivery_cost, payment.goods_total, payment.custom_fee
		FROM 
			model, delivery, payment, items
		WHERE order_uid=$1
	`
	q2 := `
		SELECT
			chrt_id, track_number, price, rid, name,
			sale, size, total_price, nm_id, brand, status
		FROM items
		WHERE track_number=$1
	`
	row := conn.QueryRow(ctx, q1, id)
	err = row.Scan(
		&m.OrderUid, &m.Tracknumber, &m.Entry, &m.Locale,
		&m.Internalsignature, &m.Customerid, &m.Deliveryservice,
		&m.Shardkey, &m.Smid, &m.Datecreated, &m.Oofshard,
		&m.Delivery.Name, &m.Delivery.Phone, &m.Delivery.Zip,
		&m.Delivery.City, &m.Delivery.Address, &m.Delivery.Region,
		&m.Delivery.Email, &m.Payment.Transaction, &m.Payment.Requestid,
		&m.Payment.Currency, &m.Payment.Provider, &m.Payment.Amount,
		&m.Payment.Paymentdt, &m.Payment.Bank, &m.Payment.Deliverycost,
		&m.Payment.Goodstotal, &m.Payment.Customfee)
	if err != nil {
		log.Println(err)
		tx.Rollback(ctx)
		return models.Model{}, err
	}
	var items = []models.Item{}
	rows, err := conn.Query(ctx, q2, m.Tracknumber)
	if err != nil {
		log.Println(err)
		tx.Rollback(ctx)
		return models.Model{}, err
	}
	defer rows.Close()
	for rows.Next() {
		item := models.Item{}
		err = rows.Scan(
			&item.Chrtid, &item.Tracknumber, &item.Price, &item.Rid, &item.Name, &item.Sale,
			&item.Size, &item.Totalprice, &item.Nmid, &item.Brand, &item.Status)
		if err != nil {
			log.Println(err)
			continue
		}
		items = append(items, item)
	}
	m.Items = items
	err = tx.Commit(ctx)
	return m, err
}

func (db *database) GetAll(ctx context.Context) ([]models.Model, error) {
	var out []models.Model
	var model models.Model
	conn, err := db.conn.GetClient()
	if err != nil {
		log.Println(err)
		return []models.Model{}, err
	}

	rows, err := conn.Query(ctx, "SELECT order_uid FROM model")
	if err != nil {
		log.Println(err)
		return []models.Model{}, err
	}
	defer rows.Close()

	var ids []string
	var id string
	for rows.Next() {
		err = rows.Scan(&id)
		if err != nil {
			log.Println(err)
			continue
		}
		ids = append(ids, id)
	}

	for _, v := range ids {
		model, err = db.Get(ctx, v)
		if err != nil {
			log.Println(err)
			return []models.Model{}, err
		}
		out = append(out, model)
	}
	return out, nil
}

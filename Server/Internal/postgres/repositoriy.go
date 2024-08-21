package postgres

import (
	"context"

	"github.com/ViPDanger/L0/Server/Internal/structures"
	"github.com/nats-io/nats.go"
)

// Описание репозитория клиента
type Repository struct {
	client Client
}

func NewRepository(client Client, nats *nats.Conn) *Repository {
	return &Repository{
		client: client,
	}
}

// Функции репозитория

// Внесения заказа в PG
func (r *Repository) PutOrder(order structures.Order) error {
	ctx := context.Background()
	tx, err := r.client.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	// таблица orders
	query := "insert into wb_schema.orders (order_uid,track_number,entry,locale,internal_signature,customer_id,delivery_service,shardkey,sm_id,date_created,oof_shard) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)"
	_, err = tx.Exec(ctx, query, order.Order_uid, order.Track_number, order.Entry, order.Locale, order.Internal_signature, order.Customer_id, order.Delivery_service, order.Shardkey, order.Sm_id, order.Date_created, order.Oof_shard)
	if err != nil {
		return err
	}
	// таблица delivery
	query = "insert into wb_schema.delivery (order_uid,name,phone,zip,city,address,region,email) VALUES ($1,$2,$3,$4,$5,$6,$7,$8)"
	_, err = tx.Exec(ctx, query, order.Order_uid, order.Delivery.Name, order.Delivery.Phone, order.Delivery.Zip, order.Delivery.City, order.Delivery.Address, order.Delivery.Region, order.Delivery.Email)
	if err != nil {
		return err
	}
	// таблица payment
	query = "insert into wb_schema.payment (order_uid,request_id,currency,provider,amount,payment_dt,bank,delivery_cost,goods_total,custom_fee) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)"
	_, err = tx.Exec(ctx, query, order.Order_uid, order.Payment.Request_id, order.Payment.Currency, order.Payment.Provider, order.Payment.Amount, order.Payment.Payment_dt, order.Payment.Bank, order.Payment.Delivery_cost, order.Payment.Goods_total, order.Payment.Custom_fee)
	if err != nil {
		return err
	}
	// таблица items
	query = "insert into wb_schema.items (chrt_id,track_number,price,rid,name,sale,size,total_price,nm_id,brand,status) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)"
	for _, item := range order.Items {
		_, err = tx.Exec(ctx, query, item.Chrt_id, item.Track_number, item.Price, item.Rid, item.Name, item.Sale, item.Size, item.Total_price, item.Nm_id, item.Brand, item.Status)
		if err != nil {
			return err
		}
		_, err = tx.Exec(ctx, "insert into wb_schema.order_item (order_uid,chrt_id) VALUES ($1,$2)", order.Order_uid, item.Chrt_id)
		if err != nil {
			return err
		}
	}
	tx.Commit(ctx)
	return err
}

// Взятие заказа из PG
func (r *Repository) GetOrder(order_uid string) (structures.Order, error) {
	var order structures.Order
	// таблица orders
	query := "SELECT order_uid,track_number,entry,locale,internal_signature,customer_id,delivery_service,shardkey,sm_id,CONCAT(TO_CHAR(date_created,'YYYY-MM-DDT'),TO_CHAR(date_created,'HH24:MM:SSZ')) as date_created,oof_shard FROM wb_schema.orders WHERE order_uid = $1"
	err := r.client.QueryRow(context.Background(), query, order_uid).Scan(&order.Order_uid, &order.Track_number, &order.Entry, &order.Locale, &order.Internal_signature, &order.Customer_id, &order.Delivery_service, &order.Shardkey, &order.Sm_id, &order.Date_created, &order.Oof_shard)
	if err != nil {
		return order, err
	}
	// таблица delivery
	query = "SELECT order_uid,name,phone,zip,city,address,region,email FROM wb_schema.delivery WHERE order_uid = $1"
	err = r.client.QueryRow(context.Background(), query, order_uid).Scan(&order.Order_uid, &order.Delivery.Name, &order.Delivery.Phone, &order.Delivery.Zip, &order.Delivery.City, &order.Delivery.Address, &order.Delivery.Region, &order.Delivery.Email)
	if err != nil {
		return order, err
	}
	// таблица payment
	query = "SELECT order_uid,request_id,currency,provider,amount,payment_dt,bank,delivery_cost,goods_total,custom_fee FROM wb_schema.payment WHERE order_uid = $1"
	err = r.client.QueryRow(context.Background(), query, order_uid).Scan(&order.Order_uid, &order.Payment.Request_id, &order.Payment.Currency, &order.Payment.Provider, &order.Payment.Amount, &order.Payment.Payment_dt, &order.Payment.Bank, &order.Payment.Delivery_cost, &order.Payment.Goods_total, &order.Payment.Custom_fee)
	if err != nil {
		return order, err
	}
	query = "SELECT items.chrt_id,track_number,price,rid,name,sale,size,total_price,nm_id,brand,status FROM wb_schema.items INNER JOIN wb_schema.order_item ON order_item.chrt_id = items.chrt_id WHERE order_uid = $1"
	rows, err := r.client.Query(context.Background(), query, order_uid)
	if err != nil {
		return order, err
	}
	defer rows.Close()
	// таблицы items и order_item)
	for i := 0; rows.Next(); i++ {
		order.Items = append(order.Items, structures.Items{})
		err = rows.Scan(&order.Items[i].Chrt_id, &order.Items[i].Track_number, &order.Items[i].Price, &order.Items[i].Rid, &order.Items[i].Name, &order.Items[i].Sale, &order.Items[i].Size, &order.Items[i].Total_price, &order.Items[i].Nm_id, &order.Items[i].Brand, &order.Items[i].Status)
		if err != nil {
			return order, err
		}
	}
	return order, err
}

// Удаление заказа из PG
func (r *Repository) DeleteOrder(order_uid string) error {
	ctx := context.Background()
	tx, err := r.client.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	//  поиск и удаление items связанных с order_item
	query := "DELETE FROM wb_schema.items USING wb_schema.order_item WHERE items.chrt_id = order_item.chrt_id AND order_item.order_uid = $1"
	_, err = tx.Exec(ctx, query, order_uid)
	if err != nil {
		return err
	}
	// удаление order_item
	query = "DELETE FROM wb_schema.order_item WHERE order_uid = $1"
	_, err = tx.Exec(ctx, query, order_uid)
	if err != nil {
		return err
	}
	// удаление delivery
	query = "DELETE FROM wb_schema.delivery WHERE order_uid = $1"
	_, err = tx.Exec(ctx, query, order_uid)
	if err != nil {
		return err
	}
	// удаление payment
	query = "DELETE FROM wb_schema.payment WHERE order_uid = $1"
	_, err = tx.Exec(ctx, query, order_uid)
	if err != nil {
		return err
	}
	// удаление orders
	query = "DELETE FROM wb_schema.orders WHERE order_uid = $1"
	_, err = tx.Exec(ctx, query, order_uid)
	if err != nil {
		return err
	}
	tx.Commit(ctx)
	return err
}

package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"syscall"

	c "github.com/ViPDanger/L0/Server/Internal/config"
	n "github.com/ViPDanger/L0/Server/Internal/nats"
	pg "github.com/ViPDanger/L0/Server/Internal/postgres"
	"github.com/ViPDanger/L0/Server/Internal/structures"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill, syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	config := c.Read_Config()
	go func() {

		client, _ := pg.NewClient(context.Background(), config)
		rep := pg.NewRepository(client)
		if rep != nil {
		}
		natsConnection, _ := n.ConnectToNATS(config)
		if natsConnection != nil {
		}

		data := []byte("{\"order_uid\": \"b563feb7b2b84b6test\",\"track_number\": \"WBILMTESTTRACK\",\"entry\": \"WBIL\",\"delivery\": {\"name\": \"Test Testov\",\"phone\": \"+9720000000\",\"zip\": \"2639809\",\"city\": \"Kiryat Mozkin\",\"address\": \"Ploshad Mira 15\",\"region\": \"Kraiot\",\"email\": \"test@gmail.com\"},\"payment\": {\"transaction\": \"b563feb7b2b84b6test\",\"request_id\": \"\",\"currency\": \"USD\",\"provider\": \"wbpay\",\"amount\": 1817,\"payment_dt\": 1637907727,\"bank\": \"alpha\",\"delivery_cost\": 1500,\"goods_total\": 317,\"custom_fee\": 0},\"items\": [{\"chrt_id\": 9934930,\"track_number\": \"WBILMTESTTRACK\",\"price\": 453,\"rid\": \"ab4219087a764ae0btest\",\"name\": \"Mascaras\",\"sale\": 30,\"size\": \"0\",\"total_price\": 317,\"nm_id\": 2389212,\"brand\": \"Vivienne Sabo\",\"status\": 202}],\"locale\": \"en\",\"internal_signature\": \"\",\"customer_id\": \"test\",\"delivery_service\": \"meest\",\"shardkey\": \"9\",\"sm_id\": 99,\"date_created\": \"2021-11-26T06:22:19Z\",\"oof_shard\": \"1\"}")
		var order structures.Order
		err := json.Unmarshal(data, &order)
		if err != nil {
			log.Fatalln("ой", err)
		}
		err = rep.PutOrder(order)
		if err != nil {
			log.Println("Trying to put order:", err)
		}

		order, err = rep.GetOrder(order.Order_uid)
		if err != nil {
			log.Println("Trying to Get order:", err)
		}
		err = rep.DeleteOrder(order.Order_uid)
		if err != nil {
			log.Println("Trying to Get order:", err)
		}
		log.Println("Done")
	}()
	<-ctx.Done()
	log.Println("Server was killed.")

}

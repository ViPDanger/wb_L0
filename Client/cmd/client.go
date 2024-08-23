package main

import (
	"log"
	"time"

	c "github.com/ViPDanger/L0/Client/Internal/config"
	handlers "github.com/ViPDanger/L0/Client/Internal/handler"
	n "github.com/ViPDanger/L0/Client/Internal/nats"
)

func main() {
	config := c.Read_Config()
	var NewServer handlers.HttpClient
	natsConnection, _ := n.ConnectToNATS(config)
	natsHandler := handlers.NatsHanlder{}
	natsHandler.ClientName = "Client"
	natsHandler.NatsConnection = natsConnection
	go natsHandler.InitJetstreamHandler()

	time.Sleep(4 * time.Second)
	data := []byte("{\"order_uid\": \"b563feb7b2b84b6test\",\"track_number\": \"WBILMTESTTRACK\",\"entry\": \"WBIL\",\"delivery\": {\"name\": \"Test Testov\",\"phone\": \"+9720000000\",\"zip\": \"2639809\",\"city\": \"Kiryat Mozkin\",\"address\": \"Ploshad Mira 15\",\"region\": \"Kraiot\",\"email\": \"test@gmail.com\"},\"payment\": {\"transaction\": \"b563feb7b2b84b6test\",\"request_id\": \"\",\"currency\": \"USD\",\"provider\": \"wbpay\",\"amount\": 1817,\"payment_dt\": 1637907727,\"bank\": \"alpha\",\"delivery_cost\": 1500,\"goods_total\": 317,\"custom_fee\": 0},\"items\": [{\"chrt_id\": 9934930,\"track_number\": \"WBILMTESTTRACK\",\"price\": 453,\"rid\": \"ab4219087a764ae0btest\",\"name\": \"Mascaras\",\"sale\": 30,\"size\": \"0\",\"total_price\": 317,\"nm_id\": 2389212,\"brand\": \"Vivienne Sabo\",\"status\": 202}],\"locale\": \"en\",\"internal_signature\": \"\",\"customer_id\": \"test\",\"delivery_service\": \"meest\",\"shardkey\": \"9\",\"sm_id\": 99,\"date_created\": \"2021-11-26T06:22:19Z\",\"oof_shard\": \"1\"}")
	n.PublishMsg(*natsHandler.JetstreamContext, "PutOrder.Server", data)
	n.PublishMsg(*natsHandler.JetstreamContext, "GetOrder.Server", []byte("{\"order_uid\": \"b563feb7b2b84b6test\"}"))
	if err := NewServer.Run(config.Adress, config.Port, natsHandler); err != nil {
		log.Fatal(err)
	}

}

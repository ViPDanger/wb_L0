package main

import (
	"encoding/json"
	"errors"
	"log"
	"sync"
	"time"

	c "github.com/ViPDanger/L0/Client/Internal/config"
	handlers "github.com/ViPDanger/L0/Client/Internal/handler"
	n "github.com/ViPDanger/L0/Client/Internal/nats"
	"github.com/ViPDanger/L0/Client/Internal/structures"
	"github.com/bluele/gcache"
	"github.com/nats-io/nats.go"
)

func main() {
	var mute sync.Mutex
	config := c.Read_Config()
	var NewClient handlers.HttpClient
	natsConnection, _ := n.ConnectToNATS(config)
	natsHandler := handlers.NatsHanlder{}
	natsHandler.ClientName = "Client"
	natsHandler.NatsConnection = natsConnection
	// Запуск NATS Хэндлера
	go func() {
		mute.Lock()
		err := errors.New("")
		for count := config.Con_Attempts; err != nil || count > 0; count-- {
			if err = natsHandler.InitJetstreamHandler(&mute); err != nil {
				log.Println("Retrying")
			}
			time.Sleep(3 * time.Second)
		}
		if err != nil {
			log.Fatalln("Can't reach the server: Number of attempts exceeded.")
		}
	}()
	time.Sleep(1 * time.Second)
	// Ожидание запуска Jetstream
	mute.Lock()
	mute.Unlock()
	//data := []byte("{\"order_uid\": \"b563feb7b2b84b6test\",\"track_number\": \"WBILMTESTTRACK\",\"entry\": \"WBIL\",\"delivery\": {\"name\": \"Test Testov\",\"phone\": \"+9720000000\",\"zip\": \"2639809\",\"city\": \"Kiryat Mozkin\",\"address\": \"Ploshad Mira 15\",\"region\": \"Kraiot\",\"email\": \"test@gmail.com\"},\"payment\": {\"transaction\": \"b563feb7b2b84b6test\",\"request_id\": \"\",\"currency\": \"USD\",\"provider\": \"wbpay\",\"amount\": 1817,\"payment_dt\": 1637907727,\"bank\": \"alpha\",\"delivery_cost\": 1500,\"goods_total\": 317,\"custom_fee\": 0},\"items\": [{\"chrt_id\": 9934930,\"track_number\": \"WBILMTESTTRACK\",\"price\": 453,\"rid\": \"ab4219087a764ae0btest\",\"name\": \"Mascaras\",\"sale\": 30,\"size\": \"0\",\"total_price\": 317,\"nm_id\": 2389212,\"brand\": \"Vivienne Sabo\",\"status\": 202}],\"locale\": \"en\",\"internal_signature\": \"\",\"customer_id\": \"test\",\"delivery_service\": \"meest\",\"shardkey\": \"9\",\"sm_id\": 99,\"date_created\": \"2021-11-26T06:22:19Z\",\"oof_shard\": \"1\"}")
	cache := gcache.New(20).
		LRU().
		Build()
	// Запись в кэш данных из БД
	var msg *nats.Msg
	n.PublishMsg(*natsHandler.JetstreamContext, "AllOrder."+natsHandler.ClientName, []byte(""))
	var allOrder_uid structures.OrderuidList
	msg = <-natsHandler.AllChannel
	json.Unmarshal(msg.Data, &allOrder_uid)
	var order structures.Order
	for _, order_uid := range allOrder_uid.List {
		data, _ := json.Marshal(structures.Order_uid{
			Order_uid: order_uid,
		})
		n.PublishMsg(*natsHandler.JetstreamContext, "GetOrder."+natsHandler.ClientName, data)
		msg = <-natsHandler.GetChannel
		json.Unmarshal(msg.Data, &order)
		cache.Set(order_uid, order)
	}
	// Запуск Http Клиента
	if err := NewClient.Run(config.Adress, config.Port, natsHandler, &cache); err != nil {
		log.Fatal(err)
	}
	natsHandler.CloseJetstreamHandler()
}

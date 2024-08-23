package handler

import (
	"context"
	"log"
	"time"

	"encoding/json"

	"strings"

	"sync"

	n "github.com/ViPDanger/L0/Server/Internal/nats"
	pg "github.com/ViPDanger/L0/Server/Internal/postgres"
	"github.com/ViPDanger/L0/Server/Internal/structures"
	"github.com/nats-io/nats.go"
)

type NatsHanlder struct {
	NatsConnection   *nats.Conn
	PgRepository     *pg.Repository
	JetstreamContext *nats.JetStreamContext
}

// инициализация NATS Jetstream хэндлера для сервера
func (nh *NatsHanlder) InitJetstreamHandler() {
	var mutex sync.Mutex
	Jetstream, _ := n.NatsJetStream(nh.NatsConnection)
	nh.JetstreamContext = &Jetstream
	ctxJetstream, JetCancel := context.WithCancel(context.Background())
	defer JetCancel()
	n.ClearStream(Jetstream, "GetOrder")
	n.ClearStream(Jetstream, "PutOrder")
	n.ClearStream(Jetstream, "DelOrder")
	n.CreateStream(ctxJetstream, Jetstream, "stream_GetOrder", []string{"GetOrder", "GetOrder.Server"})
	n.CreateStream(ctxJetstream, Jetstream, "stream_PutOrder", []string{"PutOrder", "PutOrder.Server"})
	n.CreateStream(ctxJetstream, Jetstream, "stream_DelOrder", []string{"DelOrder", "DelOrder.Server"})

	// GetOrder fetch
	go func() {
		ctxConsumer, ConCancel := context.WithCancel(context.Background())
		defer ConCancel()
		n.CreateConsumer(ctxConsumer, Jetstream, "GetOrder", "stream_GetOrder")
		getOrderSub, _ := n.Subscribe(ctxConsumer, Jetstream, "GetOrder", "GetOrder", "stream_GetOrder")
		for {
			msg, _ := n.FetchOne(ctxConsumer, getOrderSub)
			if msg != nil {
				mutex.Lock()
				var order_uid structures.Order_uid

				err := json.Unmarshal(msg.Data, &order_uid)
				if err != nil {
					log.Println(msg.Subject, ": Unmarshal error")

					break
				}

				order, err := nh.PgRepository.GetOrder(order_uid.Order_uid)
				if err != nil {
					log.Println(msg.Subject, ": GetOrder in repository error", err)
					n.PublishMsg(*nh.JetstreamContext, strings.TrimPrefix(msg.Subject, "GetOrder.")+".GetOrder", []byte(err.Error()))
				} else {
					data, err := json.Marshal(order)
					if err != nil {
						log.Println(msg.Subject, ":  Marshal error", err)
						n.PublishMsg(*nh.JetstreamContext, strings.TrimPrefix(msg.Subject, "GetOrder.")+".GetOrder", []byte(err.Error()))
					} else {
						n.PublishMsg(*nh.JetstreamContext, strings.TrimPrefix(msg.Subject, "GetOrder.")+".GetOrder", data)
						log.Println(msg.Subject, ": ", string(msg.Data))
					}
				}
				mutex.Unlock()
			}
		}
	}()

	// DelOrder fetch
	go func() {
		ctxConsumer, ConCancel := context.WithCancel(context.Background())
		defer ConCancel()
		n.CreateConsumer(ctxConsumer, Jetstream, "DelOrder", "stream_DelOrder")
		delOrderSub, _ := n.Subscribe(ctxConsumer, Jetstream, "DelOrder", "DelOrder", "stream_DelOrder")
		for {
			msg, _ := n.FetchOne(ctxConsumer, delOrderSub)
			if msg != nil {
				mutex.Lock()
				var order_uid structures.Order_uid
				msg.Ack()
				log.Println(msg.Subject, ": ", string(msg.Data))
				err := json.Unmarshal(msg.Data, &order_uid)
				if err != nil {
					log.Println(msg.Subject, "del: Unmarshal error")
					n.PublishMsg(*nh.JetstreamContext, strings.TrimPrefix(msg.Subject, "DelOrder.")+".DelOrder", []byte(err.Error()))
				} else {
					err = nh.PgRepository.DeleteOrder(order_uid.Order_uid)
					if err != nil {
						log.Println(msg.Subject, ": DelOrder in repository error", err)
						n.PublishMsg(*nh.JetstreamContext, strings.TrimPrefix(msg.Subject, "DelOrder.")+".DelOrder", []byte(err.Error()))
					} else {
						n.PublishMsg(*nh.JetstreamContext, strings.TrimPrefix(msg.Subject, "DelOrder.")+".DelOrder", []byte(order_uid.Order_uid+"Удалён успешно"))
					}
				}

				mutex.Unlock()
			}
		}
	}()

	// PutOrder fetch
	func() {
		ctxConsumer, ConCancel := context.WithCancel(context.Background())
		defer ConCancel()
		n.CreateConsumer(ctxConsumer, Jetstream, "PutOrder", "stream_PutOrder")
		putOrderSub, _ := n.Subscribe(ctxConsumer, Jetstream, "PutOrder", "PutOrder", "stream_PutOrder")
		for {
			msg, _ := n.FetchOne(ctxConsumer, putOrderSub)

			if msg != nil {
				mutex.Lock()
				var order structures.Order
				msg.Ack()
				log.Println(msg.Subject, ": ", string(msg.Data))
				err := json.Unmarshal(msg.Data, &order)
				if err != nil {
					log.Println(msg.Subject, ": Unmarshal error")
					n.PublishMsg(*nh.JetstreamContext, strings.TrimPrefix(msg.Subject, "PutOrder.")+".PutOrder", []byte(err.Error()))
				} else {
					err = nh.PgRepository.PutOrder(order)
					if err != nil {
						log.Println(msg.Subject, ": PutOrder in repository error", err)
						n.PublishMsg(*nh.JetstreamContext, strings.TrimPrefix(msg.Subject, "PutOrder.")+".PutOrder", []byte(err.Error()))
					} else {
						n.PublishMsg(*nh.JetstreamContext, strings.TrimPrefix(msg.Subject, "PutOrder.")+".PutOrder", []byte(order.Order_uid+"Добавлен успешно"))
					}
				}
				mutex.Unlock()
			}
		}
	}()

}

func (nh *NatsHanlder) CloseJetstreamHandler() {
	(*nh.JetstreamContext).DeleteStream(<-(*nh.JetstreamContext).StreamNames())
	time.Sleep(2 * time.Second)
	nh.NatsConnection.Close()
	log.Println("Jetstream: Closed")
}

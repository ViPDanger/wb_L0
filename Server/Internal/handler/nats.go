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
	jetCancel        *context.CancelFunc
}

// инициализация NATS Jetstream хэндлера для сервера
func (nh *NatsHanlder) InitJetstreamHandler() {
	var mutex sync.Mutex
	Jetstream, _ := n.NatsJetStream(nh.NatsConnection)
	nh.JetstreamContext = &Jetstream
	ctxJetstream, JetCancel := context.WithCancel(context.Background())
	defer JetCancel()
	nh.jetCancel = &JetCancel
	n.CreateStream(ctxJetstream, Jetstream, "stream_GetOrder", []string{"GetOrder", "GetOrder.Server"})
	n.CreateStream(ctxJetstream, Jetstream, "stream_PutOrder", []string{"PutOrder", "PutOrder.Server"})
	n.CreateStream(ctxJetstream, Jetstream, "stream_DelOrder", []string{"DelOrder", "DelOrder.Server"})
	n.CreateStream(ctxJetstream, Jetstream, "stream_AllOrder", []string{"AllOrder", "AllOrder.Server"})
	// GetOrder fetch
	go func() {
		ctxConsumer, ConCancel := context.WithCancel(ctxJetstream)
		defer ConCancel()
		n.CreateConsumer(ctxConsumer, Jetstream, "GetOrder", "stream_GetOrder")
		getOrderSub, _ := n.Subscribe(ctxConsumer, Jetstream, "GetOrder", "GetOrder", "stream_GetOrder")
		for ctxJetstream.Err() == nil {
			msg, _ := n.FetchOne(ctxConsumer, getOrderSub)
			if msg != nil {
				mutex.Lock()
				var order_uid structures.Order_uid
				msg.Ack()
				log.Println(msg.Subject, ": ", " Data recieved")
				err := json.Unmarshal(msg.Data, &order_uid)
				if err != nil {
					log.Println(msg.Subject, ": Unmarshal error")
					break
				}

				order, err := nh.PgRepository.GetOrder(order_uid.Order_uid)
				if err != nil {
					log.Println(msg.Subject, ": GetOrder Repository error", err)
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
		ctxConsumer, ConCancel := context.WithCancel(ctxJetstream)
		defer ConCancel()
		n.CreateConsumer(ctxConsumer, Jetstream, "DelOrder", "stream_DelOrder")
		delOrderSub, _ := n.Subscribe(ctxConsumer, Jetstream, "DelOrder", "DelOrder", "stream_DelOrder")
		for ctxJetstream.Err() == nil {
			msg, _ := n.FetchOne(ctxConsumer, delOrderSub)
			if msg != nil {
				mutex.Lock()
				var order_uid structures.Order_uid
				msg.Ack()
				log.Println(msg.Subject, ": ", " Data recieved")
				err := json.Unmarshal(msg.Data, &order_uid)
				if err != nil {
					log.Println(msg.Subject, "del: Unmarshal error")
					n.PublishMsg(*nh.JetstreamContext, strings.TrimPrefix(msg.Subject, "DelOrder.")+".DelOrder", []byte(err.Error()))
				} else {
					err = nh.PgRepository.DeleteOrder(order_uid.Order_uid)
					if err != nil {
						log.Println(msg.Subject, ": DelOrder Repository error", err)
						n.PublishMsg(*nh.JetstreamContext, strings.TrimPrefix(msg.Subject, "DelOrder.")+".DelOrder", []byte(err.Error()))
					} else {
						n.PublishMsg(*nh.JetstreamContext, strings.TrimPrefix(msg.Subject, "DelOrder.")+".DelOrder", []byte("Order with order_uid: "+order_uid.Order_uid+" has been removed"))
					}
				}

				mutex.Unlock()
			}
		}
	}()
	// AllOrder fetch
	go func() {
		ctxConsumer, ConCancel := context.WithCancel(ctxJetstream)
		defer ConCancel()
		n.CreateConsumer(ctxConsumer, Jetstream, "AllOrder", "stream_AllOrder")
		allOrderSub, _ := n.Subscribe(ctxConsumer, Jetstream, "AllOrder", "AllOrder", "stream_AllOrder")
		for ctxJetstream.Err() == nil {
			msg, _ := n.FetchOne(ctxConsumer, allOrderSub)
			if msg != nil {
				mutex.Lock()
				msg.Ack()
				log.Println(msg.Subject, ": ", " msg recieved")
				orderuidList, err := nh.PgRepository.AllOrder_uid()
				if err != nil {
					log.Println(msg.Subject, ": AllOrder Repository error", err)
					n.PublishMsg(*nh.JetstreamContext, strings.TrimPrefix(msg.Subject, "AllOrder.")+".AllOrder", []byte(err.Error()))
				} else {
					data, err := json.Marshal(orderuidList)
					if err != nil {
						log.Println(msg.Subject, ": AllOrder Repository error", err)
						n.PublishMsg(*nh.JetstreamContext, strings.TrimPrefix(msg.Subject, "AllOrder.")+".AllOrder", []byte(err.Error()))
					}
					n.PublishMsg(*nh.JetstreamContext, strings.TrimPrefix(msg.Subject, "AllOrder.")+".AllOrder", data)
				}
				mutex.Unlock()
			}
		}
	}()
	// PutOrder fetch
	func() {
		ctxConsumer, ConCancel := context.WithCancel(ctxJetstream)
		defer ConCancel()
		n.CreateConsumer(ctxConsumer, Jetstream, "PutOrder", "stream_PutOrder")
		putOrderSub, _ := n.Subscribe(ctxConsumer, Jetstream, "PutOrder", "PutOrder", "stream_PutOrder")
		for ctxJetstream.Err() == nil {
			msg, _ := n.FetchOne(ctxConsumer, putOrderSub)

			if msg != nil {
				mutex.Lock()
				var order structures.Order
				msg.Ack()
				log.Println(msg.Subject, ": ", " Data recieved")
				err := json.Unmarshal(msg.Data, &order)
				if err != nil {
					log.Println(msg.Subject, ": Unmarshal error")
					n.PublishMsg(*nh.JetstreamContext, strings.TrimPrefix(msg.Subject, "PutOrder.")+".PutOrder", []byte(err.Error()))
				} else {
					err = nh.PgRepository.PutOrder(order)
					if err != nil {
						log.Println(msg.Subject, ": PutOrder Repository error", err)
						n.PublishMsg(*nh.JetstreamContext, strings.TrimPrefix(msg.Subject, "PutOrder.")+".PutOrder", []byte(err.Error()))
					} else {
						n.PublishMsg(*nh.JetstreamContext, strings.TrimPrefix(msg.Subject, "PutOrder.")+".PutOrder", []byte("Order with order_uid: "+order.Order_uid+" has been added to the base"))
					}
				}
				mutex.Unlock()
			}
		}
	}()

}

func (nh *NatsHanlder) CloseJetstreamHandler() {
	(*nh.jetCancel)()
	time.Sleep(1 * time.Second)
	(*nh.JetstreamContext).DeleteConsumer("stream_DelOrder", "DelOrder")
	(*nh.JetstreamContext).DeleteStream("stream_DelOrder")
	log.Println("Jetstream: stream DelOrder closed")
	(*nh.JetstreamContext).DeleteConsumer("stream_GetOrder", "GetOrder")
	(*nh.JetstreamContext).DeleteStream("stream_GetOrder")
	log.Println("Jetstream: stream GetOrder closed")
	(*nh.JetstreamContext).DeleteConsumer("stream_PutOrder", "PutOrder")
	(*nh.JetstreamContext).DeleteStream("stream_PutOrder")
	log.Println("Jetstream: stream PutOrder closed")
	(*nh.JetstreamContext).DeleteConsumer("stream_AllOrder", "AllOrder")
	(*nh.JetstreamContext).DeleteStream("stream_AllOrder")
	log.Println("Jetstream: stream AllOrder closed")
	nh.NatsConnection.Close()
	log.Println("Jetstream: Closed")
}

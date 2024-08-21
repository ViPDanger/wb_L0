package handler

import (
	"encoding/json"
	"log"
	"time"

	pg "github.com/ViPDanger/L0/Server/Internal/postgres"
	"github.com/ViPDanger/L0/Server/Internal/structures"
	"github.com/nats-io/nats.go"
)

type NatsHanlder struct {
	natsConnection *nats.Conn
	pgRepository   *pg.Repository
	cache          gcache.LRUCache
}

func (nh *NatsHanlder) initNatsHandler() error {
	nh.natsConnection.Subscribe("orders", func(m *nats.Msg) {
		order := &structures.Order{}
		err := json.Unmarshal(m.Data, order)
		if err != nil {
			log.Println("NATS Handler: ", err)
			return
		}

		err = nh.pgRepository.PutOrder(*order)
		if err != nil {
			log.Println("NATS Handler: ", err)
			return
		}
		log.Println("NATS Handler: ", err)
	})
}

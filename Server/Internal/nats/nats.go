package nats

import (
	"log"

	"github.com/ViPDanger/L0/Server/Internal/config"
	"github.com/nats-io/nats.go"
)

func ConnectToNATS(conf config.CFG) (*nats.Conn, error) {
	var err error
	var nc *nats.Conn
	for attempts := conf.Con_Attempts; attempts > 0; attempts-- {
		nc, err = nats.Connect("nats://" + conf.Nats_host + ":" + conf.Nats_port)
		if err != nil {
			log.Println("Connecting to Nats: Nats didn't respound")
		}
	}
	if err != nil {
		log.Fatalln("Connecting to Nats: Number of attempts exceeded.")

	}
	return nc, err
}

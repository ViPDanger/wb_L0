package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	c "github.com/ViPDanger/L0/Server/Internal/config"
	natsHandler "github.com/ViPDanger/L0/Server/Internal/handler"
	n "github.com/ViPDanger/L0/Server/Internal/nats"
	pg "github.com/ViPDanger/L0/Server/Internal/postgres"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill, syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	config := c.Read_Config()
	jetstreamHandler := natsHandler.NatsHanlder{}
	go func() {
		client, _ := pg.NewClient(context.Background(), config)
		pgRep := pg.NewRepository(client)
		natsConnection, _ := n.ConnectToNATS(config)
		jetstreamHandler.NatsConnection = natsConnection
		jetstreamHandler.PgRepository = pgRep
		go jetstreamHandler.InitJetstreamHandler()
	}()
	<-ctx.Done()
	log.Println("Closing Jetstream Handler")
	jetstreamHandler.CloseJetstreamHandler()
	log.Println("Server was killed.")

}

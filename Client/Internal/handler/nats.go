package handler

import (
	"context"
	"log"
	"strings"
	"sync"
	"time"

	n "github.com/ViPDanger/L0/Client/Internal/nats"
	"github.com/nats-io/nats.go"
)

type NatsHanlder struct {
	ClientName       string
	NatsConnection   *nats.Conn
	JetstreamContext *nats.JetStreamContext
	jetCancel        *context.CancelFunc
	DelChannel       chan *nats.Msg
	PutChannel       chan *nats.Msg
	GetChannel       chan *nats.Msg
	AllChannel       chan *nats.Msg
}

// инициализация NATS Jetstream хэндлера для клиента
func (nh *NatsHanlder) InitJetstreamHandler(mute *sync.Mutex) error {
	nh.DelChannel = make(chan *nats.Msg, 1)
	nh.PutChannel = make(chan *nats.Msg, 1)
	nh.GetChannel = make(chan *nats.Msg, 1)
	nh.AllChannel = make(chan *nats.Msg, 1)
	Jetstream, _ := n.NatsJetStream(nh.NatsConnection)
	nh.JetstreamContext = &Jetstream
	ctxJetstream, JetCancel := context.WithCancel(context.Background())
	defer JetCancel()
	nh.jetCancel = &JetCancel
	streamGetOrder, err := Jetstream.StreamInfo("stream_GetOrder")
	if err != nil {
		log.Println("StreamInfo: can't reach Server")
		return err
	}
	streamPutOrder, err := Jetstream.StreamInfo("stream_PutOrder")
	if err != nil {

		log.Println("StreamInfo: can't reach Server ")
		return err
	}
	stream_DelOrder, err := Jetstream.StreamInfo("stream_DelOrder")
	if err != nil {
		log.Println("StreamInfo: can't reach Server")
		return err
	}
	stream_AllOrder, err := Jetstream.StreamInfo("stream_AllOrder")
	if err != nil {
		log.Println("StreamInfo: can't reach Server")
		return err
	}
	streamGetOrder.Config.Subjects = append(streamGetOrder.Config.Subjects, "GetOrder."+nh.ClientName)
	_, err = Jetstream.UpdateStream(&streamGetOrder.Config)
	if err != nil {
		if err.Error() != "nats: duplicate subjects detected" {
			log.Println("StreamInfo: can't update stream_GetOrder ", err)
			return err
		}
	}
	streamPutOrder.Config.Subjects = append(streamPutOrder.Config.Subjects, "PutOrder."+nh.ClientName)
	_, err = Jetstream.UpdateStream(&streamPutOrder.Config)
	if err != nil {
		if err.Error() != "nats: duplicate subjects detected" {
			log.Println("StreamInfo: can't update stream_PutOrder ", err)
			return err
		}
	}
	stream_DelOrder.Config.Subjects = append(stream_DelOrder.Config.Subjects, "DelOrder."+nh.ClientName)
	_, err = Jetstream.UpdateStream(&stream_DelOrder.Config)
	if err != nil {
		if err.Error() != "nats: duplicate subjects detected" {
			log.Println("StreamInfo: can't update stream_DelOrder ", err)
			return err
		}
	}
	stream_AllOrder.Config.Subjects = append(stream_AllOrder.Config.Subjects, "AllOrder."+nh.ClientName)
	_, err = Jetstream.UpdateStream(&stream_AllOrder.Config)
	if err != nil {
		if err.Error() != "nats: duplicate subjects detected" {
			log.Println("StreamInfo: can't update stream_AllOrder ", err)
			return err
		}
	}
	n.CreateStream(ctxJetstream, Jetstream, "stream_"+nh.ClientName, []string{nh.ClientName + ".DelOrder", nh.ClientName + ".PutOrder", nh.ClientName + ".GetOrder", nh.ClientName + ".AllOrder", nh.ClientName})
	n.CreateConsumer(ctxJetstream, Jetstream, nh.ClientName, "stream_"+nh.ClientName)
	ClientSubscription, err := n.Subscribe(ctxJetstream, *nh.JetstreamContext, nh.ClientName, nh.ClientName, "stream_"+nh.ClientName)
	mute.Unlock()
	for ctxJetstream.Err() == nil {
		msg, err := n.FetchOne(ctxJetstream, ClientSubscription)
		if err == nil {
			log.Println(msg.Subject, ": Data Recieved")
			switch strings.TrimPrefix(msg.Subject, nh.ClientName+".") {
			case "DelOrder":
				msg.Ack()
				nh.DelChannel <- msg
			case "PutOrder":
				msg.Ack()
				nh.PutChannel <- msg
			case "GetOrder":
				msg.Ack()
				nh.GetChannel <- msg
			case "AllOrder":
				msg.Ack()
				nh.AllChannel <- msg
			default:
			}
		}
	}
	return nil
}
func (nh *NatsHanlder) CloseJetstreamHandler() {
	(*nh.jetCancel)()
	time.Sleep(1 * time.Second)
	(*nh.JetstreamContext).DeleteConsumer("stream_"+nh.ClientName, nh.ClientName)
	(*nh.JetstreamContext).DeleteStream("stream_" + nh.ClientName)
	log.Println("Jetstream: stream", nh.ClientName, "closed")
	nh.NatsConnection.Close()
	log.Println("Jetstream: Closed")
}

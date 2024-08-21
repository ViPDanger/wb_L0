package nats

import (
	"context"
	"log"
	"time"

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
		} else {
			break
		}
		time.Sleep(5 * time.Second)
	}
	if err != nil {
		log.Fatalln("Connecting to Nats: Number of attempts exceeded.")

	}

	log.Println("Connecting to Nats: Successfully ")
	return nc, err
}

func NatsJetStream(nc *nats.Conn) (nats.JetStreamContext, error) {
	jsCtx, err := nc.JetStream()
	if err != nil {
		log.Println("Jetstream: ", err)
		return nil, err
	}
	log.Println("Jetstream: Successfully ")
	return jsCtx, nil
}
func ClearStream(jsCtx nats.JetStreamContext, subjectName string) error {
	streamName, err := jsCtx.StreamNameBySubject(subjectName)
	if err != nil {
		log.Println("Jetstream delete stream by subject", subjectName, ": ", err)
	}
	err = jsCtx.DeleteStream(streamName)
	if err == nil {
		log.Println("Jetstream delete stream ", streamName, ": Successfully")
	}
	return err
}
func CreateStream(ctx context.Context, jsCtx nats.JetStreamContext, StreamName string, subjects []string) (*nats.StreamInfo, error) {
	err := jsCtx.DeleteStream(StreamName)
	if err == nil {
		log.Println("Jetstream delete stream: Successfully")
	}
	stream, err := jsCtx.AddStream(&nats.StreamConfig{
		Name:              StreamName,
		Subjects:          subjects,
		Retention:         nats.InterestPolicy, // remove acked messages
		Discard:           nats.DiscardOld,     // when the stream is full, discard old messages
		MaxAge:            24 * time.Hour,      // max age of stored messages is 7 days
		Storage:           nats.FileStorage,    // type of message storage
		MaxMsgsPerSubject: 100_000_000,         // max stored messages per subject
		MaxMsgSize:        4 << 20,             // max single message size is 4 MB
		NoAck:             false,               // we need the "ack" system for the message queue system
	}, nats.Context(ctx))
	if err != nil {
		log.Println("Jetstream add stream:", err)
		return nil, err
	}
	log.Println("Jetstream add stream: Successfully")
	return stream, nil
}

func PublishMsg(jsCtx nats.JetStreamContext, subject string, payload []byte) error {
	ack, err := jsCtx.Publish(subject, payload)
	if err != nil {
		log.Println("Jetstream publish: ", err)
		return err
	}
	log.Println(ack.Stream, subject, "publish: Successfully")
	return nil
}

func CreateConsumer(ctx context.Context, jsCtx nats.JetStreamContext, consumerGroupName, streamName string) (*nats.ConsumerInfo, error) {
	consumer, err := jsCtx.AddConsumer(streamName, &nats.ConsumerConfig{
		Durable:       consumerGroupName,      // durable name is the same as consumer group name
		DeliverPolicy: nats.DeliverAllPolicy,  // deliver all messages, even if they were sent before the consumer was created
		AckPolicy:     nats.AckExplicitPolicy, // ack messages manually
		AckWait:       5 * time.Second,        // wait for ack for 5 seconds
		MaxAckPending: -1,                     // unlimited number of pending acks
	}, nats.Context(ctx))
	if err != nil {
		log.Println("Jetstream add consumer:", err)
		return nil, err
	}
	log.Println("Jetstream add consumer: Successfully")
	return consumer, nil
}

func Subscribe(ctx context.Context, js nats.JetStreamContext, subject, consumerGroupName, streamName string) (*nats.Subscription, error) {
	pullSub, err := js.PullSubscribe(
		subject,
		consumerGroupName,
		nats.ManualAck(),                         // ack messages manually
		nats.Bind(streamName, consumerGroupName), // bind consumer to the stream
		nats.Context(ctx),                        // use context to cancel the subscription
	)
	if err != nil {
		log.Println(subject, " subscribe:", err)
		return nil, err
	}
	log.Println(subject, " subscribe: Successfully")
	return pullSub, nil
}

func FetchOne(ctx context.Context, pullSub *nats.Subscription) (*nats.Msg, error) {
	msgs, err := pullSub.Fetch(1, nats.Context(ctx))
	if err != nil {
		if err.Error() != "context deadline exceeded" {
			log.Println("Jetstream fetch:", err)
		}
		return nil, err
	}
	if len(msgs) == 0 {
		log.Println("Jetstream fetch: No Masseges")
		return nil, nil
	}
	msgs[0].Ack()
	return msgs[0], nil
}

package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/jami/kafka-eb-collector/src"
)

const (
	// CollectorGroupID group id for kafka
	CollectorGroupID = "DefaultCollectorGroup"
)

var (
	config           src.CLIConfiguration
	eventbusProducer *src.Producer
	eventbusConsumer *src.Consumer
	store            src.CollectorStore
)

type eventbusListener struct {
	src.ConsumeHandler
}

// Consume a chunk
func (ebl eventbusListener) Consume(b []byte) error {
	header := src.EventEnvelop{}
	if err := header.FromBytes(b); err != nil {
		return err
	}

	// check for valid event type
	var gce src.GroupCollectorEnumeration
	if err := gce.FromString(header.Type); err != nil {
		return err
	}

	// handle types of interests
	switch gce {
	case src.EventGroupCollectorCreate:
		return src.GroupCreateHandler(store, b, eventbusProducer)
	case src.EventGroupCollectorEntityDone:
		return src.EntityDoneHandler(store, b, eventbusProducer)
	}

	return nil
}

// run the listeners
func run() {
	go func() {
		eventListener := eventbusListener{}
		log.Fatal(eventbusConsumer.Listen(eventListener, []string{
			config.EventBusTopic,
		}))
	}()
}

func main() {
	var err error

	log.Println("kafka event bus collector")
	err = config.Parse()

	if err != nil {
		log.Println(err.Error())
		flag.PrintDefaults()
		os.Exit(1)
	}

	if config.ShowHelp {
		flag.PrintDefaults()
		os.Exit(0)
	}

	// producer configuration
	eventbusProducer, err = src.NewProducer(&src.ProducerConfiguration{
		Broker: config.Broker,
		Topic:  config.EventBusTopic,
	})

	if err != nil {
		log.Fatal(err)
	}

	// consumer configuration
	consumerConfig := src.NewConsumerConfiguration()
	consumerConfig.GroupID = CollectorGroupID
	consumerConfig.Broker = config.Broker
	eventbusConsumer = src.NewConsumer(consumerConfig)

	// simple in mem store
	store = src.CreateSimpleInMemStore()

	// signal channel
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, os.Interrupt, syscall.SIGTERM)

	// run the application
	run()

	// wait for signals
	<-sigchan
}

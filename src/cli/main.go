package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/jami/kafka-eb-collector/src"
)

var (
	config           src.CLIConfiguration
	eventbusProducer *src.Producer
	eventbusConsumer *src.Consumer
)

func run() {
	go func() {
		eventListener := src.NewEventListener(a)
		eventbusConsumer.Listen(eventListener, []string{
			a.config.EventBusTopic,
			a.config.MetaTopic,
		})
	}()

	go func() {
		log.Fatal(src.RestListenerListen(config.ListenerEndpoint))
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

	eventbusProducer, err = src.NewProducer(&src.ProducerConfiguration{
		Broker: config.Broker,
		Topic:  config.EventBusTopic,
	})

	if err != nil {
		log.Fatal(err)
	}

	// signal channel
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, os.Interrupt, syscall.SIGTERM)

	// run the application
	run()

	// wait for signals
	<-sigchan
}

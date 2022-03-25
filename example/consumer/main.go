package main

import (
	"errors"
	"flag"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/nsqio/go-nsq"
	"github.com/sirupsen/logrus"
)

var (
	topices  = flag.String("tpc", "test_topic", "topices for consume")
	channel  = flag.String("ch", "test_chan", "channel name for consume")
	addrs    = flag.String("addr", "localhost:4161", "nsqlookupd cluster http host:port")
	addrType = flag.String("type", "nsqlookupd", "nsqlookupd|nsqd")
)

func main() {
	flag.Parse()

	cfg := nsq.NewConfig()

	for _, topic := range strings.Split(*topices, ",") {
		consumer, err := nsq.NewConsumer(topic, *channel, cfg)
		if err != nil {
			panic(err)
		}
		defer consumer.Stop()
		consumer.AddConcurrentHandlers(nsq.HandlerFunc(func(m *nsq.Message) error {
			if len(m.Body) == 0 {
				return errors.New("body is blank re-enqueue message")
			}
			logrus.Printf("receive: %s", m.Body)
			return nil
		}), 100)

		if *addrType == "nsqd" {
			if err := consumer.ConnectToNSQDs(strings.Split(*addrs, ",")); err != nil {
				panic(err)
			}
		} else {
			if err := consumer.ConnectToNSQLookupds(strings.Split(*addrs, ",")); err != nil {
				panic(err)
			}
		}
	}

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT)
	<-shutdown
}

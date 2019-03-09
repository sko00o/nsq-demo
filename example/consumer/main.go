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

/**
nsqlookupd & 
nsqd --lookupd-tcp-address=127.0.0.1:4160 &
nsqadmin --lookupd-http-address=127.0.0.1:4161 &
**/


var (
	topic   = flag.String("tpc", "test_topic", "topic for consume")
	channel = flag.String("ch", "test_chan", "channel name for consume")
	addrs   = flag.String("addr", "localhost:4161", "nsqlookupd cluster http host:port")
)

func main() {
	cfg := nsq.NewConfig()

	consumer, err := nsq.NewConsumer(*topic, *channel, cfg)
	if err != nil {
		panic(err)
	}
	consumer.AddConcurrentHandlers(nsq.HandlerFunc(func(m *nsq.Message) error {
		if len(m.Body) == 0 {
			return errors.New("body is blank re-enqueue message")
		}
		logrus.Printf("receive: %s", m.Body)
		return nil
	}), 100)
	if err := consumer.ConnectToNSQLookupds(strings.Split(*addrs, ",")); err != nil {
		panic(err)
	}

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT)
	for {
		select {
		case <-consumer.StopChan:
			return
		case <-shutdown:
			consumer.Stop()
		}
	}
}

package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/nsqio/go-nsq"
	"github.com/sirupsen/logrus"
)

var (
	topic = flag.String("tpc", "test_topic", "topic for consume")
	addr  = flag.String("addr", "localhost:4150", "nsqd http host:port")
)

func main() {
	cfg := nsq.NewConfig()
	producer, err := nsq.NewProducer(*addr, cfg)
	if err != nil {
		panic(err)
	}

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT)
	ticker := time.NewTicker(2 * time.Second)
	for {
		select {
		case <-shutdown:
			producer.Stop()
			return
		case <-ticker.C:
			msg := []byte(time.Now().String())
			if err := producer.Publish(*topic, msg); err != nil {
				logrus.Errorf("putlish error: %s", err)
				continue
			}
			logrus.Infof("send: %s", msg)
		}
	}
}

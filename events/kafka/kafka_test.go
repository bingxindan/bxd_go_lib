package kafka

import (
	"context"
	"fmt"
	"github.com/bingxindan/bxd_go_lib/events/event"
	"os"
	"os/signal"
	"syscall"
	"testing"
)

func TestReceive(t *testing.T) {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	receiver, err := NewKafkaReceiver([]string{"localhost:9092"}, "kratos")
	if err != nil {
		panic(err)
	}
	receive(receiver)

	<-sigs
	_ = receiver.Close()
}

func receive(receiver event.Receiver) {
	fmt.Println("start receiver")
	err := receiver.Receive(context.Background(), func(ctx context.Context, msg event.Event) error {
		fmt.Printf("key:%s, value:%s\n", msg.Key(), msg.Value())
		return nil
	})
	if err != nil {
		return
	}
}

func TestSend(t *testing.T) {
	sender, err := NewKafkaSender([]string{"localhost:9092"}, "kratos")
	if err != nil {
		panic(err)
	}

	for i := 0; i < 50; i++ {
		send(sender)
	}

	_ = sender.Close()
}

func send(sender event.Sender) {
	msg := NewMessage("kratos", []byte("hello world"))
	err := sender.Send(context.Background(), msg)
	if err != nil {
		panic(err)
	}
	fmt.Printf("key:%s, value:%s\n", msg.Key(), msg.Value())
}

package kafka

import (
	"context"
	"fmt"
	"github.com/bingxindan/bxd_go_lib/events/event"
	"github.com/bingxindan/bxd_go_lib/logger"
	"testing"
)

func TestReceive(t *testing.T) {
	receiver, err := NewKafkaReceiver([]string{"192.168.31.156:29092"}, "topic", "content")
	if err != nil {
		panic(err)
	}
	receive(receiver)
	defer receiver.Close()
}

func receive(receiver event.Receiver) {
	fmt.Println("start receiver")
	err := receiver.Receive(context.Background(), func(ctx context.Context, msg event.Event) error {
		fmt.Printf("key:%s, value:%s\n", msg.Key(), msg.Value())
		logger.Ix(ctx, "kafka.receive", "key:%s, value:%s\n", msg.Key(), msg.Value())
		return nil
	})
	fmt.Printf("receive: %+v\n", err)
	if err != nil {
		return
	}
}

func TestSend(t *testing.T) {
	sender, err := NewKafkaSender([]string{"192.168.31.156:29092"}, "topic")
	if err != nil {
		panic(err)
	}

	for i := 0; i < 50; i++ {
		send(sender)
	}

	_ = sender.Close()
}

func send(sender event.Sender) {
	msg := NewMessage("topic", []byte("hello world"))
	err := sender.Send(context.Background(), msg)
	if err != nil {
		panic(err)
	}
	fmt.Printf("key:%s, value:%s\n", msg.Key(), msg.Value())
}

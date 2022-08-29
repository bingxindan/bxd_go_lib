package discovery

import (
	"context"
	"fmt"
	"github.com/bingxindan/bxd_go_lib/logger"
	"github.com/bingxindan/bxd_go_lib/register/etcd/discovery/helloworld"
	"google.golang.org/grpc"
	"google.golang.org/grpc/resolver"
	"net"
	"testing"
	"time"
)

var (
	etcdAddrs = []string{"0.0.0.0:2379"}
)

func TestReso(t *testing.T) {
	r := NewResolver(etcdAddrs)

	resolver.Register(r)

	conn, err := grpc.Dial("etcd:///hello", grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"round_robin"}`), grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		logger.F("aaa", "fail.to.dial: %+v", err)
	}

	defer conn.Close()

	c := helloworld.NewGreeterClient(conn)

	for i := 0; i < 10; i++ {
		resp, err := c.SayHello(context.Background(), &helloworld.HelloRequest{Name: "acc"})
		if err != nil {
			t.Fatalf("say hello failed %v", err)
		}
		println(resp.Message)
		time.Sleep(100 * time.Millisecond)
	}

	time.Sleep(200000 * time.Minute)
	select {}
}

func TestResolver(t *testing.T) {
	r := NewResolver(etcdAddrs)

	resolver.Register(r)

	// etcd注册5个服务
	go newServer(t, ":1001", "1.0", 1)
	go newServer(t, ":1002", "1.0", 1)
	go newServer(t, ":1003", "1.0", 1)
	go newServer(t, ":1004", "1.0", 1)
	go newServer(t, ":1005", "1.0", 10)

	time.Sleep(5 * time.Second)

	conn, err := grpc.Dial("etcd:///hello", grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"round_robin"}`), grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		logger.F("aaa", "fail.to.dial: %+v", err)
	}

	defer conn.Close()

	c := helloworld.NewGreeterClient(conn)

	for i := 0; i < 10; i++ {
		resp, err := c.SayHello(context.Background(), &helloworld.HelloRequest{Name: "acc"})
		if err != nil {
			t.Fatalf("say hello failed %v", err)
		}
		println(resp.Message)
		time.Sleep(100 * time.Millisecond)
	}

	time.Sleep(2 * time.Second)
}

type server struct {
	Port string
}

// SayHello implements helloworld.GreeterServer
func (s *server) SayHello(ctx context.Context, in *helloworld.HelloRequest) (*helloworld.HelloReply, error) {
	return &helloworld.HelloReply{Message: fmt.Sprintf("Hello From %s", s.Port)}, nil
}

func newServer(t *testing.T, port, version string, weight int64) {
	register := NewRegister(etcdAddrs)
	defer register.Stop()

	listen, err := net.Listen("tcp", port)
	if err != nil {
		logger.F("TestResolver", "net.Listen.err")
		return
	}

	s := grpc.NewServer()
	helloworld.RegisterGreeterServer(s, &server{Port: port})

	info := Server{
		Name:    "hello",
		Addr:    fmt.Sprintf("127.0.0.1%s", port),
		Version: version,
		Weight:  weight,
	}

	register.Register(context.Background(), info, 10)

	if err = s.Serve(listen); err != nil {
		logger.F("newServer", "fail.to.server, %+v", err)
	}
}

func TestEtcdResolver(t *testing.T) {
	r := NewResolver(etcdAddrs)
	resolver.Register(r)

	r.StartResolver(context.Background(), r)

	time.Sleep(5 * time.Second)

	/*conn, err := grpc.Dial("etcd:///hello", grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"round_robin"}`), grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		logger.F("aaa", "fail.to.dial: %+v", err)
	}
	defer conn.Close()

	c := helloworld.NewGreeterClient(conn)

	for i := 0; i < 10; i++ {
		resp, err := c.SayHello(context.Background(), &helloworld.HelloRequest{Name: "acc"})
		if err != nil {
			t.Fatalf("say hello failed %v", err)
		}
		println(resp.Message)
		time.Sleep(100 * time.Millisecond)
	}*/

	time.Sleep(2 * time.Second)
}

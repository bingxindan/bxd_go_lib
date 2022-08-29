package discovery

import (
	"context"
	"fmt"
	"github.com/bingxindan/bxd_go_lib/logger"
	"github.com/bingxindan/bxd_go_lib/register/etcd/discovery/helloworld"
	"google.golang.org/grpc"
	"net"
	"testing"
	"time"
)

func TestReg(t *testing.T) {
	var ctx context.Context
	adds := []string{"0.0.0.0:2379"}
	etcdRegister := NewRegister(adds)
	fmt.Printf("etcdRegister: %+v\n", etcdRegister)

	node := Server{
		Name:    "user",
		Addr:    "0.0.0.1:10091",
		Version: "1.0.0",
		Weight:  2,
	}
	if _, err := etcdRegister.Register(ctx, node, 10); err != nil {
		fmt.Printf("err: %+v\n", err)
		panic(fmt.Sprintf("server register failed: %v", err))
	}

	fmt.Printf("11111\n")
	time.Sleep(100 * time.Minute)
}

func TestDisRegister(t *testing.T) {

	/*go EtcdRegister(t, ":1001", "1.0", 1)
	go EtcdRegister(t, ":1002", "1.0", 2)
	go EtcdRegister(t, ":1003", "1.0", 3)
	go EtcdRegister(t, ":1004", "1.0", 4)
	go EtcdRegister(t, ":1005", "1.0", 5)
	go EtcdRegister(t, ":1006", "1.0", 6)
	go EtcdRegister(t, ":1007", "1.0", 7)
	go EtcdRegister(t, ":1008", "1.0", 8)
	go EtcdRegister(t, ":1009", "1.0", 1)
	go EtcdRegister(t, ":1010", "1.0", 1)
	go EtcdRegister(t, ":1011", "1.0", 1)*/

	go EtcdRegister(t, ":1012", "1.0", 1)
	go EtcdRegister(t, ":1013", "1.0", 2)
	go EtcdRegister(t, ":1014", "1.0", 3)
	go EtcdRegister(t, ":1015", "1.0", 4)
	go EtcdRegister(t, ":1016", "1.0", 5)
	go EtcdRegister(t, ":1017", "1.0", 6)
	go EtcdRegister(t, ":1018", "1.0", 7)
	go EtcdRegister(t, ":1019", "1.0", 8)
	go EtcdRegister(t, ":1020", "1.0", 1)
	go EtcdRegister(t, ":1021", "1.0", 1)
	go EtcdRegister(t, ":1022", "1.0", 1)

	//time.Sleep(9999999 * time.Minute)
	select {}
}

func EtcdRegister(t *testing.T, port, version string, weight int64) {
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

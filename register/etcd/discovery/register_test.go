package discovery

import (
	"fmt"
	"testing"
	"time"
)

func TestReg(t *testing.T) {
	adds := []string{"0.0.0.0:2379"}
	etcdRegister := NewRegister(adds)
	fmt.Printf("etcdRegister: %+v\n", etcdRegister)

	node := Server{
		Name:    "user",
		Addr:    "0.0.0.1:10091",
		Version: "1.0.0",
		Weight:  2,
	}
	if _, err := etcdRegister.Register(node, 10); err != nil {
		fmt.Printf("err: %+v\n", err)
		panic(fmt.Sprintf("server register failed: %v", err))
	}

	fmt.Printf("11111\n")
	time.Sleep(100 * time.Minute)
}

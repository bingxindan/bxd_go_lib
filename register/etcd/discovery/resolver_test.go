package discovery

import (
	"fmt"
	"testing"
	"time"
)

func TestReso(t *testing.T) {
	etcdAddrs := []string{"0.0.0.0:2379"}
	infoRes, err := StartResolver(etcdAddrs)
	fmt.Printf("infoRes: %+v, err: %+v\n", infoRes, err)
	time.Sleep(2 * time.Second)
}

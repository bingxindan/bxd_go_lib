package config

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func TestPrepareConfigs(t *testing.T) {
	ctx := context.Background()

	PrepareConfigs(ctx, nil)

	time.Sleep(3 * time.Second)

	s := String("server.http.addr")
	fmt.Println("sssssss", s)

	select {}
}

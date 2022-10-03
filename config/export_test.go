package config

import (
	"context"
	"fmt"
	"testing"
)

func TestExport(t *testing.T) {
	var (
		ctx       = context.Background()
		filePaths = []string{"config.yaml"}
	)
	config, err := NewConfig(ctx, filePaths)
	fmt.Printf("aaaaaaa, %+v, err: %+v\n", config, err)
}

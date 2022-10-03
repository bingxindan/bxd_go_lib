package config

import (
	"context"
	"fmt"
	"testing"
)

func TestFile(t *testing.T) {
	println("aaaa")

	var (
		ctx        = context.Background()
		path       = "config.yaml"
		typ        = "auto"
		autoReload = true
	)

	ret, err := NewFileConfigSource(ctx, path, typ, autoReload)

	fmt.Printf("bbbbbb: %+v, err: %+v\n", ret, err)
}

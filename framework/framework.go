package framework

import (
	"context"
	"github.com/bingxindan/bxd_go_lib/config"
	"sync"
)

var (
	configOnce sync.Once
)

func Init(ctx context.Context, ns []string) {
	configOnce.Do(func() {
		config.PrepareConfigs(ctx, ns)
	})
}

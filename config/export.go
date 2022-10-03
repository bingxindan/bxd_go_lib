package config

import (
	"context"
	"sync"
)

type ExportConfig Config
type ExportSource Source

var (
	mutex         sync.RWMutex
	defaultConfig ExportConfig
)

func NewConfig(ctx context.Context, ns []string, sources ...ExportSource) (ExportConfig, error) {
	var internalSources []Source

	for _, source := range sources {
		if source == nil {
			continue
		}
		internalSources = append(internalSources, source)
	}

	return NewConfigImpl(ctx, ns, internalSources...)
}

func SetGlobalConfig(config ExportConfig) {
	mutex.Lock()
	defer mutex.Unlock()

	defaultConfig = config
}

func GlobalConfig() ExportConfig {
	mutex.RLock()
	defer mutex.RUnlock()

	return defaultConfig
}

func String(key string) string {
	return GlobalConfig().String(key)
}

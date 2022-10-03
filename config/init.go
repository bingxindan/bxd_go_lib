package config

import (
	"context"
	"log"
)

func PrepareConfigs(ctx context.Context, ns []string) {
	fileConfigSource, err := NewFileConfigSource()
	if err != nil {
		return
	}

	bootstrapConfigPreparer := func(ExportConfig) (ExportConfig, error) {
		bootstrapConfig, err := NewConfig(ctx, ns, fileConfigSource)
		return bootstrapConfig, err
	}

	bootstrapConfig, err := bootstrapConfigPreparer(nil)
	if err != nil {
		log.Fatalf("[ERROR] Create bootstreap config failed, err = %+v", err)
	}

	configPreparer := func(cfg ExportConfig) (ExportConfig, error) {
		return bootstrapConfig, nil
	}

	defaultCfg, err := configPreparer(bootstrapConfig)
	if err != nil {
		log.Fatalf("[ERROR] Create config failed, err = %+v", err)
	}

	SetGlobalConfig(defaultCfg)
}

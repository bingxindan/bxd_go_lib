package framework

import (
	"github.com/bingxindan/bxd_go_lib/config"
	"time"
)

type Options struct {
	finalizeTimeout      time.Duration
	compactUsage         bool
	flagsToShow          []string
	flagsToHide          []string
	configPreparer       func(exportConfig config.ExportConfig) (config.ExportConfig, error)
	beforeConfigPreparer []func()
}

type Option func(*Options) error

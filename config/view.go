package config

import (
	"context"
	"github.com/bingxindan/bxd_go_lib/config/watcher"
	"strings"
	"sync"
)

type ValueGetter interface {
	Get(key string) (string, bool)
	GetByPrefix(prefix string) map[string]string
}

type DataSource interface {
	Get(key, prefix string, raw bool) (string, int)
	AppendPrefix(ctx context.Context, prefix []string) error
	GetByPrefix(prefix string) map[string]string
}

type View struct {
	Value

	filter   *Filter
	source   DataSource
	prefix   string
	ns       []string
	password [][16]byte
	notifier *watcher.Notifier
	once     sync.Once

	mu        sync.RWMutex
	viewCache map[string]*View
	cancel    func()
}

func NewConfigView(ctx context.Context, source DataSource, prefix string, ns []string, password [][16]byte) *View {
	lowerNs := make([]string, len(ns))
	for i, n := range ns {
		lowerNs[i] = strings.ToLower(n)
	}

	if len(lowerNs) != 0 {
		_ = source.AppendPrefix(ctx, lowerNs)
	}

	filter := NewFilter(source, prefix, ns, password)
	return &View{
		Value:     NewValue(filter),
		filter:    filter,
		source:    source,
		prefix:    prefix,
		ns:        ns,
		notifier:  watcher.NewNotifier(),
		viewCache: map[string]*View{},
	}
}

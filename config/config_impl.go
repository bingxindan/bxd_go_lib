package config

import (
	"context"
	"github.com/bingxindan/bxd_go_lib/config/watcher"
	"github.com/bingxindan/bxd_go_lib/tools/helper"
	"log"
	"strings"
	"sync"
	"time"
)

type ValueNode struct {
	level int
	value string
}

type ConfigImpl struct {
	source   []*SourceWrapper
	notifier *watcher.Notifier
	cache    *PrefixTree
	closed   uint32

	mu     sync.RWMutex
	cancel []func()
}

func NewConfigImpl(ctx context.Context, ns []string, sources ...Source) (config Config, err error) {
	cfg := &ConfigImpl{
		notifier: watcher.NewNotifier(),
		cache:    NewPrefixTree("."),
		closed:   0,
	}

	if len(ns) == 0 {
		ns = []string{""}
	}

	for idx, source := range sources {
		sourceWrapper, err := NewSourceWrapper(ctx, source)
		if err != nil {
			return nil, err
		}
		cfg.cancel = append(cfg.cancel, sourceWrapper.Watch(ctx, &WatcherClosure{idx, cfg.onUpdate, cfg.onDelete, cfg.onSync}))
		cfg.source = append(cfg.source, sourceWrapper)
	}

	defer func() {
		if e := cfg.sync(ctx); e != nil {
			err = e
		}
	}()

	return NewConfigView(ctx, cfg, "", ns, nil), nil
}

func (s *ConfigImpl) sync(ctx context.Context) error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for level := range s.source {
		if val, err := s.source[level].Sync(ctx); err != nil {
			return err
		} else {
			notify := map[string]string{}
			for k, v := range val {
				k = strings.ToLower(k)
				if valNode, ok := s.cache.Get(k).(ValueNode); !ok || valNode.value != v {
					notify[k] = v
				}
				s.cache.Put(k, ValueNode{level, v})
			}

			s.notifier.OnUpdate(notify)
		}
	}

	return nil
}

func (s *ConfigImpl) AppendPrefix(ctx context.Context, prefix []string) error {
	defer helper.TimeoutGuardWithFunc(ctx, 30*time.Second, nil, nil, func(...interface{}) {
		log.Fatalf("Append.config.prefix.with.prefix: %v", prefix)
	})()

	s.mu.RLock()
	defer s.mu.RUnlock()

	var err error
	for _, source := range s.source {
		if e := source.AppendPrefix(prefix); e != nil {
			err = e
		}
	}

	return err
}

func (s *ConfigImpl) Get(key, prefix string, raw bool) (string, int) {
	key = strings.ToLower(key)

	if raw {
		raw, hasRaw := s.cache.Get(key).(ValueNode)
		if !hasRaw || !s.source[raw.level].IgnoreNamespace() {
			return "", -1
		} else {
			return raw.value, raw.level
		}
	}

	v, hasV := s.cache.Get(s.withPrefix(key, prefix)).(ValueNode)
	if !hasV {
		return "", -1
	} else {
		return v.value, v.level
	}
}

func (s *ConfigImpl) withPrefix(key, prefix string) string {
	return prefix + key
}

func (s *ConfigImpl) GetByPrefix(prefix string) map[string]string {
	result := map[string]string{}

	for k, v := range s.cache.SearchByPrefix(strings.ToLower(prefix)) {
		if v, ok := v.(ValueNode); ok {
			result[k] = v.value
		}
	}

	return result
}

func (s *ConfigImpl) onUpdate(level int, val map[string]string) {
	notify := map[string]string{}

	for k, v := range val {
		k = strings.ToLower(k)
		if valNode, ok := s.cache.Get(k).(ValueNode); !ok || valNode.level <= level {
			s.cache.Put(k, ValueNode{level, v})
			if !ok || valNode.value != v {
				notify[k] = v
			}
		}
	}

	s.notifier.OnUpdate(notify)
}

func (s *ConfigImpl) onDelete(level int, val []string) {
	updateNotify := map[string]string{}
	var deleteNotify []string

	s.mu.RLock()
	defer s.mu.RUnlock()

out:
	for _, k := range val {
		k = strings.ToLower(k)
		if valNode, ok := s.cache.Get(k).(ValueNode); ok && valNode.level == level {
			for i := level - 1; i >= 0; i-- {
				if v, ok := s.source[i].Get(k); ok {
					s.cache.Put(k, ValueNode{i, v})
					if v != valNode.value {
						updateNotify[k] = v
					}
					continue out
				}
			}
			s.cache.Del(k)
			deleteNotify = append(deleteNotify, k)
		}
	}

	s.notifier.OnUpdate(updateNotify)
	s.notifier.OnDelete(deleteNotify)
}

func (s *ConfigImpl) onSync(level int, val map[string]string) {
	updateNotify := map[string]string{}
	var deleteNotify []string

	lowerVal := make(map[string]string, len(val))
	for k, v := range val {
		lowerVal[strings.ToLower(k)] = v
	}

	for k, v := range lowerVal {
		if valNode, ok := s.cache.Get(k).(ValueNode); !ok || valNode.level <= level {
			s.cache.Put(k, ValueNode{level, v})
			if !ok || valNode.value != v {
				updateNotify[k] = v
			}
		}
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

out:
	for _, k := range s.cache.Keys() {
		k = strings.ToLower(k)
		valNode, _ := s.cache.Get(k).(ValueNode)
		if _, ok := lowerVal[k]; !ok && level == valNode.level {
			for i := level - 1; i >= 0; i-- {
				if v, ok := s.source[i].Get(k); ok {
					s.cache.Put(k, ValueNode{i, v})
					if v != valNode.value {
						updateNotify[k] = v
					}
					continue out
				}
			}
			s.cache.Del(k)
			deleteNotify = append(deleteNotify, k)
		}
	}

	s.notifier.OnUpdate(updateNotify)
	s.notifier.OnDelete(deleteNotify)
}

type WatcherClosure struct {
	level int

	OnUpdateFunc func(int, map[string]string)
	OnDeleteFunc func(int, []string)
	OnSyncFunc   func(int, map[string]string)
}

func (s *WatcherClosure) OnUpdate(val map[string]string) {
	s.OnUpdateFunc(s.level, val)
}

func (s *WatcherClosure) OnDelete(val []string) {
	s.OnDeleteFunc(s.level, val)
}

func (s *WatcherClosure) OnSync(val map[string]string) {
	s.OnSyncFunc(s.level, val)
}

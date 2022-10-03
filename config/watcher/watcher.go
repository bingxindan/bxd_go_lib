package watcher

import "strings"

type Watcher interface {
	OnUpdate(map[string]string)
	OnDelete([]string)
}

type WatcherHelper struct {
	keyToWatch  string
	watcherFunc func(v string, deleted bool)
}

type SourceWatcher interface {
	Watcher
	OnSync(map[string]string)
}

func NewWatcherHelper(keyToWatch string, watcher func(v string, deleted bool)) *WatcherHelper {
	return &WatcherHelper{
		keyToWatch:  keyToWatch,
		watcherFunc: watcher,
	}
}

func (w *WatcherHelper) OnUpdate(val map[string]string) {
	if v, ok := val[strings.ToLower(w.keyToWatch)]; ok {
		w.watcherFunc(v, false)
	}
}

func (w *WatcherHelper) OnDelete(val []string) {
	for _, k := range val {
		if strings.ToLower(k) != w.keyToWatch {
			continue
		}
		w.watcherFunc("", true)
		return
	}
}

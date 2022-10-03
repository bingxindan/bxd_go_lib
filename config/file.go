package config

import (
	"context"
	"flag"
	"github.com/bingxindan/bxd_go_lib/config/watcher"
	"github.com/fsnotify/fsnotify"
	"github.com/pkg/errors"
	"log"
	"path/filepath"
	"sync"
	"time"
)

type FileConfigSource struct {
	path      string
	reader    FileReader
	notifier  *watcher.Notifier
	watchDone chan struct{}
	reload    bool

	once sync.Once
}

var (
	// flagconf is the config flag.
	flagconf       string
	flagtyp        string
	flagAutoReload bool

	readerFactory = map[string]func(path string) (FileReader, error){
		"yaml": NewYamlReader,
	}
)

func init() {
	flag.StringVar(&flagconf, "conf", "", "config path, eg: -conf config.yaml")
	flag.StringVar(&flagtyp, "typ", "", "typ, eg: -typ auto")
	flag.BoolVar(&flagAutoReload, "autoReload", true, "autoReload, eg: -autoReload true")
}

func NewFileConfigSource() (*FileConfigSource, error) {
	flag.Parse()

	if flagtyp == "" {
		flagtyp = "auto"
	}

	if reader, err := createReader(flagconf, flagtyp); err != nil {
		return nil, err
	} else {
		return &FileConfigSource{
			path:      flagconf,
			reader:    reader,
			notifier:  watcher.NewNotifier(),
			watchDone: make(chan struct{}),
			reload:    flagAutoReload,
		}, nil
	}
}

func createReader(filePath, fileType string) (FileReader, error) {
	if fileType == "auto" {
		fileType = filepath.Ext(filePath)
		if fileType != "" && fileType[0] == '.' {
			fileType = fileType[1:]
		}
	}

	if factory, ok := readerFactory[fileType]; ok {
		return factory(filePath)
	} else {
		return nil, errors.Errorf("unknown file type, type = %s", fileType)
	}
}

func (s *FileConfigSource) Sync(ctx context.Context) (map[string]string, error) {
	return s.reader.Read(ctx)
}

func (s *FileConfigSource) Watch(ctx context.Context, watcher watcher.SourceWatcher) (cancel func()) {
	if !s.reload {
		return
	}

	ctxW, cancelFunc := context.WithCancel(ctx)
	cancelWatch := s.notifier.Watch(ctxW, watcher)

	cancel = func() {
		cancelFunc()
		cancelWatch()
	}

	s.once.Do(func() {
		go func() {
			w, err := fsnotify.NewWatcher()
			if err != nil {
				log.Printf("Register fsnotify watcher failed, err = %v\n", err)
				return
			}

			defer func() {
				_ = w.Close()
			}()

			err = w.Add(s.path)
			if err != nil {
				log.Printf("Add config path to fsnotify watcher failed, path = %s, err = %v\n", s.path, err)
				return
			}

			go func() {
				if fileProvider, ok := s.reader.(AdditionalFileProvider); !ok {
					return
				} else {
					for {
						select {
						case v := <-fileProvider.Watch():
							if v.Canceled {
								return
							}

							if err := w.Add(v.FilePath); err != nil {
								log.Printf("Add config path to fsnotify watcher failed, path = %s, err = %v\n", v.FilePath, err)
							}
						}
					}
				}
			}()

			for {
				select {
				case event, ok := <-w.Events:
					if !ok {
						log.Printf("Bump from watcher.Events failed, err = %v\n", err)
						return
					}

					if event.Op&fsnotify.Rename == fsnotify.Rename {
						time.Sleep(time.Second)
						err = w.Add(event.Name)
						if err != nil {
							log.Printf("Add config path to fsnotify watcher failed, path = %s, err = %v\n", s.path, err)
							return
						}
					}

					go func() {
						if (event.Op&fsnotify.Write == fsnotify.Write) || (event.Op&fsnotify.Rename == fsnotify.Rename) {
							if val, err := s.Sync(ctxW); err == nil {
								s.notifier.OnSync(val)
							}
						}
					}()
				case err, ok := <-w.Errors:
					if !ok {
						log.Printf("Bump from watcher.Errors failed, err = %v\n", err)
						return
					}
				case <-s.watchDone:
					return
				case <-ctxW.Done():
					return
				}
			}
		}()
	})

	return cancel
}

func (s *FileConfigSource) AppendPrefix(prefix []string) error {
	return nil
}

func (s *FileConfigSource) Close() error {
	return nil
}

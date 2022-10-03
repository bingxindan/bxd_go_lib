package config

import "context"

type FileReader interface {
	Read(ctx context.Context) (map[string]string, error)
	Close() error
}

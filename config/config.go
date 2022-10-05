package config

import "time"

type Value interface {
	Bool(key string) bool
	Int(key string) int
	Uint(key string) uint
	Int64(key string) int64
	Uint64(key string) uint64
	Float64(key string) float64
	Duration(key string) time.Duration
	Bytes(key string) []byte
	String(key string) string
	IntList(key string) []int
	UintList(key string) []uint
	Int64List(key string) []int64
	Uint64List(key string) []uint64
	StringList(key string) []string
	StringMap(key string) map[string]string
	GetByPrefix(prefix string) map[string]string
	Has(key string) bool
}

type Config interface {
	Value
}

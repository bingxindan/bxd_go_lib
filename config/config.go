package config

type Value interface {
	String(key string) string
}

type Config interface {
	Value
}

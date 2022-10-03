package config

type ValueImpl struct {
	ValueGetter
}

func NewValue(getter ValueGetter) *ValueImpl {
	return &ValueImpl{getter}
}

func (s *ValueImpl) String(key string) string {
	if v, ok := s.Get(key); ok {
		return v
	} else {
		return ""
	}
}

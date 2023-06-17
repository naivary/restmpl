package logging

import "github.com/knadh/koanf/v2"

func NewEnvManager(k *koanf.Koanf) *manager {
	m := newManager()
	attrs := make([]any, 0)
	attrs = append(attrs, commonEnvAttrs(k))
	m.AddCommonAttrs(attrs)
	return m
}

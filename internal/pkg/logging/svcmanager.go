package logging

import "github.com/naivary/restmpl/internal/pkg/service"

func NewSvcManager(svc service.Service) *manager {
	m := newManager()
	attrs := make([]any, 0)
	attrs = append(attrs, commonSvcAttrs(svc))
	m.AddCommonAttrs(attrs)
	return m
}

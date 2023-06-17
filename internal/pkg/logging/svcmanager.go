package logging

import "github.com/naivary/restmpl/internal/pkg/service"

func NewSvcManager(svc service.Service) *manager {
	m := NewManager()
	m.AddCommonAttrs(commonSvcAttrs(svc))
	return m
}

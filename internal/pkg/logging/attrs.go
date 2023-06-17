package logging

import (
	"github.com/naivary/restmpl/internal/pkg/service"
	"golang.org/x/exp/slog"
)

func commonSvcAttrs(svc service.Service) slog.Attr {
	return slog.Group(
		"service",
		slog.String("id", svc.ID()),
		slog.String("name", svc.Name()),
	)
}

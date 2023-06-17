package logging

import (
	"github.com/knadh/koanf/v2"
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

func commonEnvAttrs(k *koanf.Koanf) slog.Attr {
	return slog.Group(
		"env",
		slog.String("version", k.String("version")),
		slog.String("name", k.String("name")),
	)
}

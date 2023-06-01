package builder

import (
	"context"
	"net/http"
	"time"

	"github.com/knadh/koanf/v2"
	"github.com/naivary/instance/internal/pkg/service"
	"golang.org/x/exp/slog"
)

var _ Recorder = (*envBuilder)(nil)

type envBuilder struct {
	ctx context.Context
	rec slog.Record
}

func NewEnvBuilder(ctx context.Context, level slog.Leveler, msg string) *envBuilder {
	return &envBuilder{
		ctx: ctx,
		rec: slog.NewRecord(time.Now(), level.Level(), msg, 0),
	}
}

func (e *envBuilder) Data() (context.Context, slog.Record) {
	return e.ctx, e.rec
}

func (e *envBuilder) APIServerStart(k *koanf.Koanf, srv *http.Server) {
	api := slog.Group(
		"api",
		slog.String("name", k.String("name")),
		slog.String("version", k.String("version")),
		slog.String("used_config_file", k.String("cfgFile")),
	)

	srvConfig := slog.Group(
		"server",
		slog.String("addr", srv.Addr),
	)

	e.rec.AddAttrs(api, srvConfig)
}

func (e *envBuilder) ServiceShutdown(svc service.Service) {
	svcShutdown := slog.Group(
		"service",
		slog.String("name", svc.Name()),
		slog.String("id", svc.ID()),
	)
	e.rec.AddAttrs(svcShutdown)
}

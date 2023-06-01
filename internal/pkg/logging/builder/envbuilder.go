package builder

import (
	"context"
	"net/http"
	"time"

	"github.com/knadh/koanf/v2"
	"github.com/naivary/instance/internal/pkg/service"
	"golang.org/x/exp/slog"
)

var _ Recorder = (*EnvBuilder)(nil)

type EnvBuilder struct {
	ctx context.Context
	rec slog.Record
}

func NewEnvBuilder(ctx context.Context, level slog.Leveler, msg string) *EnvBuilder {
	return &EnvBuilder{
		ctx: ctx,
		rec: slog.NewRecord(time.Now(), level.Level(), msg, 0),
	}
}

func (e *EnvBuilder) Data() (context.Context, slog.Record) {
	return e.ctx, e.rec
}

func (e *EnvBuilder) APIServerStart(k *koanf.Koanf, srv *http.Server) *EnvBuilder {
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
	return e
}

func (e *EnvBuilder) ServiceShutdown(svc service.Service) *EnvBuilder {
	svcShutdown := slog.Group(
		"service",
		slog.String("name", svc.Name()),
		slog.String("id", svc.ID()),
	)
	e.rec.AddAttrs(svcShutdown)
	return e
}

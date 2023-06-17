package builder

import (
	"context"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/knadh/koanf/v2"
	"github.com/naivary/restmpl/internal/pkg/service"
	"golang.org/x/exp/slog"
)

type Recorder interface {
	Data() (context.Context, slog.Record)
}

var _ Recorder = (*Builder)(nil)

type Builder struct {
	ctx    context.Context
	record slog.Record
}

func New(ctx context.Context, level slog.Leveler, msg string) *Builder {
	return &Builder{
		ctx:    ctx,
		record: slog.NewRecord(time.Now(), level.Level(), msg, 0),
	}
}

func (b *Builder) Data() (context.Context, slog.Record) {
	return b.ctx, b.record
}

func (b *Builder) IncomingRequest(req *http.Request) *Builder {
	id := middleware.GetReqID(req.Context())
	attr := slog.Group(
		"request",
		slog.String("id", id),
		slog.String("method", req.Method),
		slog.String("host", req.Host),
		slog.String("remote_addr", req.RemoteAddr),
		slog.String("endpoint", req.URL.Path),
		slog.String("protocol_version", req.Proto),
		slog.String("user_agent", req.UserAgent()),
	)
	b.record.AddAttrs(attr)
	return b
}

func (b *Builder) APIServerStart(k *koanf.Koanf, srv *http.Server) *Builder {
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

	b.record.AddAttrs(api, srvConfig)
	return b
}

func (b *Builder) ServiceInfo(svc service.Service) *Builder {
	svcShutdown := slog.Group(
		"service",
		slog.String("name", svc.Name()),
		slog.String("id", svc.ID()),
	)
	b.record.AddAttrs(svcShutdown)
	return b
}

func (b *Builder) ServiceInit(svc service.Service) *Builder {
	attr := slog.Group(
		"service",
		slog.String("name", svc.Name()),
		slog.String("id", svc.ID()),
		slog.String("endpoint", svc.Pattern()),
	)
	b.record.AddAttrs(attr)
	return b
}

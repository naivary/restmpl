package builder

import (
	"context"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/knadh/koanf/v2"
	"golang.org/x/exp/slog"
)

type Recorder interface {
	Data() (context.Context, slog.Record)
}

var _ Recorder = (*record)(nil)

type record struct {
	slogRecord slog.Record
	ctx        context.Context
}

func New(ctx context.Context, level slog.Leveler, msg string) record {
	return record{
		ctx:        ctx,
		slogRecord: slog.NewRecord(time.Now(), level.Level(), msg, 0),
	}
}

func (r record) Data() (context.Context, slog.Record) {
	return r.ctx, r.slogRecord
}

func (r *record) IncomingRequest(req *http.Request) {
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
	r.slogRecord.AddAttrs(attr)
}

func (r *record) APIServerStart(k *koanf.Koanf, srv *http.Server) {
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

	r.slogRecord.AddAttrs(api, srvConfig)

}

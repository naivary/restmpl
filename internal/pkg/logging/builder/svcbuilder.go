package builder

import (
	"context"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"golang.org/x/exp/slog"
)

var _ Recorder = (*svcRecord)(nil)

type svcRecord struct {
	rec slog.Record
	ctx context.Context
}

func NewSvcBuilder(ctx context.Context, level slog.Leveler, msg string) *svcRecord {
	return &svcRecord{
		ctx: ctx,
		rec: slog.NewRecord(time.Now(), level.Level(), msg, 0),
	}
}

func (r svcRecord) Data() (context.Context, slog.Record) {
	return r.ctx, r.rec
}

func (r *svcRecord) IncomingRequest(req *http.Request) {
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
	r.rec.AddAttrs(attr)
}

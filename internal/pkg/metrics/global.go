package metrics

import (
	"github.com/naivary/restmpl/internal/pkg/service"
	"github.com/prometheus/client_golang/prometheus"
)

func IncomingHTTPRequest(svc service.Service) prometheus.Counter {
	return prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "incoming_http_requests",
			Help: "number of incoming request (including 4XX and 5XX)",
			ConstLabels: prometheus.Labels{
				"name":        svc.Name(),
				"id":          svc.ID(),
				"endpoint":    svc.Pattern(),
				"description": svc.Description(),
			},
		},
	)
}

func NumberOfErrors(svc service.Service) prometheus.Counter {
	return prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "incoming_unsuccessfull_http_requests",
			Help: "number of request which finished with 4XX or 5XX",
			ConstLabels: prometheus.Labels{
				"id": svc.ID(),
			},
		},
	)
}

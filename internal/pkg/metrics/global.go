package metrics

import (
	"fmt"

	"github.com/naivary/restmpl/internal/pkg/service"
	"github.com/prometheus/client_golang/prometheus"
)

func IncomingHTTPRequest(svc service.Service) prometheus.Counter {
	help := fmt.Sprintf("number of incoming request for all the subroutes of %s", svc.Pattern())
	name := fmt.Sprintf("incoming_http_request_%s", svc.Name())
	return prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: name,
			Help: help,
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
	help := fmt.Sprintf("number of errors that appeared in the service %s", svc.Name())
	name := fmt.Sprintf("err_counter_%s", svc.Name())
	return prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: name,
			Help: help,
			ConstLabels: prometheus.Labels{
				"id": svc.ID(),
			},
		},
	)
}

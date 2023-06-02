package metrics

import "github.com/prometheus/client_golang/prometheus"

func (m manager) GetCounter(name string) prometheus.Counter {
	return m.counters[name]
}

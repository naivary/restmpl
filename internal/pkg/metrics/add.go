package metrics

import "github.com/prometheus/client_golang/prometheus"

func (m manager) AddCounter(name string, c prometheus.Counter) {
	m.counters[name] = c
}

func (m manager) AddCounterVec(name string, cVec prometheus.CounterVec) {
	m.counterVecs[name] = cVec
}

func (m manager) AddGauge(name string, g prometheus.Gauge) {
	m.gauges[name] = g
}

func (m manager) AddGaugeVec(name string, gVec prometheus.GaugeVec) {
	m.gaugeVecs[name] = gVec
}

func (m manager) AddSummary(name string, s prometheus.Summary) {
	m.summaries[name] = s
}

func (m manager) AddSummaryVec(name string, sVec prometheus.SummaryVec) {
	m.summaryVecs[name] = sVec
}

func (m manager) AddHistogram(name string, h prometheus.Histogram) {
	m.histograms[name] = h
}

func (m manager) AddHistrogramVec(name string, hVec prometheus.HistogramVec) {
	m.histogramVecs[name] = hVec
}

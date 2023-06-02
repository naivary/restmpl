package metrics

import "github.com/prometheus/client_golang/prometheus"

func (m manager) GetCounter(name string) prometheus.Counter {
	return m.counters[name]
}

func (m manager) GetCounterVec(name string) prometheus.CounterVec {
	return m.counterVecs[name]
}

func (m manager) GetGauge(name string) prometheus.Gauge {
	return m.gauges[name]
}

func (m manager) GetGaugeVec(name string) prometheus.GaugeVec {
	return m.gaugeVecs[name]
}

func (m manager) GetSummary(name string) prometheus.Summary {
	return m.summaries[name]
}

func (m manager) GetSummaryVec(name string) prometheus.SummaryVec {
	return m.summaryVecs[name]
}

func (m manager) GetHistogram(name string) prometheus.Histogram {
	return m.histograms[name]
}

func (m manager) GetHistogramVec(name string) prometheus.HistogramVec {
	return m.histogramVecs[name]
}

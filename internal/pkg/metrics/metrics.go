package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

type Managee interface {
	AddCounter(name string, c prometheus.Counter)
	AddCounterVec(name string, cVec prometheus.CounterVec)
	AddGauge(name string, g prometheus.Gauge)
	AddGaugeVec(name string, gVec prometheus.GaugeVec)
	AddSummary(name string, s prometheus.Summary)
	AddSummaryVec(name string, sVec prometheus.SummaryVec)
	AddHistogram(name string, h prometheus.Histogram)
	AddHistogramVec(name string, hVec prometheus.HistogramVec)
	GetCounter(name string) prometheus.Counter
	GetCounterVec(name string) prometheus.CounterVec
	GetGauge(name string) prometheus.Gauge
	GetGaugeVec(name string) prometheus.GaugeVec
	GetSummary(name string) prometheus.Summary
	GetSummaryVec(name string) prometheus.SummaryVec
	GetHistogram(name string) prometheus.Histogram
	GetHistogramVec(name string) prometheus.HistogramVec
	All() []prometheus.Collector
}

type Manager interface {
	Managee
	Register(...prometheus.Collector) error
	Registry() *prometheus.Registry
}

var _ Manager = (*manager)(nil)

type manager struct {
	re            *prometheus.Registry
	counterVecs   map[string]prometheus.CounterVec
	counters      map[string]prometheus.Counter
	gauges        map[string]prometheus.Gauge
	gaugeVecs     map[string]prometheus.GaugeVec
	histograms    map[string]prometheus.Histogram
	histogramVecs map[string]prometheus.HistogramVec
	summaries     map[string]prometheus.Summary
	summaryVecs   map[string]prometheus.SummaryVec
}

func NewManagee() Managee {
	return &manager{
		re:            prometheus.NewRegistry(),
		counterVecs:   make(map[string]prometheus.CounterVec),
		counters:      make(map[string]prometheus.Counter),
		gauges:        make(map[string]prometheus.Gauge),
		gaugeVecs:     make(map[string]prometheus.GaugeVec),
		histograms:    make(map[string]prometheus.Histogram),
		histogramVecs: make(map[string]prometheus.HistogramVec),
		summaries:     make(map[string]prometheus.Summary),
		summaryVecs:   make(map[string]prometheus.SummaryVec),
	}
}

func NewManager() Manager {
	return &manager{
		re:            prometheus.NewRegistry(),
		counterVecs:   make(map[string]prometheus.CounterVec),
		counters:      make(map[string]prometheus.Counter),
		gauges:        make(map[string]prometheus.Gauge),
		gaugeVecs:     make(map[string]prometheus.GaugeVec),
		histograms:    make(map[string]prometheus.Histogram),
		histogramVecs: make(map[string]prometheus.HistogramVec),
		summaries:     make(map[string]prometheus.Summary),
		summaryVecs:   make(map[string]prometheus.SummaryVec),
	}
}

func (m manager) Register(metrics ...prometheus.Collector) error {
	for _, metric := range metrics {
		if err := m.re.Register(metric); err != nil {
			return err
		}
	}
	return nil
}

func (m manager) Registry() *prometheus.Registry {
	return m.re
}

// TODO(naivary): probably baddest solution
func (m manager) All() []prometheus.Collector {
	met := make([]prometheus.Collector, 0)
	for _, c := range m.counters {
		met = append(met, c)
	}
	for _, cVec := range m.counterVecs {
		met = append(met, cVec)
	}
	for _, g := range m.gauges {
		met = append(met, g)
	}
	for _, gVec := range m.gaugeVecs {
		met = append(met, gVec)
	}
	for _, s := range m.summaries {
		met = append(met, s)
	}
	for _, sVec := range m.summaryVecs {
		met = append(met, sVec)
	}
	for _, h := range m.histograms {
		met = append(met, h)
	}
	for _, hVec := range m.histogramVecs {
		met = append(met, hVec)
	}
	return met
}

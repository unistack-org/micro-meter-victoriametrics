package victoriametrics

import (
	"io"
	"time"

	"go.unistack.org/metrics"
	"go.unistack.org/micro/v3/meter"
)

type victoriametricsMeter struct {
	set              *metrics.Set
	opts             meter.Options
	prometheusCompat bool
}

type prometheusCompatKey struct{}

func PrometheusCompat(b bool) meter.Option {
	return meter.SetOption(prometheusCompatKey{}, b)
}

func NewMeter(opts ...meter.Option) meter.Meter {
	m := &victoriametricsMeter{set: metrics.NewSet(), opts: meter.NewOptions(opts...)}
	if v, ok := m.opts.Context.Value(prometheusCompatKey{}).(bool); ok {
		m.prometheusCompat = v
	}
	return m
}

func (m *victoriametricsMeter) Name() string {
	return m.opts.Name
}

func (m *victoriametricsMeter) Clone(opts ...meter.Option) meter.Meter {
	options := m.opts
	for _, o := range opts {
		o(&options)
	}
	nm := &victoriametricsMeter{set: m.set, opts: options, prometheusCompat: m.prometheusCompat}
	if v, ok := m.opts.Context.Value(prometheusCompatKey{}).(bool); ok {
		m.prometheusCompat = v
	}
	return nm
}

func (m *victoriametricsMeter) buildName(name string, labels ...string) string {
	nl := len(m.opts.Labels) + len(labels)
	if nl == 0 {
		return name
	}

	nlabels := make([]string, 0, nl)
	nlabels = append(nlabels, m.opts.Labels...)
	nlabels = append(nlabels, labels...)

	return meter.BuildName(name, nlabels...)
}

func (m *victoriametricsMeter) Counter(name string, labels ...string) meter.Counter {
	return m.set.GetOrCreateCounter(m.buildName(name, labels...))
}

func (m *victoriametricsMeter) FloatCounter(name string, labels ...string) meter.FloatCounter {
	return m.set.GetOrCreateFloatCounter(m.buildName(name, labels...))
}

func (m *victoriametricsMeter) Gauge(name string, f func() float64, labels ...string) meter.Gauge {
	return m.set.GetOrCreateGauge(m.buildName(name, labels...), f)
}

func (m *victoriametricsMeter) Histogram(name string, labels ...string) meter.Histogram {
	if m.prometheusCompat {
		return m.set.GetOrCreateCompatibleHistogram(m.buildName(name, labels...))
	}
	return m.set.GetOrCreateHistogram(m.buildName(name, labels...))
}

func (m *victoriametricsMeter) Summary(name string, labels ...string) meter.Summary {
	return m.set.GetOrCreateSummary(m.buildName(name, labels...))
}

func (m *victoriametricsMeter) SummaryExt(name string, window time.Duration, quantiles []float64, labels ...string) meter.Summary {
	return m.set.GetOrCreateSummaryExt(m.buildName(name, labels...), window, quantiles)
}

func (m *victoriametricsMeter) Set(opts ...meter.Option) meter.Meter {
	nm := &victoriametricsMeter{opts: m.opts}
	for _, o := range opts {
		o(&nm.opts)
	}
	nm.set = metrics.NewSet()
	if v, ok := nm.opts.Context.Value(prometheusCompatKey{}).(bool); ok {
		nm.prometheusCompat = v
	}
	return nm
}

func (m *victoriametricsMeter) Init(opts ...meter.Option) error {
	for _, o := range opts {
		o(&m.opts)
	}
	if v, ok := m.opts.Context.Value(prometheusCompatKey{}).(bool); ok {
		m.prometheusCompat = v
	}
	return nil
}

func (m *victoriametricsMeter) Write(w io.Writer, opts ...meter.Option) error {
	options := m.opts
	for _, o := range opts {
		o(&options)
	}

	m.set.WritePrometheus(w)
	if options.WriteProcessMetrics {
		metrics.WriteProcessMetrics(w)
	}
	if options.WriteFDMetrics {
		metrics.WriteFDMetrics(w)
	}

	return nil
}

func (m *victoriametricsMeter) Options() meter.Options {
	return m.opts
}

func (m *victoriametricsMeter) String() string {
	return "victoriametrics"
}

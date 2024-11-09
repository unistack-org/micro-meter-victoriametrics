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
	if v, ok := m.opts.Context.Value(prometheusCompatKey{}).(bool); ok && v {
		m.prometheusCompat = true
	}
	return m
}

func (r *victoriametricsMeter) Name() string {
	return r.opts.Name
}

func (r *victoriametricsMeter) Clone(opts ...meter.Option) meter.Meter {
	options := r.opts
	for _, o := range opts {
		o(&options)
	}
	return &victoriametricsMeter{set: r.set, opts: options}
}

func (r *victoriametricsMeter) buildName(name string, labels ...string) string {
	nl := len(r.opts.Labels) + len(labels)
	if nl == 0 {
		return name
	}

	nlabels := make([]string, 0, nl)
	nlabels = append(nlabels, r.opts.Labels...)
	nlabels = append(nlabels, labels...)

	return meter.BuildName(name, nlabels...)
}

func (r *victoriametricsMeter) Counter(name string, labels ...string) meter.Counter {
	return r.set.GetOrCreateCounter(r.buildName(name, labels...))
}

func (r *victoriametricsMeter) FloatCounter(name string, labels ...string) meter.FloatCounter {
	return r.set.GetOrCreateFloatCounter(r.buildName(name, labels...))
}

func (r *victoriametricsMeter) Gauge(name string, f func() float64, labels ...string) meter.Gauge {
	return r.set.GetOrCreateGauge(r.buildName(name, labels...), f)
}

func (r *victoriametricsMeter) Histogram(name string, labels ...string) meter.Histogram {
	if r.prometheusCompat {
		return r.set.GetOrCreateCompatibleHistogram(r.buildName(name, labels...))
	}
	return r.set.GetOrCreateHistogram(r.buildName(name, labels...))
}

func (r *victoriametricsMeter) Summary(name string, labels ...string) meter.Summary {
	return r.set.GetOrCreateSummary(r.buildName(name, labels...))
}

func (r *victoriametricsMeter) SummaryExt(name string, window time.Duration, quantiles []float64, labels ...string) meter.Summary {
	return r.set.GetOrCreateSummaryExt(r.buildName(name, labels...), window, quantiles)
}

func (r *victoriametricsMeter) Set(opts ...meter.Option) meter.Meter {
	m := &victoriametricsMeter{opts: r.opts}
	for _, o := range opts {
		o(&m.opts)
	}
	m.set = metrics.NewSet()
	return m
}

func (r *victoriametricsMeter) Init(opts ...meter.Option) error {
	for _, o := range opts {
		o(&r.opts)
	}

	return nil
}

func (r *victoriametricsMeter) Write(w io.Writer, opts ...meter.Option) error {
	options := r.opts
	for _, o := range opts {
		o(&options)
	}

	r.set.WritePrometheus(w)
	if options.WriteProcessMetrics {
		metrics.WriteProcessMetrics(w)
	}
	if options.WriteFDMetrics {
		metrics.WriteFDMetrics(w)
	}

	return nil
}

func (r *victoriametricsMeter) Options() meter.Options {
	return r.opts
}

func (r *victoriametricsMeter) String() string {
	return "victoriametrics"
}

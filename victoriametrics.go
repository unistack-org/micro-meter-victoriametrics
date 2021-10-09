package victoriametrics

import (
	"io"
	"time"

	"github.com/VictoriaMetrics/metrics"
	"go.unistack.org/micro/v3/meter"
)

type victoriametricsMeter struct {
	set  *metrics.Set
	opts meter.Options
}

func NewMeter(opts ...meter.Option) meter.Meter {
	return &victoriametricsMeter{set: metrics.NewSet(), opts: meter.NewOptions(opts...)}
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
	if len(r.opts.MetricPrefix) > 0 {
		name = r.opts.MetricPrefix + name
	}

	if len(labels) == 0 {
		return name
	}

	if len(r.opts.LabelPrefix) == 0 {
		return meter.BuildName(name, labels...)
	}

	nlabels := make([]string, len(labels))
	copy(nlabels, labels)
	for idx := 0; idx <= len(nlabels)/2; idx += 2 {
		nlabels[idx] = r.opts.LabelPrefix + nlabels[idx]
	}
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

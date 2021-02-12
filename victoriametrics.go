package victoriametrics

import (
	"io"
	"strings"
	"time"

	"github.com/VictoriaMetrics/metrics"
	"github.com/unistack-org/micro/v3/meter"
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

func (r *victoriametricsMeter) buildName(name string, opts ...meter.Option) string {
	var b strings.Builder

	options := r.opts
	for _, o := range opts {
		o(&options)
	}

	if len(options.MetricPrefix) > 0 {
		_, _ = b.WriteString(options.MetricPrefix)
	}
	labelPrefix := false
	if len(options.LabelPrefix) > 0 {
		labelPrefix = true
	}
	_, _ = b.WriteString(name)
	if options.Labels.Len() > 0 {
		_, _ = b.WriteRune('{')
		iter := options.Labels.Iter()
		var k, v string
		for idx := 0; iter.Next(&k, &v); idx++ {
			if idx > 0 {
				_, _ = b.WriteRune(',')
			}
			if labelPrefix {
				_, _ = b.WriteString(options.LabelPrefix)
			}
			_, _ = b.WriteString(k)
			_, _ = b.WriteString(`="`)
			_, _ = b.WriteString(v)
			_, _ = b.WriteString(`"`)
		}
		_, _ = b.WriteRune('}')
	}

	return b.String()
}

func (r *victoriametricsMeter) Counter(name string, opts ...meter.Option) meter.Counter {
	return r.set.GetOrCreateCounter(r.buildName(name, opts...))
}

func (r *victoriametricsMeter) FloatCounter(name string, opts ...meter.Option) meter.FloatCounter {
	return r.set.GetOrCreateFloatCounter(r.buildName(name, opts...))
}

func (r *victoriametricsMeter) Gauge(name string, f func() float64, opts ...meter.Option) meter.Gauge {
	return r.set.GetOrCreateGauge(r.buildName(name, opts...), f)
}

func (r *victoriametricsMeter) Histogram(name string, opts ...meter.Option) meter.Histogram {
	return r.set.GetOrCreateHistogram(r.buildName(name, opts...))
}

func (r *victoriametricsMeter) Summary(name string, opts ...meter.Option) meter.Summary {
	return r.set.GetOrCreateSummary(r.buildName(name, opts...))
}

func (r *victoriametricsMeter) SummaryExt(name string, window time.Duration, quantiles []float64, opts ...meter.Option) meter.Summary {
	return r.set.GetOrCreateSummaryExt(r.buildName(name, opts...), window, quantiles)
}

func (r *victoriametricsMeter) Set(opts ...meter.Option) meter.Meter {
	m := &victoriametricsMeter{opts: meter.NewOptions(opts...), set: metrics.NewSet()}
	return m
}

func (r *victoriametricsMeter) Init(opts ...meter.Option) error {
	for _, o := range opts {
		o(&r.opts)
	}

	return nil
}

func (r *victoriametricsMeter) Write(w io.Writer, withProcessMetrics bool) error {
	r.set.WritePrometheus(w)
	if withProcessMetrics {
		metrics.WriteProcessMetrics(w)
	}
	return nil
}

func (r *victoriametricsMeter) Options() meter.Options {
	return r.opts
}

func (r *victoriametricsMeter) String() string {
	return "victoriametrics"
}

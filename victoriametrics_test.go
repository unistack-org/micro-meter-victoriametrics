package victoriametrics

import (
	"context"
	"testing"

	"github.com/unistack-org/micro/v3/client"
	"github.com/unistack-org/micro/v3/codec"
	"github.com/unistack-org/micro/v3/meter"
	"github.com/unistack-org/micro/v3/meter/wrapper"
)

func TestBuildName(t *testing.T) {
	m := NewMeter()
	im := m.(*victoriametricsMeter)
	check := `micro_foo{micro_aaa="b",micro_bar="baz"}`
	name := im.buildName("foo", meter.Label("bar", "baz"), meter.Label("aaa", "b"))
	if name != check {
		t.Fatalf("metric name error: %s != %s", name, check)
	}

	cnt := m.Counter("counter", meter.Label("key", "val"))
	cnt.Inc()
	//m.Write(os.Stdout, meter.WriteProcessMetrics(true), meter.WriteFDMetrics(true))
}

func TestWrapper(t *testing.T) {
	m := NewMeter()

	w := wrapper.NewClientWrapper(
		wrapper.ServiceName("svc1"),
		wrapper.ServiceVersion("0.0.1"),
		wrapper.ServiceID("12345"),
		wrapper.Meter(m),
	)

	ctx := context.Background()

	c := client.NewClient(client.Wrap(w))
	if err := c.Init(); err != nil {
		t.Fatal(err)
	}
	rsp := &codec.Frame{}
	req := &codec.Frame{}
	err := c.Call(ctx, c.NewRequest("svc2", "Service.Method", req), rsp)
	_, _ = rsp, err
	//m.Write(os.Stdout, meter.WriteProcessMetrics(true), meter.WriteFDMetrics(true))
}

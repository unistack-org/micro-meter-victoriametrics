package victoriametrics

import (
	"bytes"
	"context"
	"testing"

	"go.unistack.org/micro/v4/client"
	"go.unistack.org/micro/v4/codec"
	"go.unistack.org/micro/v4/meter"
)

func TestBuildName(t *testing.T) {
	m := NewMeter()
	im := m.(*victoriametricsMeter)
	check := `micro_foo{aaa="b",bar="baz",ccc="d"}`
	name := im.buildName("micro_foo", "bar", "baz", "aaa", "b", "ccc", "d")
	if name != check {
		t.Fatalf("metric name error: %s != %s", name, check)
	}

	cnt := m.Counter("counter", "key", "val")
	cnt.Inc()
}

func TestWrapper(t *testing.T) {
	m := NewMeter()
	if err := m.Init(); err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()

	c := client.NewClient(client.Meter(m))
	if err := c.Init(); err != nil {
		t.Fatal(err)
	}
	rsp := &codec.Frame{}
	req := &codec.Frame{}
	err := c.Call(ctx, c.NewRequest("svc2", "Service.Method", req), rsp)
	_, _ = rsp, err
	buf := bytes.NewBuffer(nil)
	_ = m.Write(buf, meter.WriteProcessMetrics(false), meter.WriteFDMetrics(false))
	if !bytes.Contains(buf.Bytes(), []byte(`micro_client_request_total{code="500",endpoint="Service.Method",status="failure"} 1`)) {
		t.Fatalf("invalid metrics output: %s", buf.Bytes())
	}
}

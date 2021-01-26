package victoriametrics

import (
	"os"
	"testing"

	"github.com/unistack-org/micro/v3/meter"
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
	m.Write(os.Stdout, true)
}

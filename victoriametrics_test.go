package victoriametrics

import (
	"testing"
)

func TestBuildName(t *testing.T) {
	m := NewMeter()
	im := m.(*victoriametricsMeter)
	check := `micro_foo{micro_aaa="b",micro_bar="baz",micro_ccc="d"}`
	name := im.buildName("micro_foo", "micro_bar", "baz", "micro_aaa", "b", "micro_ccc", "d")
	if name != check {
		t.Fatalf("metric name error: %s != %s", name, check)
	}

	cnt := m.Counter("counter", "key", "val")
	cnt.Inc()
}

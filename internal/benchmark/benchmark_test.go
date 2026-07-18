package benchmark

import "testing"

func TestPercentile(t *testing.T) {
	values := []float64{4, 1, 2, 3}
	if got := Percentile(values, .5); got != 3 {
		t.Fatalf("p50=%v", got)
	}
	if got := Percentile(values, .95); got != 4 {
		t.Fatalf("p95=%v", got)
	}
}

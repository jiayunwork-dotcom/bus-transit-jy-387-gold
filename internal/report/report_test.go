package report

import (
	"strings"
	"testing"

	"bus-transit/internal/metrics"
)

func TestFormatText(t *testing.T) {
	rs := []metrics.RouteMetric{{Route: "R2", OnTimeRate: 0, AvgLoad: 1.2, Underperforming: true}}
	out, err := Format(rs, "text")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, "R2:") || !strings.Contains(out, "UNDERPERFORMING") {
		t.Fatalf("bad output: %q", out)
	}
}

func TestFormatJSON(t *testing.T) {
	rs := []metrics.RouteMetric{{Route: "R1", OnTimeRate: 1.0}}
	out, err := Format(rs, "json")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, "\"Route\": \"R1\"") {
		t.Fatalf("bad output: %q", out)
	}
}

func TestFormatUnsupported(t *testing.T) {
	if _, err := Format(nil, "yaml"); err == nil {
		t.Fatal("expected error for unsupported format")
	}
}

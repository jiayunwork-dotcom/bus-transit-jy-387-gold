package metrics

import (
	"math"
	"testing"

	"bus-transit/internal/parse"
)

func approx(a, b float64) bool { return math.Abs(a-b) < 1e-6 }

func TestOnTimeRate(t *testing.T) {
	stops := []parse.StopTime{
		{SchedArr: 0, ActArr: 1},
		{SchedArr: 0, ActArr: 5},
		{SchedArr: 0, ActArr: 0},
	}
	if got := OnTimeRate(stops, 3); !approx(got, 2.0/3.0) {
		t.Fatalf("expected 2/3, got %v", got)
	}
}

func TestOnTimeRateEmpty(t *testing.T) {
	if got := OnTimeRate(nil, 3); got != 1.0 {
		t.Fatalf("expected 1.0 for empty, got %v", got)
	}
}

func TestAvgLoad(t *testing.T) {
	stops := []parse.StopTime{
		{Passengers: 25, Capacity: 50},
		{Passengers: 30, Capacity: 50},
	}
	if got := AvgLoad(stops); !approx(got, 0.55) {
		t.Fatalf("expected 0.55, got %v", got)
	}
}

func TestEvaluateRoutesUnderperforming(t *testing.T) {
	stops := []parse.StopTime{
		{Route: "R1", SchedArr: 0, ActArr: 0, Passengers: 10, Capacity: 50},
		{Route: "R2", SchedArr: 0, ActArr: 10, Passengers: 60, Capacity: 50},
	}
	cfg := Config{OnTimeThresholdMin: 3, OnTimeMin: 0.8, LoadMax: 1.0}
	rs := EvaluateRoutes(stops, cfg)
	if len(rs) != 2 {
		t.Fatalf("expected 2 routes, got %d", len(rs))
	}
	for _, r := range rs {
		if r.Route == "R2" && !r.Underperforming {
			t.Fatalf("expected R2 underperforming, got %+v", r)
		}
		if r.Route == "R1" && r.Underperforming {
			t.Fatalf("expected R1 ok, got %+v", r)
		}
	}
}

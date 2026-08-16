// Package metrics computes per-route transit efficiency.
package metrics

import "bus-transit/internal/parse"

// Config holds the evaluation thresholds.
type Config struct {
	OnTimeThresholdMin int     // allowed arrival deviation in minutes
	OnTimeMin          float64 // minimum acceptable on-time rate
	LoadMax            float64 // maximum acceptable average load factor
}

// RouteMetric summarizes one route.
type RouteMetric struct {
	Route            string
	OnTimeRate       float64
	AvgLoad          float64
	Underperforming  bool
}

// OnTimeRate returns the fraction of stops arriving within the threshold.
func OnTimeRate(stops []parse.StopTime, thresholdMin int) float64 {
	if len(stops) == 0 {
		return 1.0
	}
	onTime := 0
	for _, s := range stops {
		d := s.ActArr - s.SchedArr
		if d < 0 {
			d = -d
		}
		if d <= thresholdMin {
			onTime++
		}
	}
	return float64(onTime) / float64(len(stops))
}

// AvgLoad returns the mean passenger/capacity ratio across stops.
func AvgLoad(stops []parse.StopTime) float64 {
	if len(stops) == 0 {
		return 0
	}
	var sum float64
	for _, s := range stops {
		if s.Capacity > 0 {
			sum += float64(s.Passengers) / float64(s.Capacity)
		}
	}
	return sum / float64(len(stops))
}

// EvaluateRoutes groups stop records by route and scores each.
func EvaluateRoutes(stops []parse.StopTime, cfg Config) []RouteMetric {
	byRoute := map[string][]parse.StopTime{}
	var order []string
	for _, s := range stops {
		if _, ok := byRoute[s.Route]; !ok {
			order = append(order, s.Route)
		}
		byRoute[s.Route] = append(byRoute[s.Route], s)
	}

	out := make([]RouteMetric, 0, len(order))
	for _, r := range order {
		ss := byRoute[r]
		m := RouteMetric{
			Route:      r,
			OnTimeRate: OnTimeRate(ss, cfg.OnTimeThresholdMin),
			AvgLoad:    AvgLoad(ss),
		}
		if m.OnTimeRate < cfg.OnTimeMin || m.AvgLoad > cfg.LoadMax {
			m.Underperforming = true
		}
		out = append(out, m)
	}
	return out
}

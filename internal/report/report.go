// Package report renders bus-transit route metrics.
package report

import (
	"encoding/json"
	"fmt"
	"strings"

	"bus-transit/internal/metrics"
)

// Format renders the metrics in the requested format (text or json).
func Format(routes []metrics.RouteMetric, format string) (string, error) {
	switch format {
	case "json", "":
		b, err := json.MarshalIndent(routes, "", "  ")
		if err != nil {
			return "", fmt.Errorf("marshal report: %w", err)
		}
		return string(b) + "\n", nil
	case "text":
		return textFormat(routes), nil
	default:
		return "", fmt.Errorf("unsupported format %q (want text or json)", format)
	}
}

func textFormat(routes []metrics.RouteMetric) string {
	var b strings.Builder
	for _, r := range routes {
		status := "ok"
		if r.Underperforming {
			status = "UNDERPERFORMING"
		}
		fmt.Fprintf(&b, "%s: on_time_rate=%.2f avg_load=%.2f [%s]\n",
			r.Route, r.OnTimeRate, r.AvgLoad, status)
	}
	return b.String()
}

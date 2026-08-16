package main

import (
	"flag"
	"fmt"
	"os"

	"bus-transit/internal/metrics"
	"bus-transit/internal/parse"
	"bus-transit/internal/report"
)

func main() {
	stopsPath := flag.String("stops", "", "path to stop-times CSV")
	format := flag.String("format", "text", "output format: text or json")
	flag.Parse()

	if *stopsPath == "" {
		fmt.Fprintln(os.Stderr, "usage: bus-transit -stops <csv> [-format text|json]")
		os.Exit(2)
	}

	stops, err := parse.ReadStopTimes(*stopsPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}

	cfg := metrics.Config{OnTimeThresholdMin: 3, OnTimeMin: 0.8, LoadMax: 1.0}
	routes := metrics.EvaluateRoutes(stops, cfg)

	out, err := report.Format(routes, *format)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
	fmt.Print(out)

	if hasUnderperforming(routes) {
		os.Exit(1)
	}
}

func hasUnderperforming(routes []metrics.RouteMetric) bool {
	for _, r := range routes {
		if r.Underperforming {
			return true
		}
	}
	return false
}

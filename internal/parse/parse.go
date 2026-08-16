// Package parse reads bus stop-time records.
package parse

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// StopTime is one observed stop record for a trip.
type StopTime struct {
	Route      string
	TripID     string
	Stop       string
	SchedArr   int // scheduled arrival, minutes since midnight
	ActArr     int // actual arrival, minutes since midnight
	Passengers int
	Capacity   int
}

// ReadStopTimes reads a CSV with header:
// route,trip_id,stop,sched_arr,act_arr,passengers,capacity.
func ReadStopTimes(path string) ([]StopTime, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open stops: %w", err)
	}
	defer f.Close()

	r := csv.NewReader(f)
	records, err := r.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("read stops csv: %w", err)
	}
	if len(records) < 2 {
		return nil, fmt.Errorf("stops file has no data rows")
	}
	header := records[0]
	ri, ti, si, sai, aai, pi, ci := -1, -1, -1, -1, -1, -1, -1
	for i, h := range header {
		switch strings.ToLower(strings.TrimSpace(h)) {
		case "route":
			ri = i
		case "trip_id":
			ti = i
		case "stop":
			si = i
		case "sched_arr":
			sai = i
		case "act_arr":
			aai = i
		case "passengers":
			pi = i
		case "capacity":
			ci = i
		}
	}
	if ri < 0 || ti < 0 || si < 0 || sai < 0 || aai < 0 || pi < 0 || ci < 0 {
		return nil, fmt.Errorf("stops header must contain route,trip_id,stop,sched_arr,act_arr,passengers,capacity")
	}

	var out []StopTime
	for n, row := range records[1:] {
		if len(row) <= ci {
			return nil, fmt.Errorf("row %d: missing stop column", n+2)
		}
		sched, err := parseHHMM(row[sai])
		if err != nil {
			return nil, fmt.Errorf("row %d: invalid sched_arr: %w", n+2, err)
		}
		act, err := parseHHMM(row[aai])
		if err != nil {
			return nil, fmt.Errorf("row %d: invalid act_arr: %w", n+2, err)
		}
		pass, err := strconv.Atoi(strings.TrimSpace(row[pi]))
		if err != nil {
			return nil, fmt.Errorf("row %d: invalid passengers %q: %w", n+2, row[pi], err)
		}
		capv, err := strconv.Atoi(strings.TrimSpace(row[ci]))
		if err != nil {
			return nil, fmt.Errorf("row %d: invalid capacity %q: %w", n+2, row[ci], err)
		}
		if capv <= 0 {
			return nil, fmt.Errorf("row %d: capacity must be positive", n+2)
		}
		out = append(out, StopTime{
			Route:      strings.TrimSpace(row[ri]),
			TripID:     strings.TrimSpace(row[ti]),
			Stop:       strings.TrimSpace(row[si]),
			SchedArr:   sched,
			ActArr:     act,
			Passengers: pass,
			Capacity:   capv,
		})
	}
	return out, nil
}

func parseHHMM(s string) (int, error) {
	parts := strings.Split(strings.TrimSpace(s), ":")
	if len(parts) != 2 {
		return 0, fmt.Errorf("invalid time %q", s)
	}
	h, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, err
	}
	m, err := strconv.Atoi(parts[1])
	if err != nil {
		return 0, err
	}
	if h < 0 || h > 23 || m < 0 || m > 59 {
		return 0, fmt.Errorf("invalid time %q", s)
	}
	return h*60 + m, nil
}

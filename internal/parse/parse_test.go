package parse

import (
	"os"
	"path/filepath"
	"testing"
)

const stopsCSV = `route,trip_id,stop,sched_arr,act_arr,passengers,capacity
R1,T1,S1,08:00,08:00,20,50
R1,T1,S2,08:10,08:11,25,50
`

func TestReadStopTimesValid(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "stops.csv")
	if err := os.WriteFile(p, []byte(stopsCSV), 0o644); err != nil {
		t.Fatal(err)
	}
	stops, err := ReadStopTimes(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(stops) != 2 {
		t.Fatalf("expected 2 stops, got %d", len(stops))
	}
	if stops[0].Route != "R1" || stops[0].SchedArr != 8*60 || stops[0].ActArr != 8*60 {
		t.Fatalf("unexpected stop: %+v", stops[0])
	}
}

func TestReadStopTimesMissingFile(t *testing.T) {
	if _, err := ReadStopTimes(filepath.Join(t.TempDir(), "nope.csv")); err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestReadStopTimesBadTime(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "stops.csv")
	if err := os.WriteFile(p, []byte("route,trip_id,stop,sched_arr,act_arr,passengers,capacity\nR1,T1,S1,08:00,9o'clock,10,50\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := ReadStopTimes(p); err == nil {
		t.Fatal("expected error for invalid time")
	}
}

func TestReadStopTimesBadCapacity(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "stops.csv")
	if err := os.WriteFile(p, []byte("route,trip_id,stop,sched_arr,act_arr,passengers,capacity\nR1,T1,S1,08:00,08:00,10,0\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := ReadStopTimes(p); err == nil {
		t.Fatal("expected error for non-positive capacity")
	}
}

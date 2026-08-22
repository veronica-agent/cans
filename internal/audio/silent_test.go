package audio

import (
	"math"
	"testing"
)

func TestSilentEmpty(t *testing.T) {
	if !Silent(nil, 24000) {
		t.Fatal("nil should be silent")
	}
}

func TestSilentShortDuration(t *testing.T) {
	s := make([]float32, 1200) // 50ms at 24kHz
	for i := range s {
		s[i] = 0.3
	}
	if !Silent(s, 24000) {
		t.Fatal("short duration should be silent")
	}
}

func TestSilentLowPeak(t *testing.T) {
	sr := 24000
	n := sr / 10
	s := make([]float32, n)
	for i := range s {
		s[i] = 0.01 * float32(math.Sin(2*math.Pi*440*float64(i)/float64(sr)))
	}
	if !Silent(s, sr) {
		t.Fatal("low peak should be silent")
	}
}

// TestCleanRescuesLowPeakSpeech proves that a long, quiet-but-real clone
// (200ms, peak 0.015) survives Clean: after Clean it is NOT Silent and its
// length stays on the order of 200ms, not the 30ms trim pad. SayTo runs
// Clean then Silent, so such a clone is kept instead of being rejected as a
// silent mouth.
func TestCleanRescuesLowPeakSpeech(t *testing.T) {
	sr := 24000
	n := sr * 200 / 1000 // 200ms
	src := make([]float32, n)
	for i := range src {
		src[i] = 0.015 * float32(math.Sin(2*math.Pi*440*float64(i)/float64(sr)))
	}
	if !Silent(src, sr) {
		t.Fatal("200ms peak 0.015 should be silent before Clean")
	}
	out := Clean(append([]float32(nil), src...), sr)
	if Silent(out, sr) {
		t.Fatal("200ms peak 0.015 should NOT be silent after Clean")
	}
	if len(out)*1000/sr < 100 {
		t.Fatalf("Clean collapsed 200ms to %dms; want ~200ms", len(out)*1000/sr)
	}
	var peak float32
	for _, x := range out {
		if x < 0 {
			x = -x
		}
		if x > peak {
			peak = x
		}
	}
	if peak < 0.4 {
		t.Fatalf("peak %f after Clean, want ~0.5", peak)
	}
}

// TestCleanSilentStaysSilent proves near-zero or very short audio is still
// Silent after Clean, so a glitch is rejected as a silent mouth.
func TestCleanSilentStaysSilent(t *testing.T) {
	sr := 24000
	// Near-zero: 200ms of 1e-5 noise — well below the trim threshold.
	n := sr * 200 / 1000
	flat := make([]float32, n)
	for i := range flat {
		flat[i] = 1e-5 * float32(math.Sin(2*math.Pi*440*float64(i)/float64(sr)))
	}
	out := Clean(flat, sr)
	if !Silent(out, sr) {
		t.Fatalf("near-zero should stay silent after Clean: %d samples", len(out))
	}
	// Very short: 30ms of loud tone — below the 80ms minimum.
	short := make([]float32, sr*30/1000)
	for i := range short {
		short[i] = 0.3 * float32(math.Sin(2*math.Pi*440*float64(i)/float64(sr)))
	}
	out = Clean(short, sr)
	if !Silent(out, sr) {
		t.Fatalf("30ms should stay silent after Clean: %d samples", len(out))
	}
}

func TestSilentOK(t *testing.T) {
	sr := 24000
	n := sr / 10
	s := make([]float32, n)
	for i := range s {
		s[i] = 0.3 * float32(math.Sin(2*math.Pi*440*float64(i)/float64(sr)))
	}
	if Silent(s, sr) {
		t.Fatal("normal audio should not be silent")
	}
}

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

// TestSilentLowPeakRescuedByNormalize proves that a long, quiet-but-real clone
// (200ms, peak 0.015) is Silent before Normalize but NOT Silent after Normalize
// to 0.5. SayTo runs Clean → Normalize → Silent, so such a clone survives
// instead of being rejected as a silent mouth.
func TestSilentLowPeakRescuedByNormalize(t *testing.T) {
	sr := 24000
	n := sr * 200 / 1000 // 200ms
	s := make([]float32, n)
	for i := range s {
		s[i] = 0.015 * float32(math.Sin(2*math.Pi*440*float64(i)/float64(sr)))
	}
	if !Silent(s, sr) {
		t.Fatal("200ms peak 0.015 should be silent before normalize")
	}
	out := Normalize(s, 0.5)
	if Silent(out, sr) {
		t.Fatal("200ms peak 0.015 should NOT be silent after normalize to 0.5")
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

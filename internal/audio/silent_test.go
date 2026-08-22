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

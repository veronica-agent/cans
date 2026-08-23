package audio

import (
	"math"
	"testing"
)

func TestNormalizeNil(t *testing.T) {
	if got := Normalize(nil, 0.5); got != nil {
		t.Fatalf("%v", got)
	}
}

func TestNormalizeEmpty(t *testing.T) {
	s := []float32{}
	if got := Normalize(s, 0.5); len(got) != 0 {
		t.Fatalf("%v", got)
	}
}

func TestNormalizeSineToHalf(t *testing.T) {
	sr := 24000
	n := sr / 10
	s := make([]float32, n)
	for i := range s {
		s[i] = 0.1 * float32(math.Sin(2*math.Pi*440*float64(i)/float64(sr)))
	}
	out := Normalize(s, 0.5)
	var peak float32
	for _, x := range out {
		if x < 0 {
			x = -x
		}
		if x > peak {
			peak = x
		}
	}
	if math.Abs(float64(peak)-0.5) > 0.001 {
		t.Fatalf("peak %f, want 0.5", peak)
	}
}

func TestNormalizeAlreadyLoud(t *testing.T) {
	sr := 24000
	n := sr / 10
	s := make([]float32, n)
	for i := range s {
		s[i] = 0.7 * float32(math.Sin(2*math.Pi*440*float64(i)/float64(sr)))
	}
	out := Normalize(s, 0.5)
	var peak float32
	for _, x := range out {
		if x < 0 {
			x = -x
		}
		if x > peak {
			peak = x
		}
	}
	if math.Abs(float64(peak)-0.7) > 0.001 {
		t.Fatalf("peak %f, want 0.7 unchanged", peak)
	}
}

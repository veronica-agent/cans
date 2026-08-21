package audio

import (
	"math"
	"testing"
)

func TestCleanKillsDC(t *testing.T) {
	sr := 24000
	n := sr / 2
	s := make([]float32, n)
	for i := range s {
		s[i] = 0.25 + 0.4*float32(math.Sin(2*math.Pi*200*float64(i)/float64(sr)))
	}
	out := Clean(s, sr)
	if len(out) < n/4 {
		t.Fatalf("trimmed too much: %d", len(out))
	}
	var sum float64
	for _, x := range out {
		sum += float64(x)
	}
	mean := sum / float64(len(out))
	if mean > 0.05 || mean < -0.05 {
		t.Fatalf("dc remains %f", mean)
	}
}

func TestCleanTrimsQuietEdges(t *testing.T) {
	sr := 24000
	s := make([]float32, sr)
	for i := sr / 5; i < 2*sr/5; i++ {
		s[i] = 0.4
	}
	out := Clean(s, sr)
	if len(out) >= len(s) {
		t.Fatalf("did not trim: %d", len(out))
	}
	if len(out) < sr/10 {
		t.Fatalf("trimmed too much: %d", len(out))
	}
}

func TestCleanShelvesHiss(t *testing.T) {
	sr := 24000
	n := sr / 2
	s := make([]float32, n)
	for i := range s {
		s[i] = 0.25*float32(math.Sin(2*math.Pi*1000*float64(i)/float64(sr))) +
			0.25*float32(math.Sin(2*math.Pi*10000*float64(i)/float64(sr)))
	}
	rawHi := toneRMS(s, sr, 10000)
	out := Clean(append([]float32(nil), s...), sr)
	hi := toneRMS(out, sr, 10000)
	lo := toneRMS(out, sr, 1000)
	if hi >= rawHi {
		t.Fatalf("hiss not reduced: before %f after %f", rawHi, hi)
	}
	if lo <= hi {
		t.Fatalf("shelf took the voice: lo %f hi %f", lo, hi)
	}
}

func toneRMS(s []float32, sr int, hz float64) float64 {
	var sum float64
	for i, x := range s {
		sum += float64(x) * math.Sin(2*math.Pi*hz*float64(i)/float64(sr))
	}
	return math.Abs(sum / float64(len(s)))
}

func TestCleanEmpty(t *testing.T) {
	if got := Clean(nil, 24000); got != nil {
		t.Fatalf("%v", got)
	}
}

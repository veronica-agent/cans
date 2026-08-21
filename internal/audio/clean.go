package audio

import "math"

const (
	highpassHz  = 80
	shelfHz     = 7000
	shelfGainDB = -8
	trimAmp     = 0.02 // ~-34 dB
	padMs       = 30
	fadeMs      = 8
)

// Clean high-passes rumble, shelves vocoder hiss, trims quiet edges, fades the cut.
func Clean(samples []float32, sampleRate int) []float32 {
	if len(samples) == 0 || sampleRate <= 0 {
		return samples
	}
	highpass(samples, sampleRate, highpassHz)
	highshelf(samples, sampleRate, shelfHz, shelfGainDB)
	samples = trimSilence(samples, sampleRate, trimAmp, padMs)
	fade(samples, sampleRate, fadeMs)
	return samples
}

func highpass(samples []float32, sampleRate int, cutoff float64) {
	w0 := 2 * math.Pi * cutoff / float64(sampleRate)
	cos := math.Cos(w0)
	sin := math.Sin(w0)
	alpha := sin / math.Sqrt2
	b0 := (1 + cos) / 2
	b1 := -(1 + cos)
	b2 := (1 + cos) / 2
	a0 := 1 + alpha
	a1 := -2 * cos
	a2 := 1 - alpha
	b0 /= a0
	b1 /= a0
	b2 /= a0
	a1 /= a0
	a2 /= a0
	biquad(samples, b0, b1, b2, a1, a2)
}

func highshelf(samples []float32, sampleRate int, cutoff, gainDB float64) {
	a := math.Pow(10, gainDB/40)
	w0 := 2 * math.Pi * cutoff / float64(sampleRate)
	cos := math.Cos(w0)
	sin := math.Sin(w0)
	alpha := sin / 2 * math.Sqrt2
	sa := 2 * math.Sqrt(a) * alpha
	b0 := a * ((a + 1) + (a-1)*cos + sa)
	b1 := -2 * a * ((a - 1) + (a+1)*cos)
	b2 := a * ((a + 1) + (a-1)*cos - sa)
	a0 := (a + 1) - (a-1)*cos + sa
	a1 := 2 * ((a - 1) - (a+1)*cos)
	a2 := (a + 1) - (a-1)*cos - sa
	b0 /= a0
	b1 /= a0
	b2 /= a0
	a1 /= a0
	a2 /= a0
	biquad(samples, b0, b1, b2, a1, a2)
}

func biquad(samples []float32, b0, b1, b2, a1, a2 float64) {
	var x1, x2, y1, y2 float64
	for i, x := range samples {
		xf := float64(x)
		y := b0*xf + b1*x1 + b2*x2 - a1*y1 - a2*y2
		x2, x1 = x1, xf
		y2, y1 = y1, y
		samples[i] = float32(y)
	}
}

func trimSilence(samples []float32, sampleRate int, thresh float64, padMs int) []float32 {
	n := len(samples)
	if n == 0 {
		return samples
	}
	start := 0
	for start < n && abs32(samples[start]) < thresh {
		start++
	}
	end := n - 1
	for end > start && abs32(samples[end]) < thresh {
		end--
	}
	pad := sampleRate * padMs / 1000
	if pad < 0 {
		pad = 0
	}
	start -= pad
	if start < 0 {
		start = 0
	}
	end += pad
	if end >= n {
		end = n - 1
	}
	if start == 0 && end == n-1 {
		return samples
	}
	out := make([]float32, end-start+1)
	copy(out, samples[start:end+1])
	return out
}

func fade(samples []float32, sampleRate int, ms int) {
	n := sampleRate * ms / 1000
	if n <= 0 || len(samples) < 2*n {
		return
	}
	for i := 0; i < n; i++ {
		g := float32(i) / float32(n)
		samples[i] *= g
		samples[len(samples)-1-i] *= g
	}
}

func abs32(v float32) float64 {
	if v < 0 {
		return float64(-v)
	}
	return float64(v)
}

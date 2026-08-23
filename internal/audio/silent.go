package audio

const (
	silentMinMs = 80
	silentPeak  = 0.02
)

// Silent reports whether samples are too short or too quiet to be real speech:
// fewer than 80ms of audio or a peak amplitude below 0.02.
func Silent(samples []float32, sampleRate int) bool {
	if len(samples) == 0 || sampleRate <= 0 {
		return true
	}
	if len(samples)*1000/sampleRate < silentMinMs {
		return true
	}
	var peak float32
	for _, x := range samples {
		if x < 0 {
			x = -x
		}
		if x > peak {
			peak = x
		}
	}
	return peak < silentPeak
}

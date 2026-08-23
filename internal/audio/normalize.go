package audio

// Normalize scales samples so the loudest sample reaches peak, unless the
// existing peak is already at or above peak. Empty or all-zero input is
// returned unchanged.
func Normalize(samples []float32, peak float32) []float32 {
	if len(samples) == 0 {
		return samples
	}
	var max float32
	for _, x := range samples {
		if x < 0 {
			x = -x
		}
		if x > max {
			max = x
		}
	}
	if max <= 0 || max >= peak {
		return samples
	}
	gain := peak / max
	for i := range samples {
		samples[i] *= gain
	}
	return samples
}

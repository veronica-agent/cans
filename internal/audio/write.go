package audio

import (
	"encoding/binary"
	"fmt"
	"os"
)

// WritePCM16 writes mono float32 samples as a 16-bit PCM WAV.
func WritePCM16(path string, sampleRate int, samples []float32) error {
	if sampleRate <= 0 {
		return fmt.Errorf("invalid sample rate %d", sampleRate)
	}
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	n := len(samples)
	dataBytes := n * 2
	var hdr [44]byte
	copy(hdr[0:], "RIFF")
	binary.LittleEndian.PutUint32(hdr[4:], uint32(36+dataBytes))
	copy(hdr[8:], "WAVE")
	copy(hdr[12:], "fmt ")
	binary.LittleEndian.PutUint32(hdr[16:], 16)
	binary.LittleEndian.PutUint16(hdr[20:], 1)
	binary.LittleEndian.PutUint16(hdr[22:], 1)
	binary.LittleEndian.PutUint32(hdr[24:], uint32(sampleRate))
	binary.LittleEndian.PutUint32(hdr[28:], uint32(sampleRate*2))
	binary.LittleEndian.PutUint16(hdr[32:], 2)
	binary.LittleEndian.PutUint16(hdr[34:], 16)
	copy(hdr[36:], "data")
	binary.LittleEndian.PutUint32(hdr[40:], uint32(dataBytes))
	if _, err := f.Write(hdr[:]); err != nil {
		return err
	}
	for _, s := range samples {
		if s > 1 {
			s = 1
		} else if s < -1 {
			s = -1
		}
		var b [2]byte
		binary.LittleEndian.PutUint16(b[:], uint16(int16(s*32767)))
		if _, err := f.Write(b[:]); err != nil {
			return err
		}
	}
	return nil
}

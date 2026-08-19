package audio

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"
)

// HeaderOK reports whether path is a regular RIFF/WAVE file.
func HeaderOK(path string) error {
	info, err := os.Lstat(path)
	if err != nil {
		return err
	}
	if !info.Mode().IsRegular() {
		return fmt.Errorf("not a regular file: %s", path)
	}
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()
	var h [12]byte
	if _, err := io.ReadFull(f, h[:]); err != nil {
		return fmt.Errorf("not a wav: %s", path)
	}
	if string(h[0:4]) != "RIFF" || string(h[8:12]) != "WAVE" {
		return fmt.Errorf("not a wav: %s", path)
	}
	return nil
}

// Minimal is a 24 kHz mono 16-bit silent WAV (one sample).
func Minimal() []byte {
	const dataBytes = 2
	buf := make([]byte, 44+dataBytes)
	copy(buf[0:4], "RIFF")
	binary.LittleEndian.PutUint32(buf[4:8], uint32(36+dataBytes))
	copy(buf[8:12], "WAVE")
	copy(buf[12:16], "fmt ")
	binary.LittleEndian.PutUint32(buf[16:20], 16)
	binary.LittleEndian.PutUint16(buf[20:22], 1)
	binary.LittleEndian.PutUint16(buf[22:24], 1)
	binary.LittleEndian.PutUint32(buf[24:28], 24000)
	binary.LittleEndian.PutUint32(buf[28:32], 48000)
	binary.LittleEndian.PutUint16(buf[32:34], 2)
	binary.LittleEndian.PutUint16(buf[34:36], 16)
	copy(buf[36:40], "data")
	binary.LittleEndian.PutUint32(buf[40:44], dataBytes)
	return buf
}

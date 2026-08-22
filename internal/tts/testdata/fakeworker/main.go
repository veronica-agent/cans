// Fake qwen3-tts-worker for unit tests. Emits a 0.3-amplitude 440 Hz sine
// (100 ms at 24 kHz) so happy-path tests produce audible, non-silent audio.
package main

import (
	"bufio"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"os"
)

func main() {
	fmt.Println(`{"type":"ready","protocol":"qwen3-tts-worker/v1","sample_rate":24000,"pcm_format":"f32le","streaming":false}`)
	sc := bufio.NewScanner(os.Stdin)
	for sc.Scan() {
		line := sc.Text()
		var req struct {
			Type string `json:"type"`
			ID   string `json:"id"`
			Text string `json:"text"`
		}
		_ = json.Unmarshal([]byte(line), &req)
		if req.Type == "shutdown" {
			return
		}
		if req.Type != "synthesize" {
			continue
		}
		if req.Text == "" {
			fmt.Printf("{\"type\":\"error\",\"id\":%q,\"message\":\"missing text\"}\n", req.ID)
			continue
		}
		if req.Text == "fail" {
			fmt.Printf("{\"type\":\"error\",\"id\":%q,\"message\":\"fail\"}\n", req.ID)
			continue
		}
		if req.Text == "block" {
			block()
			return
		}
		if req.Text == "silent" {
			fmt.Printf("{\"type\":\"pcm_meta\",\"id\":%q,\"sample_rate\":24000,\"format\":\"f32le\",\"n_samples\":1}\n", req.ID)
			var b [4]byte
			binary.LittleEndian.PutUint32(b[:], math.Float32bits(0))
			_, _ = os.Stdout.Write(b[:])
			fmt.Printf("\n{\"type\":\"final\",\"id\":%q}\n", req.ID)
			continue
		}
		emitSine(req.ID)
	}
}

func emitSine(id string) {
	const sr = 24000
	n := 2400
	fmt.Printf("{\"type\":\"pcm_meta\",\"id\":%q,\"sample_rate\":%d,\"format\":\"f32le\",\"n_samples\":%d}\n", id, sr, n)
	buf := make([]byte, n*4)
	for i := 0; i < n; i++ {
		s := float32(0.3 * math.Sin(2*math.Pi*440*float64(i)/float64(sr)))
		binary.LittleEndian.PutUint32(buf[i*4:], math.Float32bits(s))
	}
	_, _ = os.Stdout.Write(buf)
	fmt.Printf("\n{\"type\":\"final\",\"id\":%q}\n", id)
}

// block answers nothing: it reports that synthesis started by creating
// CANS_FAKE_BLOCK_FILE, then waits in a read until the parent signals it. It
// stands in for the mouth mid-utterance, which cannot be aborted.
func block() {
	if p := os.Getenv("CANS_FAKE_BLOCK_FILE"); p != "" {
		_ = os.WriteFile(p, []byte("x"), 0o644)
	}
	_, _ = io.Copy(io.Discard, os.Stdin)
}

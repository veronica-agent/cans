package tts

import (
	"context"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"time"
)

func (c *Client) readPCM(ctx context.Context, start time.Time) (*pcmResult, error) {
	out := &pcmResult{}
	for {
		if err := ctx.Err(); err != nil {
			return nil, err
		}
		line, err := readLine(ctx, c.stdout)
		if err != nil {
			return nil, err
		}
		if line == "" {
			continue
		}
		var head struct {
			Type       string `json:"type"`
			Message    string `json:"message"`
			SampleRate int    `json:"sample_rate"`
			NSamples   int    `json:"n_samples"`
		}
		if err := json.Unmarshal([]byte(line), &head); err != nil {
			return nil, fmt.Errorf("line json: %w (%q)", err, line)
		}
		switch head.Type {
		case "error":
			return nil, fmt.Errorf("worker: %s", head.Message)
		case "pcm_meta":
			if err := c.appendPCM(head.SampleRate, head.NSamples, out); err != nil {
				return nil, err
			}
		case "final":
			if len(out.samples) == 0 {
				return nil, fmt.Errorf("final without pcm")
			}
			out.wall = time.Since(start)
			return out, nil
		case "generating":
		default:
		}
	}
}

func (c *Client) appendPCM(rate, n int, out *pcmResult) error {
	if n <= 0 {
		return fmt.Errorf("pcm_meta n_samples=%d", n)
	}
	buf := make([]byte, n*4)
	if _, err := io.ReadFull(c.stdout, buf); err != nil {
		return fmt.Errorf("pcm read: %w", err)
	}
	if b, err := c.stdout.ReadByte(); err == nil && b != '\n' {
		_ = c.stdout.UnreadByte()
	}
	samples := make([]float32, n)
	for i := 0; i < n; i++ {
		samples[i] = math.Float32frombits(binary.LittleEndian.Uint32(buf[i*4:]))
	}
	out.sampleRate = rate
	out.samples = append(out.samples, samples...)
	return nil
}

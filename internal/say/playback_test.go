package say

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/veronica-agent/cans/internal/audio"
)

func TestRunPlaybackCancellation(t *testing.T) {
	for _, stream := range []bool{false, true} {
		for _, save := range []bool{false, true} {
			name := "once"
			if stream {
				name = "stream"
			}
			if save {
				name += "-saved"
			}
			t.Run(name, func(t *testing.T) {
				if stream {
					fakeWorkerEnv(t)
				} else {
					sayBinEnv(t)
				}
				t.Setenv("CANS_NOPLAY", "")
				dir := t.TempDir()
				started := filepath.Join(dir, "playing")
				t.Setenv("CANS_TEST_PLAYING", started)
				// Replace the player process with sleep so cancellation must
				// stop the actual child, without leaving a shell grandchild.
				body := "#!/bin/sh\nfor arg do wav=$arg; done\nprintf '%s' \"$wav\" > \"$CANS_TEST_PLAYING\"\nexec /bin/sleep 2\n"
				for _, player := range []string{"afplay", "ffplay"} {
					if err := os.WriteFile(filepath.Join(dir, player), []byte(body), 0o755); err != nil {
						t.Fatal(err)
					}
				}
				t.Setenv("PATH", dir+string(os.PathListSeparator)+os.Getenv("PATH"))
				o := DefaultOptions()
				o.Stream = stream
				if stream {
					o.Text = ""
				} else {
					o.Text = "Put the cans on."
				}
				if save {
					o.Play = true
					o.Out = filepath.Join(dir, "saved.wav")
					if stream {
						o.Out = filepath.Join(dir, "%03d.wav")
					}
				}
				ctx, cancel := context.WithCancel(context.Background())
				defer cancel()
				var out, stderr bytes.Buffer
				done := make(chan int, 1)
				go func() {
					done <- Run(ctx, o, strings.NewReader("Put the cans on.\nSecond line.\n"), &out, &stderr)
				}()
				waitFor(t, func() bool {
					data, err := os.ReadFile(started)
					return err == nil && len(data) > 0
				}, "player never started")
				began := time.Now()
				cancel()
				select {
				case code := <-done:
					if code != ExitInterrupted {
						t.Errorf("exit %d, want 130; stderr %q", code, stderr.String())
					}
				case <-time.After(5 * time.Second):
					t.Fatal("playback did not finish after cancellation")
				}
				if elapsed := time.Since(began); elapsed > time.Second {
					t.Errorf("playback kept running for %s after cancellation", elapsed)
				}
				want := "say: interrupted\n"
				if stream {
					want = "interrupted before the first line\n"
				}
				if stderr.String() != want {
					t.Errorf("stderr %q, want %q", stderr.String(), want)
				}
				data, err := os.ReadFile(started)
				if err != nil {
					t.Fatal(err)
				}
				wav := string(data)
				if save {
					if err := audio.HeaderOK(wav); err != nil {
						t.Errorf("saved WAV was lost: %v", err)
					}
				} else if _, err := os.Stat(wav); !os.IsNotExist(err) {
					t.Errorf("temporary WAV remains: %v", err)
				}
			})
		}
	}
}

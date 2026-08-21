//go:build darwin && arm64

package ship

import _ "embed"

// Patched worker_main (clone temp 0.2, max_tokens 360) linked against v0.1.0 dylib.
//
//go:embed overlay/qwen3-tts-worker-darwin-arm64
var workerOverlay []byte

# D001 — Go Charm booth, Python clone sidecar

**Decision:** The TUI is Go (bubbletea/lipgloss). The 0.6B Qwen Base clone runs in a Python sidecar (`mlx_audio`).

**Why:** Viral object is the Charm GIF; the clone API we already know is Python/MLX (`qwen.py` `_clone`). Pure Go ONNX of Qwen3-TTS is out of scope. Pure Rust is not the booth look and fights the mouth.

**Not:** Two user-facing tools. One `cans` binary that execs the sidecar.

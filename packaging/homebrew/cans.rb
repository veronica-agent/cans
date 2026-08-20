# Mirrors veronica-agent/homebrew-tap Formula/cans.rb
class Cans < Formula
  desc "Type a line. She speaks it"
  homepage "https://github.com/veronica-agent/cans"
  license "MIT"
  head "https://github.com/veronica-agent/cans.git", branch: "main"

  depends_on "go" => :build
  depends_on arch: :arm64
  depends_on :macos

  def install
    ldflags = "-s -w -X github.com/veronica-agent/cans/internal/ship.Version=#{version}"
    system "go", "build", *std_go_args(ldflags:), "./cmd/cans"
  end

  def caveats
    <<~EOS
      Apple Silicon. Unpack qwen3-tts-native into ~/.cans/native
      so bin/qwen3-tts-worker exists, then: cans doctor
    EOS
  end

  test do
    assert_match "cans", shell_output("#{bin}/cans version")
  end
end

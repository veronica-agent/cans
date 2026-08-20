package booth

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/veronica-agent/cans/internal/keep"
)

func TestEmptyEnterStaysListen(t *testing.T) {
	m := New("Put the cans on.", keep.Current{Wav: "/x", RefText: "hi"})
	next, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	got := next.(model)
	if got.status != "listen" || got.busy {
		t.Fatalf("status=%s busy=%v", got.status, got.busy)
	}
	if cmd != nil {
		t.Fatal("empty enter should not speak")
	}
}

func TestViewShowsQuote(t *testing.T) {
	m := New("Put the cans on.", keep.Current{Wav: "/x", RefText: "hi"})
	if !strings.Contains(m.View(), "Put the cans on.") {
		t.Fatalf("%q", m.View())
	}
	if !strings.Contains(m.View(), "listen") {
		t.Fatalf("%q", m.View())
	}
}

func TestEscQuits(t *testing.T) {
	m := New("x", keep.Current{Wav: "/x", RefText: "hi"})
	_, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEsc})
	if cmd == nil {
		t.Fatal("expected quit")
	}
}

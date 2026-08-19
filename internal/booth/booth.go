package booth

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/veronica-agent/cans/internal/play"
	"github.com/veronica-agent/cans/internal/tts"
)

var (
	chrome = lipgloss.NewStyle().Foreground(lipgloss.Color("175"))
	muted  = lipgloss.NewStyle().Foreground(lipgloss.Color("245"))
	errSt  = lipgloss.NewStyle().Foreground(lipgloss.Color("203"))
	box    = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).Padding(1, 2).Width(56)
)

type model struct {
	input  textinput.Model
	status string
	ttfa   string
	quote  string
	busy   bool
	err    string
}

type spokenMsg struct {
	ttfa int
	err  error
}

func New(quote string) model {
	ti := textinput.New()
	ti.Placeholder = "type a line"
	ti.Focus()
	ti.CharLimit = 280
	ti.Width = 48
	if quote == "" {
		quote = "Put the cans on."
	}
	return model{input: ti, status: "listen", quote: quote}
}

func (m model) Init() tea.Cmd { return textinput.Blink }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		case tea.KeyEnter:
			if m.busy {
				return m, nil
			}
			line := strings.TrimSpace(m.input.Value())
			if line == "" {
				return m, nil
			}
			m.busy = true
			m.status = "speaking"
			m.err = ""
			m.input.SetValue("")
			return m, speak(line)
		}
	case spokenMsg:
		m.busy = false
		m.status = "listen"
		if msg.err != nil {
			m.err = msg.err.Error()
			m.ttfa = ""
			return m, nil
		}
		m.ttfa = fmt.Sprintf("%d ms", msg.ttfa)
		return m, nil
	}
	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)
	return m, cmd
}

func speak(line string) tea.Cmd {
	return func() tea.Msg {
		r, err := tts.Say(line)
		if err != nil {
			return spokenMsg{err: err}
		}
		if err := play.File(r.Wav); err != nil {
			return spokenMsg{err: err}
		}
		return spokenMsg{ttfa: r.TTFAMs}
	}
}

func (m model) View() string {
	ttfa := m.ttfa
	if ttfa == "" {
		ttfa = "—"
	}
	head := chrome.Render("cans") + "   " + muted.Render("ttfa "+ttfa)
	body := muted.Render(m.quote) + "\n\n" + m.input.View() + "\n\n" + chrome.Render(m.status)
	if m.err != "" {
		body += "\n" + errSt.Render(m.err)
	}
	return box.Render(head+"\n\n"+body) + "\n" + muted.Render("enter speaks · esc leaves") + "\n"
}

// Run the booth. Blocks.
func Run(quote string) error {
	p := tea.NewProgram(New(quote), tea.WithOutput(os.Stderr))
	_, err := p.Run()
	return err
}

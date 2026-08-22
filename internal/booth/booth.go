package booth

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"
	"github.com/veronica-agent/cans/internal/keep"
	"github.com/veronica-agent/cans/internal/play"
	"github.com/veronica-agent/cans/internal/tts"
)

var (
	rose   = lipgloss.Color("#C45C7A")
	dust   = lipgloss.Color("#A890A0")
	snow   = lipgloss.Color("#F8EDE8")
	chrome = lipgloss.NewStyle().Foreground(rose).Bold(true)
	muted  = lipgloss.NewStyle().Foreground(dust)
	errSt  = lipgloss.NewStyle().Foreground(lipgloss.Color("#E07070"))
	box    = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(rose).
		Padding(1, 2).
		Width(56)
)

type model struct {
	input  textinput.Model
	status string
	ttfa   string
	quote  string
	throat keep.Current
	sess   *tts.Session
	ctx    context.Context
	busy   bool
	err    string
}

type spokenMsg struct {
	ttfa int
	err  error
}

func New(quote string, throat keep.Current, sess *tts.Session, ctx context.Context) model {
	ti := textinput.New()
	ti.Placeholder = "type a line"
	ti.PromptStyle = chrome
	ti.TextStyle = lipgloss.NewStyle().Foreground(snow)
	ti.PlaceholderStyle = muted
	ti.Cursor.Style = chrome
	ti.Focus()
	ti.CharLimit = 280
	ti.Width = 48
	if quote == "" {
		quote = "Put the cans on."
	}
	if ctx == nil {
		ctx = context.Background()
	}
	return model{input: ti, status: "listen", quote: quote, throat: throat, sess: sess, ctx: ctx}
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
			return m, speak(m.ctx, line, m.throat, m.sess)
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

func speak(ctx context.Context, line string, throat keep.Current, sess *tts.Session) tea.Cmd {
	return func() tea.Msg {
		var (
			r   tts.Result
			err error
		)
		if sess != nil {
			r, err = sess.Say(ctx, line, throat)
		} else {
			r, err = tts.SayWith(ctx, line, throat)
		}
		if err != nil {
			return spokenMsg{err: err}
		}
		playErr := play.File(r.Wav)
		tts.RemoveTemp(r.Wav)
		if playErr != nil {
			return spokenMsg{err: playErr}
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

// Run the booth. Throat is frozen for the session. Worker stays warm.
func Run(ctx context.Context, quote string, throat keep.Current) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	lipgloss.SetColorProfile(termenv.TrueColor)
	lipgloss.SetHasDarkBackground(true)
	var sess *tts.Session
	if os.Getenv("CANS_SAY_BIN") == "" {
		s, err := tts.OpenWith(ctx, tts.Options{
			Wait:   -1,
			OnWait: func() { fmt.Fprintln(os.Stderr, "waiting for the mouth…") },
		})
		if err != nil {
			return err
		}
		sess = s
		defer sess.Close()
	}
	p := tea.NewProgram(New(quote, throat, sess, ctx))
	_, err := p.Run()
	return err
}

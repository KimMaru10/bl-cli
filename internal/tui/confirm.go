package tui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type confirmModel struct {
	message  string
	confirmed bool
	done     bool
}

func (m confirmModel) Init() tea.Cmd { return nil }

func (m confirmModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "y", "Y":
			m.confirmed = true
			m.done = true
			return m, tea.Quit
		case "n", "N", "q", "ctrl+c", "esc":
			m.confirmed = false
			m.done = true
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m confirmModel) View() string {
	if m.done {
		return ""
	}
	prompt := lipgloss.NewStyle().Bold(true).Render(m.message)
	return fmt.Sprintf("%s [y/N] ", prompt)
}

// Confirm shows a yes/no confirmation prompt.
// Returns true if the user confirmed.
func Confirm(message string) bool {
	m := confirmModel{message: message}
	p := tea.NewProgram(m)
	finalModel, err := p.Run()
	if err != nil {
		return false
	}
	return finalModel.(confirmModel).confirmed
}

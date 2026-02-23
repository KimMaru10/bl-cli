package tui

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type inputModel struct {
	input    textinput.Model
	value    string
	quitting bool
}

func (m inputModel) Init() tea.Cmd { return textinput.Blink }

func (m inputModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			m.quitting = true
			return m, tea.Quit
		case tea.KeyEnter:
			m.value = m.input.Value()
			return m, tea.Quit
		}
	}
	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)
	return m, cmd
}

func (m inputModel) View() string {
	return m.input.View()
}

// Input shows a text input prompt and returns the entered value.
// Returns empty string and false if the user cancelled.
func Input(prompt, placeholder string) (string, bool) {
	ti := textinput.New()
	ti.Prompt = prompt
	ti.Placeholder = placeholder
	ti.Focus()
	ti.Width = 50

	m := inputModel{input: ti}
	p := tea.NewProgram(m)
	finalModel, err := p.Run()
	if err != nil {
		return "", false
	}

	result := finalModel.(inputModel)
	if result.quitting {
		return "", false
	}
	return result.value, true
}

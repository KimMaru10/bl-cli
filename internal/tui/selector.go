package tui

import (
	"fmt"
	"io"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// SelectItem represents a selectable item with an ID and label.
type SelectItem struct {
	ID    int
	Label string
}

func (i SelectItem) FilterValue() string { return i.Label }

type selectDelegate struct{}

func (d selectDelegate) Height() int                             { return 1 }
func (d selectDelegate) Spacing() int                            { return 0 }
func (d selectDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }

func (d selectDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(SelectItem)
	if !ok {
		return
	}
	if index == m.Index() {
		fmt.Fprint(w, lipgloss.NewStyle().Foreground(lipgloss.Color("6")).Render("> "+i.Label))
	} else {
		fmt.Fprint(w, "  "+i.Label)
	}
}

type selectModel struct {
	list     list.Model
	selected *SelectItem
	quitting bool
}

func (m selectModel) Init() tea.Cmd { return nil }

func (m selectModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			m.quitting = true
			return m, tea.Quit
		case tea.KeyEnter:
			if item, ok := m.list.SelectedItem().(SelectItem); ok {
				m.selected = &item
			}
			return m, tea.Quit
		}
	}
	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m selectModel) View() string {
	if m.selected != nil {
		return ""
	}
	return m.list.View()
}

// Select shows an interactive list and returns the selected item.
// Returns nil if the user cancelled.
func Select(title string, items []SelectItem) *SelectItem {
	listItems := make([]list.Item, len(items))
	for i, item := range items {
		listItems[i] = item
	}

	l := list.New(listItems, selectDelegate{}, 50, min(len(items)+6, 20))
	l.Title = title
	l.SetShowStatusBar(false)
	l.SetShowHelp(true)

	m := selectModel{list: l}
	p := tea.NewProgram(m)
	finalModel, err := p.Run()
	if err != nil {
		return nil
	}

	result := finalModel.(selectModel)
	if result.quitting || result.selected == nil {
		return nil
	}
	return result.selected
}

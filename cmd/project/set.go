package project

import (
	"fmt"
	"io"

	"github.com/KimMaru10/bl-cli/internal/cmdutil"
	"github.com/KimMaru10/bl-cli/internal/config"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

var successStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("2"))

// projectItem implements list.Item for project selection.
type projectItem struct {
	key  string
	name string
}

func (i projectItem) FilterValue() string { return i.key + " " + i.name }

// projectDelegate renders each item in the list.
type projectDelegate struct{}

func (d projectDelegate) Height() int                             { return 1 }
func (d projectDelegate) Spacing() int                            { return 0 }
func (d projectDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }

func (d projectDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(projectItem)
	if !ok {
		return
	}

	str := fmt.Sprintf("%s\t%s", i.key, i.name)
	if index == m.Index() {
		str = lipgloss.NewStyle().Foreground(lipgloss.Color("6")).Render("> " + str)
	} else {
		str = "  " + str
	}
	fmt.Fprint(w, str)
}

type setModel struct {
	list     list.Model
	selected string
	quitting bool
}

func (m setModel) Init() tea.Cmd { return nil }

func (m setModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			m.quitting = true
			return m, tea.Quit
		case tea.KeyEnter:
			if item, ok := m.list.SelectedItem().(projectItem); ok {
				m.selected = item.key
			}
			return m, tea.Quit
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m setModel) View() string {
	if m.selected != "" {
		return ""
	}
	return m.list.View()
}

func newSetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "set",
		Short: "デフォルトプロジェクトを設定する",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, client, err := cmdutil.LoadConfigAndClient()
			if err != nil {
				return err
			}

			projects, err := client.GetProjects()
			if err != nil {
				return err
			}

			items := make([]list.Item, len(projects))
			for i, p := range projects {
				items[i] = projectItem{key: p.ProjectKey, name: p.Name}
			}

			l := list.New(items, projectDelegate{}, 40, 14)
			l.Title = "デフォルトプロジェクトを選択"
			l.SetShowStatusBar(false)
			l.SetShowHelp(true)

			m := setModel{list: l}
			p := tea.NewProgram(m)
			finalModel, err := p.Run()
			if err != nil {
				return fmt.Errorf("TUI の実行に失敗しました: %w", err)
			}

			result := finalModel.(setModel)
			if result.quitting || result.selected == "" {
				return nil
			}

			cfg.DefaultProject = result.selected
			if err := config.Save(cfg); err != nil {
				return fmt.Errorf("設定の保存に失敗しました: %w", err)
			}

			fmt.Println(successStyle.Render("✔ デフォルトプロジェクトを " + result.selected + " に設定しました"))
			return nil
		},
	}
}

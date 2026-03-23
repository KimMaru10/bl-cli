package auth

import (
	"fmt"
	"io"

	"github.com/KimMaru10/bl-cli/internal/config"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

type spaceItem struct {
	alias    string
	url      string
	current  bool
}

func (i spaceItem) FilterValue() string { return i.alias + " " + i.url }

type spaceDelegate struct{}

func (d spaceDelegate) Height() int                             { return 1 }
func (d spaceDelegate) Spacing() int                            { return 0 }
func (d spaceDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }

func (d spaceDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(spaceItem)
	if !ok {
		return
	}

	marker := "  "
	if i.current {
		marker = "* "
	}

	str := fmt.Sprintf("%s%s\t%s", marker, i.alias, i.url)
	if index == m.Index() {
		str = lipgloss.NewStyle().Foreground(lipgloss.Color("6")).Render(str)
	}
	fmt.Fprint(w, str)
}

type switchModel struct {
	list     list.Model
	selected string
	quitting bool
}

func (m switchModel) Init() tea.Cmd { return nil }

func (m switchModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			m.quitting = true
			return m, tea.Quit
		case tea.KeyEnter:
			if item, ok := m.list.SelectedItem().(spaceItem); ok {
				m.selected = item.alias
			}
			return m, tea.Quit
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m switchModel) View() string {
	if m.selected != "" {
		return ""
	}
	return m.list.View()
}

func newSwitchCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "switch",
		Short: "使用するスペースを切り替える",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.Load()
			if err != nil {
				return fmt.Errorf("設定の読み込みに失敗しました: %w", err)
			}

			if len(cfg.Spaces) == 0 {
				fmt.Println(errorStyle.Render("✗ 登録されたスペースがありません。bl auth login を実行してください"))
				return nil
			}

			if len(cfg.Spaces) == 1 {
				fmt.Println(infoStyle.Render("登録されているスペースは1つだけです"))
				return nil
			}

			names := cfg.SpaceNames()
			items := make([]list.Item, len(names))
			for i, name := range names {
				s := cfg.Spaces[name]
				items[i] = spaceItem{
					alias:   name,
					url:     s.SpaceURL,
					current: name == cfg.CurrentSpace,
				}
			}

			l := list.New(items, spaceDelegate{}, 60, min(len(items)+6, 20))
			l.Title = "スペースを選択"
			l.SetShowStatusBar(false)
			l.SetShowHelp(true)

			m := switchModel{list: l}
			p := tea.NewProgram(m)
			finalModel, err := p.Run()
			if err != nil {
				return fmt.Errorf("TUI の実行に失敗しました: %w", err)
			}

			result := finalModel.(switchModel)
			if result.quitting || result.selected == "" {
				return nil
			}

			cfg.CurrentSpace = result.selected
			if err := config.Save(cfg); err != nil {
				return fmt.Errorf("設定の保存に失敗しました: %w", err)
			}

			fmt.Println(successStyle.Render("✔ スペースを " + result.selected + " に切り替えました"))
			return nil
		},
	}
}

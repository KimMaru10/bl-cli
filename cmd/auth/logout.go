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

type logoutItem struct {
	alias string
	url   string
}

func (i logoutItem) FilterValue() string { return i.alias + " " + i.url }

type logoutDelegate struct{}

func (d logoutDelegate) Height() int                             { return 1 }
func (d logoutDelegate) Spacing() int                            { return 0 }
func (d logoutDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }

func (d logoutDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(logoutItem)
	if !ok {
		return
	}

	str := fmt.Sprintf("%s\t%s", i.alias, i.url)
	if index == m.Index() {
		str = lipgloss.NewStyle().Foreground(lipgloss.Color("6")).Render("> " + str)
	} else {
		str = "  " + str
	}
	fmt.Fprint(w, str)
}

type logoutModel struct {
	list     list.Model
	selected string
	quitting bool
}

func (m logoutModel) Init() tea.Cmd { return nil }

func (m logoutModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			m.quitting = true
			return m, tea.Quit
		case tea.KeyEnter:
			if item, ok := m.list.SelectedItem().(logoutItem); ok {
				m.selected = item.alias
			}
			return m, tea.Quit
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m logoutModel) View() string {
	if m.selected != "" {
		return ""
	}
	return m.list.View()
}

func newLogoutCmd() *cobra.Command {
	var all bool

	cmd := &cobra.Command{
		Use:   "logout",
		Short: "認証情報を削除する",
		RunE: func(cmd *cobra.Command, args []string) error {
			if all {
				if err := config.Delete(); err != nil {
					return fmt.Errorf("認証情報の削除に失敗しました: %w", err)
				}
				fmt.Println(successStyle.Render("✔ すべての認証情報を削除しました"))
				return nil
			}

			cfg, err := config.Load()
			if err != nil {
				return fmt.Errorf("設定の読み込みに失敗しました: %w", err)
			}

			if len(cfg.Spaces) == 0 {
				fmt.Println(infoStyle.Render("登録されたスペースがありません"))
				return nil
			}

			// If only one space, delete it directly
			if len(cfg.Spaces) == 1 {
				if err := config.Delete(); err != nil {
					return fmt.Errorf("認証情報の削除に失敗しました: %w", err)
				}
				fmt.Println(successStyle.Render("✔ 認証情報を削除しました"))
				return nil
			}

			// Multiple spaces: show selector
			names := cfg.SpaceNames()
			items := make([]list.Item, len(names))
			for i, name := range names {
				s := cfg.Spaces[name]
				items[i] = logoutItem{alias: name, url: s.SpaceURL}
			}

			l := list.New(items, logoutDelegate{}, 60, min(len(items)+6, 20))
			l.Title = "削除するスペースを選択"
			l.SetShowStatusBar(false)
			l.SetShowHelp(true)

			m := logoutModel{list: l}
			p := tea.NewProgram(m)
			finalModel, err := p.Run()
			if err != nil {
				return fmt.Errorf("TUI の実行に失敗しました: %w", err)
			}

			result := finalModel.(logoutModel)
			if result.quitting || result.selected == "" {
				return nil
			}

			delete(cfg.Spaces, result.selected)

			// If we deleted the current space, switch to another
			if cfg.CurrentSpace == result.selected {
				cfg.CurrentSpace = ""
				for name := range cfg.Spaces {
					cfg.CurrentSpace = name
					break
				}
			}

			if len(cfg.Spaces) == 0 {
				if err := config.Delete(); err != nil {
					return fmt.Errorf("認証情報の削除に失敗しました: %w", err)
				}
			} else {
				if err := config.Save(cfg); err != nil {
					return fmt.Errorf("設定の保存に失敗しました: %w", err)
				}
			}

			fmt.Println(successStyle.Render("✔ スペース " + result.selected + " の認証情報を削除しました"))
			return nil
		},
	}

	cmd.Flags().BoolVar(&all, "all", false, "すべてのスペースの認証情報を削除する")
	return cmd
}

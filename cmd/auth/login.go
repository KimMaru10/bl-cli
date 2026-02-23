package auth

import (
	"fmt"
	"strings"

	"github.com/KimMaru10/bl-cli/internal/api"
	"github.com/KimMaru10/bl-cli/internal/config"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

var (
	successStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("2"))
	errorStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("1"))
)

type loginStep int

const (
	stepSpaceURL loginStep = iota
	stepAPIKey
	stepDone
)

type loginModel struct {
	step     loginStep
	inputs   []textinput.Model
	err      error
	quitting bool
}

func newLoginModel() loginModel {
	spaceInput := textinput.New()
	spaceInput.Placeholder = "https://myteam.backlog.com"
	spaceInput.Focus()
	spaceInput.Width = 50
	spaceInput.Prompt = "Backlog スペース URL: "

	apiKeyInput := textinput.New()
	apiKeyInput.Placeholder = "APIキーを入力"
	apiKeyInput.Width = 50
	apiKeyInput.EchoMode = textinput.EchoPassword
	apiKeyInput.Prompt = "API キー: "

	return loginModel{
		step:   stepSpaceURL,
		inputs: []textinput.Model{spaceInput, apiKeyInput},
	}
}

func (m loginModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m loginModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			m.quitting = true
			return m, tea.Quit
		case tea.KeyEnter:
			if m.step == stepSpaceURL {
				m.step = stepAPIKey
				m.inputs[0].Blur()
				m.inputs[1].Focus()
				return m, textinput.Blink
			}
			if m.step == stepAPIKey {
				m.step = stepDone
				return m, tea.Quit
			}
		}
	}

	var cmd tea.Cmd
	m.inputs[m.step], cmd = m.inputs[m.step].Update(msg)
	return m, cmd
}

func (m loginModel) View() string {
	var b strings.Builder

	if m.step >= stepSpaceURL {
		b.WriteString(m.inputs[0].View())
		b.WriteString("\n")
	}
	if m.step >= stepAPIKey {
		b.WriteString(m.inputs[1].View())
		b.WriteString("\n")
	}

	return b.String()
}

func newLoginCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "login",
		Short: "Backlog に認証する",
		RunE: func(cmd *cobra.Command, args []string) error {
			m := newLoginModel()
			p := tea.NewProgram(m)
			finalModel, err := p.Run()
			if err != nil {
				return fmt.Errorf("TUI の実行に失敗しました: %w", err)
			}

			result := finalModel.(loginModel)
			if result.quitting {
				return nil
			}

			spaceURL := strings.TrimRight(strings.TrimSpace(result.inputs[0].Value()), "/")
			apiKey := strings.TrimSpace(result.inputs[1].Value())

			if spaceURL == "" || apiKey == "" {
				fmt.Println(errorStyle.Render("✗ スペースURLとAPIキーは必須です"))
				return nil
			}

			client := api.NewClient(spaceURL, apiKey)
			user, err := client.GetMyself()
			if err != nil {
				fmt.Println(errorStyle.Render("✗ 認証に失敗しました: " + err.Error()))
				return nil
			}

			cfg := &config.Config{
				SpaceURL: spaceURL,
				APIKey:   apiKey,
			}

			existing, _ := config.Load()
			if existing != nil && existing.DefaultProject != "" {
				cfg.DefaultProject = existing.DefaultProject
			}

			if err := config.Save(cfg); err != nil {
				return fmt.Errorf("設定の保存に失敗しました: %w", err)
			}

			fmt.Println(successStyle.Render("✔ " + user.Name + " として認証しました"))
			return nil
		},
	}
}

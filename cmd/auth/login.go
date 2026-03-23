package auth

import (
	"fmt"
	"strings"

	"github.com/KimMaru10/bl-cli/internal/api"
	"github.com/KimMaru10/bl-cli/internal/config"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

type loginStep int

const (
	stepSpaceURL loginStep = iota
	stepAPIKey
	stepAlias
	stepDone
)

type loginModel struct {
	step     loginStep
	inputs   []textinput.Model
	err      error
	quitting bool
}

func newLoginModel(suggestedAlias string) loginModel {
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

	aliasInput := textinput.New()
	aliasInput.Placeholder = suggestedAlias
	aliasInput.Width = 50
	aliasInput.Prompt = "スペースのエイリアス名: "

	return loginModel{
		step:   stepSpaceURL,
		inputs: []textinput.Model{spaceInput, apiKeyInput, aliasInput},
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
			switch m.step {
			case stepSpaceURL:
				m.step = stepAPIKey
				m.inputs[0].Blur()
				m.inputs[1].Focus()
				return m, textinput.Blink
			case stepAPIKey:
				// Derive suggested alias from URL for placeholder
				spaceURL := strings.TrimRight(strings.TrimSpace(m.inputs[0].Value()), "/")
				suggested := config.ExtractAlias(spaceURL)
				m.inputs[2].Placeholder = suggested

				m.step = stepAlias
				m.inputs[1].Blur()
				m.inputs[2].Focus()
				return m, textinput.Blink
			case stepAlias:
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
	if m.step >= stepAlias {
		b.WriteString(m.inputs[2].View())
		b.WriteString("\n")
	}

	return b.String()
}

func newLoginCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "login",
		Short: "Backlog スペースを追加認証する",
		RunE: func(cmd *cobra.Command, args []string) error {
			m := newLoginModel("myteam")
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
			alias := strings.TrimSpace(result.inputs[2].Value())

			if spaceURL == "" || apiKey == "" {
				fmt.Println(errorStyle.Render("✗ スペースURLとAPIキーは必須です"))
				return nil
			}

			// Use suggested alias if empty
			if alias == "" {
				alias = config.ExtractAlias(spaceURL)
			}

			// Validate credentials
			client := api.NewClient(spaceURL, apiKey)
			user, err := client.GetMyself()
			if err != nil {
				fmt.Println(errorStyle.Render("✗ 認証に失敗しました: " + err.Error()))
				return nil
			}

			// Load existing config and add/update space
			cfg, _ := config.Load()
			if cfg.Spaces == nil {
				cfg.Spaces = make(map[string]config.SpaceConfig)
			}

			cfg.Spaces[alias] = config.SpaceConfig{
				SpaceURL:       spaceURL,
				APIKey:         apiKey,
				DefaultProject: cfg.Spaces[alias].DefaultProject, // preserve existing default_project
			}

			// Set as current if it's the first space or no current is set
			if cfg.CurrentSpace == "" || len(cfg.Spaces) == 1 {
				cfg.CurrentSpace = alias
			}

			if err := config.Save(cfg); err != nil {
				return fmt.Errorf("設定の保存に失敗しました: %w", err)
			}

			fmt.Println(successStyle.Render(fmt.Sprintf("✔ %s として認証しました（スペース: %s）", user.Name, alias)))
			if cfg.CurrentSpace == alias {
				fmt.Println(infoStyle.Render("  現在のスペースに設定されました"))
			}
			return nil
		},
	}
}

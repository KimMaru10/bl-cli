package issue

import (
	"fmt"
	"strings"

	"github.com/KimMaru10/bl-cli/internal/browser"
	"github.com/KimMaru10/bl-cli/internal/cmdutil"
	"github.com/KimMaru10/bl-cli/internal/git"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

var (
	titleStyle   = lipgloss.NewStyle().Bold(true)
	labelStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("8"))
	urlStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("4")).Underline(true)
)

// resolveIssueKey resolves an issue key from args or the current git branch.
func resolveIssueKey(args []string) (string, error) {
	if len(args) > 0 {
		return args[0], nil
	}
	branch, err := git.GetCurrentBranch()
	if err != nil {
		return "", fmt.Errorf("課題キーを指定するか、課題キーを含むブランチに切り替えてください")
	}
	key := git.ExtractIssueKey(branch)
	if key == "" {
		return "", fmt.Errorf("課題キーを指定するか、課題キーを含むブランチに切り替えてください")
	}
	return key, nil
}

func newViewCmd() *cobra.Command {
	var web bool

	cmd := &cobra.Command{
		Use:   "view [issueKey]",
		Short: "課題の詳細を表示する",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, client, err := cmdutil.LoadConfigAndClient()
			if err != nil {
				return err
			}

			issueKey, err := resolveIssueKey(args)
			if err != nil {
				return err
			}

			if web {
				url := cfg.SpaceURL + "/view/" + issueKey
				return browser.Open(url)
			}

			issue, err := client.GetIssue(issueKey)
			if err != nil {
				return err
			}

			// Title
			fmt.Println(titleStyle.Render(issue.IssueKey + " " + issue.Summary))
			fmt.Println()

			// Status | Priority | IssueType
			var meta []string
			if issue.Status != nil {
				meta = append(meta, statusColor(issue.Status.Name).Render(issue.Status.Name))
			}
			if issue.Priority != nil {
				meta = append(meta, issue.Priority.Name)
			}
			if issue.IssueType != nil {
				meta = append(meta, issue.IssueType.Name)
			}
			if len(meta) > 0 {
				fmt.Println(strings.Join(meta, " | "))
			}

			// Assignee | CreatedUser
			var people []string
			if issue.Assignee != nil {
				people = append(people, labelStyle.Render("担当者: ")+issue.Assignee.Name)
			}
			if issue.CreatedUser != nil {
				people = append(people, labelStyle.Render("作成者: ")+issue.CreatedUser.Name)
			}
			if len(people) > 0 {
				fmt.Println(strings.Join(people, " | "))
			}

			// Due date
			if issue.DueDate != "" {
				fmt.Println(labelStyle.Render("期日: ") + issue.DueDate)
			}

			// Milestones
			if len(issue.Milestone) > 0 {
				var names []string
				for _, m := range issue.Milestone {
					names = append(names, m.Name)
				}
				fmt.Println(labelStyle.Render("マイルストーン: ") + strings.Join(names, ", "))
			}

			// Description
			if issue.Description != "" {
				fmt.Println()
				fmt.Println(issue.Description)
			}

			// URL
			fmt.Println()
			fmt.Println(labelStyle.Render("URL: ") + urlStyle.Render(cfg.SpaceURL+"/view/"+issue.IssueKey))

			return nil
		},
	}

	cmd.Flags().BoolVarP(&web, "web", "w", false, "ブラウザで開く")

	return cmd
}

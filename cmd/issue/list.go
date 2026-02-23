package issue

import (
	"fmt"
	"strings"

	"github.com/KimMaru10/bl-cli/internal/api"
	"github.com/KimMaru10/bl-cli/internal/browser"
	"github.com/KimMaru10/bl-cli/internal/cmdutil"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

func statusColor(name string) lipgloss.Style {
	switch {
	case strings.Contains(name, "未対応"):
		return lipgloss.NewStyle().Foreground(lipgloss.Color("1")) // red
	case strings.Contains(name, "処理中"):
		return lipgloss.NewStyle().Foreground(lipgloss.Color("4")) // blue
	case strings.Contains(name, "処理済み"):
		return lipgloss.NewStyle().Foreground(lipgloss.Color("3")) // yellow
	case strings.Contains(name, "完了"):
		return lipgloss.NewStyle().Foreground(lipgloss.Color("2")) // green
	default:
		return lipgloss.NewStyle()
	}
}

func newListCmd() *cobra.Command {
	var (
		assignee  string
		status    string
		milestone string
		project   string
		count     int
		web       bool
	)

	cmd := &cobra.Command{
		Use:   "list",
		Short: "課題一覧を表示する",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, client, err := cmdutil.LoadConfigAndClient()
			if err != nil {
				return err
			}

			projectKey := project
			if projectKey == "" {
				projectKey = cfg.DefaultProject
			}
			if projectKey == "" {
				return fmt.Errorf("プロジェクトを指定してください（--project または bl project set）")
			}

			if web {
				url := cfg.SpaceURL + "/find/" + projectKey
				return browser.Open(url)
			}

			proj, err := client.GetProject(projectKey)
			if err != nil {
				return err
			}

			opts := &api.GetIssuesOptions{
				ProjectIDs: []int{proj.ID},
				Count:      count,
				Sort:       "updated",
				Order:      "desc",
			}

			if assignee == "@me" {
				me, err := client.GetMyself()
				if err != nil {
					return err
				}
				opts.AssigneeIDs = []int{me.ID}
			} else if assignee != "" {
				users, err := client.GetProjectUsers(projectKey)
				if err != nil {
					return err
				}
				for _, u := range users {
					if u.Name == assignee {
						opts.AssigneeIDs = []int{u.ID}
						break
					}
				}
			}

			if status != "" {
				statuses, err := client.GetStatuses(projectKey)
				if err != nil {
					return err
				}
				for _, s := range statuses {
					if s.Name == status {
						opts.StatusIDs = []int{s.ID}
						break
					}
				}
			}

			if milestone != "" {
				milestones, err := client.GetMilestones(projectKey)
				if err != nil {
					return err
				}
				for _, m := range milestones {
					if m.Name == milestone {
						opts.MilestoneIDs = []int{m.ID}
						break
					}
				}
			}

			issues, err := client.GetIssues(opts)
			if err != nil {
				return err
			}

			if len(issues) == 0 {
				fmt.Println("該当する課題はありません")
				return nil
			}

			headerStyle := lipgloss.NewStyle().Bold(true)
			fmt.Printf("%s\t%s\t%s\t%s\n",
				headerStyle.Render("KEY"),
				headerStyle.Render("STATUS"),
				headerStyle.Render("ASSIGNEE"),
				headerStyle.Render("TITLE"),
			)

			for _, issue := range issues {
				statusName := ""
				if issue.Status != nil {
					statusName = statusColor(issue.Status.Name).Render(issue.Status.Name)
				}

				assigneeName := ""
				if issue.Assignee != nil {
					assigneeName = issue.Assignee.Name
				}

				fmt.Printf("%s\t%s\t%s\t%s\n",
					issue.IssueKey,
					statusName,
					assigneeName,
					issue.Summary,
				)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&assignee, "assignee", "a", "", "担当者名（@me で自分）")
	cmd.Flags().StringVarP(&status, "status", "s", "", "ステータス名")
	cmd.Flags().StringVarP(&milestone, "milestone", "m", "", "マイルストーン名")
	cmd.Flags().StringVarP(&project, "project", "p", "", "プロジェクトキー")
	cmd.Flags().IntVarP(&count, "count", "c", 20, "表示件数")
	cmd.Flags().BoolVarP(&web, "web", "w", false, "ブラウザで開く")

	return cmd
}

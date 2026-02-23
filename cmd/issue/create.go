package issue

import (
	"fmt"

	"github.com/KimMaru10/bl-cli/internal/api"
	"github.com/KimMaru10/bl-cli/internal/cmdutil"
	"github.com/KimMaru10/bl-cli/internal/tui"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

var successStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("2"))

func newCreateCmd() *cobra.Command {
	var (
		summary     string
		typeName    string
		priority    string
		assignee    string
		description string
		dueDate     string
		milestone   string
		project     string
	)

	cmd := &cobra.Command{
		Use:   "create",
		Short: "課題を作成する",
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

			proj, err := client.GetProject(projectKey)
			if err != nil {
				return err
			}

			opts := &api.CreateIssueOptions{
				ProjectID: proj.ID,
			}

			if summary != "" {
				// Flag mode
				if typeName == "" || priority == "" {
					return fmt.Errorf("--type と --priority は必須です")
				}

				opts.Summary = summary
				opts.Description = description
				opts.DueDate = dueDate

				issueTypes, err := client.GetIssueTypes(projectKey)
				if err != nil {
					return err
				}
				for _, t := range issueTypes {
					if t.Name == typeName {
						opts.IssueTypeID = t.ID
						break
					}
				}
				if opts.IssueTypeID == 0 {
					return fmt.Errorf("課題種別 '%s' が見つかりません", typeName)
				}

				priorities, err := client.GetPriorities()
				if err != nil {
					return err
				}
				for _, p := range priorities {
					if p.Name == priority {
						opts.PriorityID = p.ID
						break
					}
				}
				if opts.PriorityID == 0 {
					return fmt.Errorf("優先度 '%s' が見つかりません", priority)
				}

				if assignee != "" {
					users, err := client.GetProjectUsers(projectKey)
					if err != nil {
						return err
					}
					for _, u := range users {
						if u.Name == assignee {
							opts.AssigneeID = u.ID
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
			} else {
				// Interactive mode
				s, ok := tui.Input("タイトル: ", "課題のタイトルを入力")
				if !ok || s == "" {
					return nil
				}
				opts.Summary = s

				// Issue type
				issueTypes, err := client.GetIssueTypes(projectKey)
				if err != nil {
					return err
				}
				typeItems := make([]tui.SelectItem, len(issueTypes))
				for i, t := range issueTypes {
					typeItems[i] = tui.SelectItem{ID: t.ID, Label: t.Name}
				}
				selected := tui.Select("課題種別を選択", typeItems)
				if selected == nil {
					return nil
				}
				opts.IssueTypeID = selected.ID

				// Priority
				priorities, err := client.GetPriorities()
				if err != nil {
					return err
				}
				prioItems := make([]tui.SelectItem, len(priorities))
				for i, p := range priorities {
					prioItems[i] = tui.SelectItem{ID: p.ID, Label: p.Name}
				}
				selected = tui.Select("優先度を選択", prioItems)
				if selected == nil {
					return nil
				}
				opts.PriorityID = selected.ID

				// Assignee
				users, err := client.GetProjectUsers(projectKey)
				if err != nil {
					return err
				}
				userItems := []tui.SelectItem{{ID: 0, Label: "未設定"}}
				for _, u := range users {
					userItems = append(userItems, tui.SelectItem{ID: u.ID, Label: u.Name})
				}
				selected = tui.Select("担当者を選択", userItems)
				if selected == nil {
					return nil
				}
				if selected.ID > 0 {
					opts.AssigneeID = selected.ID
				}

				// Due date
				d, ok := tui.Input("期日 (yyyy-MM-dd, 空欄で省略): ", "2025-12-31")
				if !ok {
					return nil
				}
				opts.DueDate = d

				// Description
				desc, ok := tui.Input("説明 (空欄で省略): ", "")
				if !ok {
					return nil
				}
				opts.Description = desc

				// Confirm
				if !tui.Confirm("この内容で課題を作成しますか？") {
					fmt.Println("キャンセルしました")
					return nil
				}
			}

			issue, err := client.CreateIssue(opts)
			if err != nil {
				return err
			}

			fmt.Println(successStyle.Render("✔ " + issue.IssueKey + " を作成しました"))
			return nil
		},
	}

	cmd.Flags().StringVarP(&summary, "summary", "s", "", "タイトル")
	cmd.Flags().StringVarP(&typeName, "type", "t", "", "課題種別名")
	cmd.Flags().StringVar(&priority, "priority", "", "優先度名")
	cmd.Flags().StringVarP(&assignee, "assignee", "a", "", "担当者名")
	cmd.Flags().StringVarP(&description, "description", "d", "", "説明")
	cmd.Flags().StringVar(&dueDate, "due-date", "", "期日（yyyy-MM-dd）")
	cmd.Flags().StringVarP(&milestone, "milestone", "m", "", "マイルストーン名")
	cmd.Flags().StringVarP(&project, "project", "p", "", "プロジェクトキー")

	return cmd
}

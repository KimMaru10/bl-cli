package issue

import (
	"fmt"
	"strings"

	"github.com/KimMaru10/bl-cli/internal/api"
	"github.com/KimMaru10/bl-cli/internal/cmdutil"
	"github.com/KimMaru10/bl-cli/internal/tui"
	"github.com/spf13/cobra"
)

func intPtr(v int) *int       { return &v }
func strPtr(v string) *string { return &v }

func newEditCmd() *cobra.Command {
	var (
		status    string
		assignee  string
		dueDate   string
		priority  string
		milestone string
		comment   string
	)

	cmd := &cobra.Command{
		Use:   "edit [issueKey]",
		Short: "課題を更新する",
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

			// Get current issue to determine project key
			currentIssue, err := client.GetIssue(issueKey)
			if err != nil {
				return err
			}

			// Extract project key from issue key (e.g. "TEST-1" -> "TEST")
			projectKey := issueKey[:strings.Index(issueKey, "-")]
			_ = cfg

			hasFlags := cmd.Flags().Changed("status") || cmd.Flags().Changed("assignee") ||
				cmd.Flags().Changed("due-date") || cmd.Flags().Changed("priority") ||
				cmd.Flags().Changed("milestone") || cmd.Flags().Changed("comment")

			opts := &api.UpdateIssueOptions{}

			if hasFlags {
				// Flag mode
				if status != "" {
					statuses, err := client.GetStatuses(projectKey)
					if err != nil {
						return err
					}
					for _, s := range statuses {
						if s.Name == status {
							opts.StatusID = intPtr(s.ID)
							break
						}
					}
					if opts.StatusID == nil {
						return fmt.Errorf("ステータス '%s' が見つかりません", status)
					}
				}

				if assignee != "" {
					users, err := client.GetProjectUsers(projectKey)
					if err != nil {
						return err
					}
					for _, u := range users {
						if u.Name == assignee {
							opts.AssigneeID = intPtr(u.ID)
							break
						}
					}
				}

				if dueDate != "" {
					opts.DueDate = strPtr(dueDate)
				}

				if priority != "" {
					priorities, err := client.GetPriorities()
					if err != nil {
						return err
					}
					for _, p := range priorities {
						if p.Name == priority {
							opts.PriorityID = intPtr(p.ID)
							break
						}
					}
					if opts.PriorityID == nil {
						return fmt.Errorf("優先度 '%s' が見つかりません", priority)
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

				if comment != "" {
					opts.Comment = strPtr(comment)
				}
			} else {
				// Interactive mode
				editItems := []tui.SelectItem{
					{ID: 1, Label: "ステータスを変更"},
					{ID: 2, Label: "担当者を変更"},
					{ID: 3, Label: "期日を変更"},
					{ID: 4, Label: "優先度を変更"},
					{ID: 5, Label: "マイルストーンを変更"},
				}

				selected := tui.Select("編集項目を選択", editItems)
				if selected == nil {
					return nil
				}

				switch selected.ID {
				case 1: // Status
					statuses, err := client.GetStatuses(projectKey)
					if err != nil {
						return err
					}
					items := make([]tui.SelectItem, len(statuses))
					for i, s := range statuses {
						label := s.Name
						if currentIssue.Status != nil && s.ID == currentIssue.Status.ID {
							label += " (現在)"
						}
						items[i] = tui.SelectItem{ID: s.ID, Label: label}
					}
					sel := tui.Select("ステータスを選択", items)
					if sel == nil {
						return nil
					}
					opts.StatusID = intPtr(sel.ID)

				case 2: // Assignee
					users, err := client.GetProjectUsers(projectKey)
					if err != nil {
						return err
					}
					items := []tui.SelectItem{{ID: 0, Label: "未設定"}}
					for _, u := range users {
						label := u.Name
						if currentIssue.Assignee != nil && u.ID == currentIssue.Assignee.ID {
							label += " (現在)"
						}
						items = append(items, tui.SelectItem{ID: u.ID, Label: label})
					}
					sel := tui.Select("担当者を選択", items)
					if sel == nil {
						return nil
					}
					opts.AssigneeID = intPtr(sel.ID)

				case 3: // Due date
					current := ""
					if currentIssue.DueDate != "" {
						current = currentIssue.DueDate
					}
					val, ok := tui.Input(fmt.Sprintf("期日 (現在: %s): ", current), "yyyy-MM-dd")
					if !ok {
						return nil
					}
					opts.DueDate = strPtr(val)

				case 4: // Priority
					priorities, err := client.GetPriorities()
					if err != nil {
						return err
					}
					items := make([]tui.SelectItem, len(priorities))
					for i, p := range priorities {
						label := p.Name
						if currentIssue.Priority != nil && p.ID == currentIssue.Priority.ID {
							label += " (現在)"
						}
						items[i] = tui.SelectItem{ID: p.ID, Label: label}
					}
					sel := tui.Select("優先度を選択", items)
					if sel == nil {
						return nil
					}
					opts.PriorityID = intPtr(sel.ID)

				case 5: // Milestone
					milestones, err := client.GetMilestones(projectKey)
					if err != nil {
						return err
					}
					items := []tui.SelectItem{{ID: 0, Label: "未設定"}}
					for _, m := range milestones {
						items = append(items, tui.SelectItem{ID: m.ID, Label: m.Name})
					}
					sel := tui.Select("マイルストーンを選択", items)
					if sel == nil {
						return nil
					}
					if sel.ID > 0 {
						opts.MilestoneIDs = []int{sel.ID}
					} else {
						opts.MilestoneIDs = []int{}
					}
				}

				if !tui.Confirm("この内容で更新しますか？") {
					fmt.Println("キャンセルしました")
					return nil
				}
			}

			_, err = client.UpdateIssue(issueKey, opts)
			if err != nil {
				return err
			}

			fmt.Println(successStyle.Render("✔ " + issueKey + " を更新しました"))
			return nil
		},
	}

	cmd.Flags().StringVar(&status, "status", "", "ステータス名")
	cmd.Flags().StringVarP(&assignee, "assignee", "a", "", "担当者名")
	cmd.Flags().StringVar(&dueDate, "due-date", "", "期日（yyyy-MM-dd）")
	cmd.Flags().StringVar(&priority, "priority", "", "優先度名")
	cmd.Flags().StringVarP(&milestone, "milestone", "m", "", "マイルストーン名")
	cmd.Flags().StringVar(&comment, "comment", "", "更新時コメント")

	return cmd
}

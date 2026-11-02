package project

import (
	"fmt"

	"github.com/KimMaru10/bl-cli/internal/cmdutil"
	"github.com/spf13/cobra"
)

func newMetaCmd() *cobra.Command {
	var project string
	cmd := &cobra.Command{
		Use:   "meta",
		Short: "ステータス・課題種別・カテゴリ・マイルストーンの一覧を表示する",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, client, err := cmdutil.LoadConfigAndClient()
			if err != nil {
				return err
			}
			key, err := cmdutil.ResolveProjectKey(cfg, project)
			if err != nil {
				return err
			}

			statuses, err := client.GetStatuses(key)
			if err != nil {
				return err
			}
			fmt.Println(headerStyle.Render("ステータス"))
			for _, s := range statuses {
				fmt.Printf("  %d\t%s\n", s.ID, s.Name)
			}

			types, err := client.GetIssueTypes(key)
			if err != nil {
				return err
			}
			fmt.Println(headerStyle.Render("課題種別"))
			for _, t := range types {
				fmt.Printf("  %d\t%s\n", t.ID, t.Name)
			}

			categories, err := client.GetCategories(key)
			if err != nil {
				return err
			}
			fmt.Println(headerStyle.Render("カテゴリ"))
			for _, c := range categories {
				fmt.Printf("  %d\t%s\n", c.ID, c.Name)
			}

			milestones, err := client.GetMilestones(key)
			if err != nil {
				return err
			}
			fmt.Println(headerStyle.Render("マイルストーン"))
			for _, m := range milestones {
				fmt.Printf("  %d\t%s\t%s\n", m.ID, m.Name, m.ReleaseDueDate)
			}
			return nil
		},
	}
	cmd.Flags().StringVarP(&project, "project", "p", "", "プロジェクトキー")
	return cmd
}

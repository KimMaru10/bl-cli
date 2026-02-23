package project

import (
	"fmt"

	"github.com/KimMaru10/bl-cli/internal/cmdutil"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

var headerStyle = lipgloss.NewStyle().Bold(true)

func newListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "プロジェクト一覧を表示する",
		RunE: func(cmd *cobra.Command, args []string) error {
			_, client, err := cmdutil.LoadConfigAndClient()
			if err != nil {
				return err
			}

			projects, err := client.GetProjects()
			if err != nil {
				return err
			}

			if len(projects) == 0 {
				fmt.Println("プロジェクトがありません")
				return nil
			}

			fmt.Printf("%s\t%s\n", headerStyle.Render("KEY"), headerStyle.Render("NAME"))
			for _, p := range projects {
				fmt.Printf("%s\t%s\n", p.ProjectKey, p.Name)
			}
			return nil
		},
	}
}

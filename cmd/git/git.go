package git

import (
	"fmt"
	"strconv"

	"github.com/KimMaru10/bl-cli/internal/cmdutil"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

var (
	headerStyle = lipgloss.NewStyle().Bold(true)
	titleStyle  = lipgloss.NewStyle().Bold(true)
	labelStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("8"))
)

// NewRepoCmd returns the repo subcommand group (Backlog Git).
func NewRepoCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "repo",
		Short: "Backlog Git リポジトリの参照",
	}
	var project string
	list := &cobra.Command{
		Use:   "list",
		Short: "リポジトリ一覧を表示する",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, client, err := cmdutil.LoadConfigAndClient()
			if err != nil {
				return err
			}
			key, err := cmdutil.ResolveProjectKey(cfg, project)
			if err != nil {
				return err
			}
			repos, err := client.GetRepositories(key)
			if err != nil {
				return err
			}
			if len(repos) == 0 {
				fmt.Println("リポジトリはありません")
				return nil
			}
			fmt.Printf("%s\t%s\t%s\n", headerStyle.Render("NAME"), headerStyle.Render("UPDATED"), headerStyle.Render("DESCRIPTION"))
			for _, r := range repos {
				fmt.Printf("%s\t%s\t%s\n", r.Name, r.Updated, cmdutil.StripControl(r.Description))
			}
			return nil
		},
	}
	list.Flags().StringVarP(&project, "project", "p", "", "プロジェクトキー")
	cmd.AddCommand(list)
	return cmd
}

// NewPRCmd returns the pr subcommand group (Backlog Git pull requests).
func NewPRCmd() *cobra.Command {
	var project, repo string
	cmd := &cobra.Command{
		Use:   "pr",
		Short: "Backlog Git プルリクエストの参照",
	}
	cmd.PersistentFlags().StringVarP(&project, "project", "p", "", "プロジェクトキー")
	cmd.PersistentFlags().StringVarP(&repo, "repo", "r", "", "リポジトリ名（必須、一覧は bl repo list）")
	cmd.MarkPersistentFlagRequired("repo")

	var count int
	list := &cobra.Command{
		Use:   "list",
		Short: "プルリクエスト一覧を表示する",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, client, err := cmdutil.LoadConfigAndClient()
			if err != nil {
				return err
			}
			key, err := cmdutil.ResolveProjectKey(cfg, project)
			if err != nil {
				return err
			}
			prs, err := client.GetPullRequests(key, repo, count)
			if err != nil {
				return err
			}
			if len(prs) == 0 {
				fmt.Println("プルリクエストはありません")
				return nil
			}
			fmt.Printf("%s\t%s\t%s\t%s\n", headerStyle.Render("#"), headerStyle.Render("STATUS"), headerStyle.Render("BRANCH"), headerStyle.Render("TITLE"))
			for _, p := range prs {
				status := ""
				if p.Status != nil {
					status = p.Status.Name
				}
				fmt.Printf("%d\t%s\t%s→%s\t%s\n", p.Number, status, cmdutil.StripControl(p.Branch), cmdutil.StripControl(p.Base), cmdutil.StripControl(p.Summary))
			}
			return nil
		},
	}
	list.Flags().IntVarP(&count, "count", "c", 20, "表示件数（最大 100）")

	view := &cobra.Command{
		Use:   "view <number>",
		Short: "プルリクエストの詳細を表示する",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			number, err := strconv.Atoi(args[0])
			if err != nil {
				return fmt.Errorf("プルリクエスト番号は数値で指定してください: %s", args[0])
			}
			cfg, client, err := cmdutil.LoadConfigAndClient()
			if err != nil {
				return err
			}
			key, err := cmdutil.ResolveProjectKey(cfg, project)
			if err != nil {
				return err
			}
			p, err := client.GetPullRequest(key, repo, number)
			if err != nil {
				return err
			}
			fmt.Println(titleStyle.Render(fmt.Sprintf("#%d %s", p.Number, cmdutil.StripControl(p.Summary))))
			if p.Status != nil {
				fmt.Println(labelStyle.Render("状態: ") + p.Status.Name)
			}
			fmt.Println(labelStyle.Render("ブランチ: ") + cmdutil.StripControl(p.Branch+" → "+p.Base))
			if p.Assignee != nil {
				fmt.Println(labelStyle.Render("担当者: ") + p.Assignee.Name)
			}
			if p.Issue != nil {
				issue := p.Issue.IssueKey
				if issue == "" {
					issue = "ID " + strconv.Itoa(p.Issue.ID)
				}
				fmt.Println(labelStyle.Render("関連課題: ") + issue)
			}
			fmt.Println(labelStyle.Render("作成: ") + p.Created + labelStyle.Render("  更新: ") + p.Updated)
			if p.Description != "" {
				fmt.Println()
				fmt.Println(cmdutil.StripControl(p.Description))
			}
			return nil
		},
	}

	cmd.AddCommand(list)
	cmd.AddCommand(view)
	return cmd
}

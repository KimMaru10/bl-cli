package wiki

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/KimMaru10/bl-cli/internal/api"
	"github.com/KimMaru10/bl-cli/internal/cmdutil"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

var (
	headerStyle = lipgloss.NewStyle().Bold(true)
	titleStyle  = lipgloss.NewStyle().Bold(true)
	labelStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("8"))
)

// NewWikiCmd returns the wiki subcommand group.
func NewWikiCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "wiki",
		Short: "Wiki の参照",
	}
	cmd.AddCommand(newListCmd())
	cmd.AddCommand(newViewCmd())
	return cmd
}

func newListCmd() *cobra.Command {
	var project, keyword string
	cmd := &cobra.Command{
		Use:   "list",
		Short: "Wiki ページ一覧を表示する",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, client, err := cmdutil.LoadConfigAndClient()
			if err != nil {
				return err
			}
			key, err := cmdutil.ResolveProjectKey(cfg, project)
			if err != nil {
				return err
			}
			wikis, err := client.GetWikis(key, keyword)
			if err != nil {
				return err
			}
			if len(wikis) == 0 {
				fmt.Println("該当する Wiki ページはありません")
				return nil
			}
			fmt.Printf("%s\t%s\t%s\n", headerStyle.Render("ID"), headerStyle.Render("UPDATED"), headerStyle.Render("NAME"))
			for _, w := range wikis {
				fmt.Printf("%d\t%s\t%s\n", w.ID, w.Updated, cmdutil.StripControl(w.Name))
			}
			return nil
		},
	}
	cmd.Flags().StringVarP(&project, "project", "p", "", "プロジェクトキー")
	cmd.Flags().StringVarP(&keyword, "keyword", "k", "", "キーワード")
	return cmd
}

func newViewCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "view <wikiId>",
		Short: "Wiki ページの本文を表示する",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := strconv.Atoi(args[0])
			if err != nil {
				return fmt.Errorf("Wiki ID は数値で指定してください: %s", args[0])
			}
			_, client, err := cmdutil.LoadConfigAndClient()
			if err != nil {
				return err
			}
			w, err := client.GetWiki(id)
			if err != nil {
				return err
			}
			fmt.Println(titleStyle.Render(cmdutil.StripControl(w.Name)))
			if len(w.Tags) > 0 {
				var names []string
				for _, t := range w.Tags {
					names = append(names, cmdutil.StripControl(t.Name))
				}
				fmt.Println(labelStyle.Render("タグ: ") + strings.Join(names, ", "))
			}
			fmt.Println(labelStyle.Render("更新: ") + w.Updated)
			if len(w.Attachments) > 0 {
				fmt.Println(labelStyle.Render("添付ファイル:"))
				for _, a := range w.Attachments {
					fmt.Printf("  %d\t%s\t%s\n", a.ID, cmdutil.StripControl(a.Name), api.HumanSize(a.Size))
				}
			}
			fmt.Println()
			fmt.Println(cmdutil.StripControl(w.Content))
			return nil
		},
	}
}

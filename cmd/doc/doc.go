package doc

import (
	"fmt"
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

// NewDocCmd returns the doc subcommand group.
func NewDocCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "doc",
		Short: "ドキュメントの参照",
	}
	cmd.AddCommand(newTreeCmd())
	cmd.AddCommand(newListCmd())
	cmd.AddCommand(newViewCmd())
	cmd.AddCommand(newAttachmentCmd())
	return cmd
}

func newTreeCmd() *cobra.Command {
	var project string
	cmd := &cobra.Command{
		Use:   "tree",
		Short: "ドキュメントツリーを表示する",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, client, err := cmdutil.LoadConfigAndClient()
			if err != nil {
				return err
			}
			key, err := cmdutil.ResolveProjectKey(cfg, project)
			if err != nil {
				return err
			}
			tree, err := client.GetDocumentTree(key)
			if err != nil {
				return err
			}
			if tree.ActiveTree == nil || len(tree.ActiveTree.Children) == 0 {
				fmt.Println("ドキュメントはありません")
				return nil
			}
			printTree(tree.ActiveTree.Children, 0)
			return nil
		},
	}
	cmd.Flags().StringVarP(&project, "project", "p", "", "プロジェクトキー")
	return cmd
}

func printTree(nodes []api.DocumentNode, depth int) {
	for _, n := range nodes {
		name := cmdutil.StripControl(n.Name)
		if name == "" {
			name = "(無題)"
		}
		if n.Emoji != "" {
			name = n.Emoji + " " + name
		}
		fmt.Printf("%s%s  %s\n", strings.Repeat("  ", depth), name, labelStyle.Render(n.ID))
		printTree(n.Children, depth+1)
	}
}

func newListCmd() *cobra.Command {
	var (
		project string
		keyword string
		count   int
	)
	cmd := &cobra.Command{
		Use:   "list",
		Short: "ドキュメント一覧を表示する",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, client, err := cmdutil.LoadConfigAndClient()
			if err != nil {
				return err
			}
			key, err := cmdutil.ResolveProjectKey(cfg, project)
			if err != nil {
				return err
			}
			proj, err := client.GetProject(key)
			if err != nil {
				return err
			}
			docs, err := client.GetDocuments(proj.ID, keyword, count)
			if err != nil {
				return err
			}
			if len(docs) == 0 {
				fmt.Println("該当するドキュメントはありません")
				return nil
			}
			fmt.Printf("%s\t%s\t%s\n", headerStyle.Render("ID"), headerStyle.Render("UPDATED"), headerStyle.Render("TITLE"))
			for _, d := range docs {
				fmt.Printf("%s\t%s\t%s\n", d.ID, d.Updated, cmdutil.StripControl(d.Title))
			}
			return nil
		},
	}
	cmd.Flags().StringVarP(&project, "project", "p", "", "プロジェクトキー")
	cmd.Flags().StringVarP(&keyword, "keyword", "k", "", "キーワード")
	cmd.Flags().IntVarP(&count, "count", "c", 20, "表示件数（最大 100）")
	return cmd
}

func newViewCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "view <documentId>",
		Short: "ドキュメントの本文を表示する",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			_, client, err := cmdutil.LoadConfigAndClient()
			if err != nil {
				return err
			}
			d, err := client.GetDocument(args[0])
			if err != nil {
				return err
			}
			fmt.Println(titleStyle.Render(cmdutil.StripControl(d.Title)))
			if len(d.Tags) > 0 {
				var names []string
				for _, t := range d.Tags {
					names = append(names, cmdutil.StripControl(t.Name))
				}
				fmt.Println(labelStyle.Render("タグ: ") + strings.Join(names, ", "))
			}
			fmt.Println(labelStyle.Render("更新: ") + d.Updated)
			if len(d.Attachments) > 0 {
				fmt.Println(labelStyle.Render("添付ファイル:"))
				for _, a := range d.Attachments {
					fmt.Printf("  %d\t%s\t%s\n", a.ID, cmdutil.StripControl(a.Name), api.HumanSize(a.Size))
				}
			}
			fmt.Println()
			fmt.Println(cmdutil.StripControl(d.Plain))
			return nil
		},
	}
}

func newAttachmentCmd() *cobra.Command {
	var (
		download bool
		id       int
		output   string
	)
	cmd := &cobra.Command{
		Use:   "attachment <documentId>",
		Short: "ドキュメントの添付ファイルを一覧・ダウンロードする",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			_, client, err := cmdutil.LoadConfigAndClient()
			if err != nil {
				return err
			}
			d, err := client.GetDocument(args[0])
			if err != nil {
				return err
			}
			if len(d.Attachments) == 0 {
				fmt.Println("添付ファイルはありません")
				return nil
			}
			if !download {
				printAttachments(d.Attachments)
				return nil
			}
			return cmdutil.DownloadAttachments(d.Attachments, id, output, func(aid int) (string, []byte, error) {
				return client.DownloadDocumentAttachment(d.ID, aid)
			})
		},
	}
	cmd.Flags().BoolVarP(&download, "download", "d", false, "ダウンロードする（--id 未指定なら全件）")
	cmd.Flags().IntVar(&id, "id", 0, "ダウンロードする添付ファイルの ID")
	cmd.Flags().StringVarP(&output, "output", "o", ".", "保存先フォルダ")
	return cmd
}

func printAttachments(attachments []api.Attachment) {
	fmt.Printf("%s\t%s\t%s\n", headerStyle.Render("ID"), headerStyle.Render("SIZE"), headerStyle.Render("NAME"))
	for _, a := range attachments {
		fmt.Printf("%d\t%s\t%s\n", a.ID, api.HumanSize(a.Size), cmdutil.StripControl(a.Name))
	}
}

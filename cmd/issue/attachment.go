package issue

import (
	"fmt"

	"github.com/KimMaru10/bl-cli/internal/api"
	"github.com/KimMaru10/bl-cli/internal/cmdutil"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

func newAttachmentCmd() *cobra.Command {
	var (
		download bool
		id       int
		output   string
	)

	cmd := &cobra.Command{
		Use:   "attachment [issueKey]",
		Short: "課題の添付ファイルを一覧・ダウンロードする",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			_, client, err := cmdutil.LoadConfigAndClient()
			if err != nil {
				return err
			}
			issueKey, err := resolveIssueKey(args)
			if err != nil {
				return err
			}

			attachments, err := client.GetIssueAttachments(issueKey)
			if err != nil {
				return err
			}
			if len(attachments) == 0 {
				fmt.Println("添付ファイルはありません")
				return nil
			}

			if !download {
				headerStyle := lipgloss.NewStyle().Bold(true)
				fmt.Printf("%s\t%s\t%s\t%s\n", headerStyle.Render("ID"), headerStyle.Render("SIZE"), headerStyle.Render("CREATED"), headerStyle.Render("NAME"))
				for _, a := range attachments {
					fmt.Printf("%d\t%s\t%s\t%s\n", a.ID, api.HumanSize(a.Size), a.Created, cmdutil.StripControl(a.Name))
				}
				return nil
			}

			return cmdutil.DownloadAttachments(attachments, id, output, func(aid int) (string, []byte, error) {
				return client.DownloadIssueAttachment(issueKey, aid)
			})
		},
	}

	cmd.Flags().BoolVarP(&download, "download", "d", false, "ダウンロードする（--id 未指定なら全件）")
	cmd.Flags().IntVar(&id, "id", 0, "ダウンロードする添付ファイルの ID")
	cmd.Flags().StringVarP(&output, "output", "o", ".", "保存先フォルダ")

	return cmd
}

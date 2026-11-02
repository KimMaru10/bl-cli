package file

import (
	"fmt"
	"strconv"

	"github.com/KimMaru10/bl-cli/internal/api"
	"github.com/KimMaru10/bl-cli/internal/cmdutil"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

var headerStyle = lipgloss.NewStyle().Bold(true)

// NewFileCmd returns the file subcommand group (shared files).
func NewFileCmd() *cobra.Command {
	var project string
	cmd := &cobra.Command{
		Use:   "file",
		Short: "共有ファイルの参照",
	}
	cmd.PersistentFlags().StringVarP(&project, "project", "p", "", "プロジェクトキー")

	list := &cobra.Command{
		Use:   "list [dir]",
		Short: "共有ファイル一覧を表示する（dir 省略時はルート）",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, client, err := cmdutil.LoadConfigAndClient()
			if err != nil {
				return err
			}
			key, err := cmdutil.ResolveProjectKey(cfg, project)
			if err != nil {
				return err
			}
			dir := ""
			if len(args) > 0 {
				dir = args[0]
			}
			files, err := client.GetSharedFiles(key, dir)
			if err != nil {
				return err
			}
			if len(files) == 0 {
				fmt.Println("共有ファイルはありません")
				return nil
			}
			fmt.Printf("%s\t%s\t%s\t%s\n", headerStyle.Render("ID"), headerStyle.Render("SIZE"), headerStyle.Render("UPDATED"), headerStyle.Render("NAME"))
			for _, f := range files {
				size := api.HumanSize(f.Size)
				name := cmdutil.StripControl(f.Dir + f.Name)
				if f.Type == "directory" {
					size = "-"
					name += "/"
				}
				fmt.Printf("%d\t%s\t%s\t%s\n", f.ID, size, f.Updated, name)
			}
			return nil
		},
	}

	var output string
	download := &cobra.Command{
		Use:   "download <fileId>",
		Short: "共有ファイルをダウンロードする",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := strconv.Atoi(args[0])
			if err != nil {
				return fmt.Errorf("ファイル ID は数値で指定してください: %s", args[0])
			}
			cfg, client, err := cmdutil.LoadConfigAndClient()
			if err != nil {
				return err
			}
			key, err := cmdutil.ResolveProjectKey(cfg, project)
			if err != nil {
				return err
			}
			name, data, err := client.DownloadSharedFile(key, id)
			if err != nil {
				return err
			}
			path, err := cmdutil.SaveFile(output, name, fmt.Sprintf("file-%d", id), data)
			if err != nil {
				return err
			}
			fmt.Println(path)
			return nil
		},
	}
	download.Flags().StringVarP(&output, "output", "o", ".", "保存先フォルダ")

	cmd.AddCommand(list)
	cmd.AddCommand(download)
	return cmd
}

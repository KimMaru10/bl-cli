package issue

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/KimMaru10/bl-cli/internal/cmdutil"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

var (
	commentHeaderStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("8"))
	separatorStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("8"))
)

func newCommentCmd() *cobra.Command {
	var body string

	cmd := &cobra.Command{
		Use:   "comment [issueKey]",
		Short: "課題にコメントを追加する",
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

			var content string

			if body != "" {
				// Inline mode
				content = body
			} else {
				// Editor mode
				editor := os.Getenv("EDITOR")
				if editor == "" {
					editor = "vim"
				}

				tmpDir := os.TempDir()
				tmpFile := filepath.Join(tmpDir, "bl-comment-"+issueKey+".md")

				template := "# 1行目以降にコメントを入力してください。# で始まる行は無視されます\n"
				if err := os.WriteFile(tmpFile, []byte(template), 0644); err != nil {
					return fmt.Errorf("一時ファイルの作成に失敗しました: %w", err)
				}
				defer os.Remove(tmpFile)

				editorCmd := exec.Command(editor, tmpFile)
				editorCmd.Stdin = os.Stdin
				editorCmd.Stdout = os.Stdout
				editorCmd.Stderr = os.Stderr
				if err := editorCmd.Run(); err != nil {
					return fmt.Errorf("エディタの起動に失敗しました: %w", err)
				}

				file, err := os.Open(tmpFile)
				if err != nil {
					return fmt.Errorf("一時ファイルの読み込みに失敗しました: %w", err)
				}
				defer file.Close()

				var lines []string
				scanner := bufio.NewScanner(file)
				for scanner.Scan() {
					line := scanner.Text()
					if !strings.HasPrefix(line, "#") {
						lines = append(lines, line)
					}
				}
				content = strings.TrimSpace(strings.Join(lines, "\n"))
			}

			if content == "" {
				fmt.Println("コメントが空のため中止しました")
				return nil
			}

			_, err = client.AddComment(issueKey, content)
			if err != nil {
				return err
			}

			fmt.Println(successStyle.Render("✔ " + issueKey + " にコメントを追加しました"))
			return nil
		},
	}

	cmd.Flags().StringVarP(&body, "body", "b", "", "コメント本文")

	// Add list subcommand
	cmd.AddCommand(newCommentListCmd())

	return cmd
}

func newCommentListCmd() *cobra.Command {
	var count int

	cmd := &cobra.Command{
		Use:   "list [issueKey]",
		Short: "コメント一覧を表示する",
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

			comments, err := client.GetComments(issueKey, count, "desc")
			if err != nil {
				return err
			}

			if len(comments) == 0 {
				fmt.Println("コメントはありません")
				return nil
			}

			changeStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("6"))
			labelStyle := lipgloss.NewStyle().Bold(true)

			for i, c := range comments {
				userName := ""
				if c.CreatedUser != nil {
					userName = c.CreatedUser.Name
				}
				header := commentHeaderStyle.Render(userName + " — " + c.Created)
				fmt.Println(header)

				// Show change logs
				for _, cl := range c.ChangeLog {
					if cl.OriginalValue != "" && cl.NewValue != "" {
						fmt.Println(changeStyle.Render(fmt.Sprintf("  %s：%s → %s", cl.Field, cl.OriginalValue, cl.NewValue)))
					} else if cl.NewValue != "" {
						fmt.Println(changeStyle.Render(fmt.Sprintf("  %s：→ %s", cl.Field, cl.NewValue)))
					} else if cl.OriginalValue != "" {
						fmt.Println(changeStyle.Render(fmt.Sprintf("  %s：%s →", cl.Field, cl.OriginalValue)))
					}
				}

				if c.Content != "" {
					fmt.Println(labelStyle.Render("comment") + "：" + c.Content)
				}

				if i < len(comments)-1 {
					fmt.Println(separatorStyle.Render("---"))
				}
			}
			return nil
		},
	}

	cmd.Flags().IntVarP(&count, "count", "c", 10, "表示件数")

	return cmd
}

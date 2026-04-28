package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/godjian/myppt-app/internal/styles"
)

// newStylesCmd 风格列表命令
func newStylesCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "styles",
		Short: "列出所有可用风格",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println()
			fmt.Println("可用风格列表")
			fmt.Println()

			styleList := styles.ListStyles()

			// 按类别分组
			categories := make(map[string][]styles.Style)
			for _, s := range styleList {
				cat := s.Category
				if cat == "" {
					cat = "其他"
				}
				categories[cat] = append(categories[cat], s)
			}

			for cat, list := range categories {
				fmt.Printf("%s\n", cat)
				fmt.Printf("--------------------------------------------------\n")
				for _, s := range list {
					fmt.Printf("  %s (%s)\n", s.ID, s.Name)
					fmt.Printf("     %s\n", s.Description)
					colors := ""
					for i, c := range s.Palette {
						if i >= 4 {
							colors += " ..."
							break
						}
						colors += " " + c
					}
					fmt.Printf("     配色: %s\n", colors)
				}
				fmt.Println()
			}

			fmt.Println("使用示例:")
			fmt.Println("  oh-my-ppt generate --topic 'AI发展趋势' --style minimal-white")
			fmt.Println()
		},
	}
}

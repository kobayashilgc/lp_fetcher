package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"lp_fetcher_golang/internal/fetcher"
)

func init() {
	cmd := &cobra.Command{
		Use:   "category",
		Short: "获取店铺分类叶子列表并写入 category_result.json",
		RunE:  runCategory,
	}
	cmd.Flags().String("output", "category_result.json", "输出 JSON 文件路径")
	rootCmd.AddCommand(cmd)
}

func runCategory(cmd *cobra.Command, args []string) error {
	output, err := cmd.Flags().GetString("output")
	if err != nil {
		return err
	}

	result := fetcher.FetchCategories(wdtoken)
	if err := finishFetchCommand(output, result.Status, result.Msg, result); err != nil {
		return err
	}
	fmt.Printf("分类获取完成，共 %d 个叶子分类，已写入 %s\n", result.Count, output)
	return nil
}

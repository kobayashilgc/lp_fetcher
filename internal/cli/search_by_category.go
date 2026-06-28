package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"lp_fetcher_golang/internal/fetcher"
)

func init() {
	cmd := &cobra.Command{
		Use:   "search_by_category",
		Short: "根据分类 ID 检索商品并写入 search_result.json",
		RunE:  runSearchByCategory,
	}
	cmd.Flags().String("cateId", "", "店铺分类 ID")
	cmd.Flags().String("output", "search_result.json", "输出 JSON 文件路径")
	cmd.MarkFlagRequired("cateId")
	rootCmd.AddCommand(cmd)
}

func runSearchByCategory(cmd *cobra.Command, args []string) error {
	cateID, err := cmd.Flags().GetString("cateId")
	if err != nil {
		return err
	}
	if cateID == "" {
		return fmt.Errorf("cateId 为必填参数")
	}

	output, err := cmd.Flags().GetString("output")
	if err != nil {
		return err
	}

	result := fetcher.FetchItemsByCategory(cateID, wdtoken, dash)
	if err := finishFetchCommand(output, result.Status, result.Msg, result); err != nil {
		return err
	}
	fmt.Printf("分类检索完成，共 %d 条结果，已写入 %s\n", result.Count, output)
	return nil
}

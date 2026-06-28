package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"lp_fetcher_golang/internal/fetcher"
)

func init() {
	cmd := &cobra.Command{
		Use:   "search",
		Short: "根据关键字检索商品并写入 search_result.json",
		RunE:  runSearch,
	}
	cmd.Flags().String("keyword", "", "店铺内搜索词")
	cmd.Flags().String("output", "search_result.json", "输出 JSON 文件路径")
	cmd.MarkFlagRequired("keyword")
	rootCmd.AddCommand(cmd)
}

func runSearch(cmd *cobra.Command, args []string) error {
	keyword, err := cmd.Flags().GetString("keyword")
	if err != nil {
		return err
	}
	if keyword == "" {
		return fmt.Errorf("keyword 为必填参数")
	}

	output, err := cmd.Flags().GetString("output")
	if err != nil {
		return err
	}

	result := fetcher.FetchItems(keyword, wdtoken, dash)
	if err := finishFetchCommand(output, result.Status, result.Msg, result); err != nil {
		return err
	}
	fmt.Printf("检索完成，共 %d 条结果，已写入 %s\n", result.Count, output)
	return nil
}

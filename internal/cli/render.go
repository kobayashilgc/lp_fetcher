package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"lp_fetcher_golang/internal/render"
)

func init() {
	cmd := &cobra.Command{
		Use:   "render",
		Short: "读取 search_result.json 生成 report.html",
		RunE:  runRender,
	}
	cmd.Flags().String("input", "search_result.json", "输入 JSON 文件路径")
	cmd.Flags().String("template", "report_template.html", "HTML 模板路径")
	cmd.Flags().String("output", "report.html", "输出 HTML 文件路径")
	rootCmd.AddCommand(cmd)
}

func runRender(cmd *cobra.Command, args []string) error {
	input, err := cmd.Flags().GetString("input")
	if err != nil {
		return err
	}
	template, err := cmd.Flags().GetString("template")
	if err != nil {
		return err
	}
	output, err := cmd.Flags().GetString("output")
	if err != nil {
		return err
	}

	result, keyword, err := render.LoadSearchResult(input)
	if err != nil {
		return err
	}
	if err := render.RenderReport(result, keyword, template, output); err != nil {
		return err
	}
	fmt.Printf("报告已生成: %s\n", output)
	return nil
}

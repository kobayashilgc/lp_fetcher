package cli

import (
	"os"

	"github.com/spf13/cobra"
)

var (
	wdtoken string
	dash    string
)

var rootCmd = &cobra.Command{
	Use:   "lp_fetcher",
	Short: "浣熊唱片商品检索",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&wdtoken, "wdtoken", "", "API 鉴权 token")
	rootCmd.PersistentFlags().StringVar(&dash, "dash", "", "URL 防缓存时间戳")
}

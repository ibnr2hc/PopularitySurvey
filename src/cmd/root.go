package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "PopularitySurvey",
	Short: "指定したユーザーのフォロワーをランキング付けする",
	Long: `指定したユーザーのフォロワーをランキング付けする。
フォロワー数をソートし、上位から数十名を表示する。`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

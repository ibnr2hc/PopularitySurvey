package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var surveyCmd = &cobra.Command{
	Use:   "survey",
	Short: "指定したユーザーのフォロワーを調査しランキングにする。",
	Long: `指定したユーザーのフォロワーを調査しランキングにする。
フォロワー数が多い順にソートして表示をする。`,
	Run: func(cmd *cobra.Command, args []string) {
		screenName, err := cmd.Flags().GetString("screen_name")
		if err != nil {
			fmt.Println(err)
		}

		fmt.Println("調査対象：ほにゃらら(@" + screenName + ") フォロワー数100人")
		fmt.Println("人気フォロワー1位：あれこれ(@xxx) フォロワー数1000人")
		fmt.Println("人気フォロワー2位：だれそれ(@xxx) フォロワー数900人")
		fmt.Println("人気フォロワー3位：それそれ(@xxx) フォロワー数800人")
	},
}

func init() {
	rootCmd.AddCommand(surveyCmd)

	// Option: screen_name
	surveyCmd.Flags().StringP("screen_name", "s", "", "Twitter user screen name")
}

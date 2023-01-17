package cmd

import (
	"github.com/ibnr2hc/PopularitySurvey/cmd/models"
	"github.com/spf13/cobra"
	"log"
)

var surveyCmd = &cobra.Command{
	Use:   "survey",
	Short: "指定したユーザーのフォロワーを調査しランキングにする。",
	Long: `指定したユーザーのフォロワーを調査しランキングにする。
フォロワー数が多い順にソートして表示をする。`,
	Run: func(cmd *cobra.Command, args []string) {
		screenName, err := cmd.Flags().GetString("screen_name")
		if err != nil || screenName == "" {
			log.Fatal("[Error] 検索する対象のScreenNameは必須です。-sオプションの後にScreenNameを指定してください。")
		}

		// 調査対象のユーザー情報の取得
		surveyTarget := models.NewUserForSurveyTarget(screenName)
		// ランキングされたフォロワーを表示する
		followers := surveyTarget.GetFollowerForTopRanking()
		models.ShowFollowers(followers)
	},
}

func init() {
	rootCmd.AddCommand(surveyCmd)

	// Option: screen_name
	surveyCmd.Flags().StringP("screen_name", "s", "", "Twitter user screen name")
}

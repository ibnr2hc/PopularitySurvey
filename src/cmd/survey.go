package cmd

import (
	"context"
	"fmt"
	"github.com/olekukonko/tablewriter"
	"github.com/sivchari/gotwtr"
	"github.com/spf13/cobra"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

var surveyCmd = &cobra.Command{
	Use:   "survey",
	Short: "指定したユーザーのフォロワーを調査しランキングにする。",
	Long: `指定したユーザーのフォロワーを調査しランキングにする。
フォロワー数が多い順にソートして表示をする。`,
	Run: func(cmd *cobra.Command, args []string) {
		surveyTargetUser, err := cmd.Flags().GetString("screen_name")
		if err != nil || surveyTargetUser == "" {
			log.Fatal("[Error] 検索する対象のScreenNameは必須です。-sオプションの後にScreenNameを指定してください。")
		}

		// フォロワーとそのユーザーのフォロワー数を取得し上位n人を取得する。
		// TODO: 上位n名を指定できるようにする。
		TOKEN := os.Getenv("TWITTER_BEARER_TOKEN")
		if TOKEN == "" {
			log.Fatal("[Error] Twitter API Bearer Tokenの環境変数(TWITTER_BEARER_TOKEN)を設定してください。")
		}
		client := gotwtr.New(TOKEN)

		// 調査対象のユーザー情報を表示する
		opt := gotwtr.RetrieveUserOption{
			UserFields: []gotwtr.UserField{
				gotwtr.UserFieldName,
				gotwtr.UserFieldUserName,
				gotwtr.UserFieldID,
			},
		}
		surveyTargetUserInfo, err := client.RetrieveSingleUserWithUserName(context.Background(), surveyTargetUser, &opt)
		if err != nil {
			if strings.Contains(err.Error(), "429") {
				log.Fatal("[Error] レート制限のためしばらくしてから実行してください。")
			}
			log.Fatal("[Error] 存在しないScreenNameか、もしくは何かしらのエラーが発生しました。")
		}
		fmt.Println("[Info] 調査対象：" + surveyTargetUserInfo.User.Name + "(@" + surveyTargetUserInfo.User.UserName + ")")

		// 調査対象のフォロワー情報を取得する
		followerOpt := gotwtr.FollowOption{
			UserFields: []gotwtr.UserField{
				gotwtr.UserFieldName,
				gotwtr.UserFieldUserName,
				gotwtr.UserFieldID,
				gotwtr.UserFieldPublicMetrics,
			},
			MaxResults: 1000,
		}
		nextToken := ""
		followerInfo := map[int]map[string]string{}
		// TODO: フォロワー数でソートし、フォロワー数をキーにしてランキング結果を表示すると同一のフォロワー数をもつユーザー1人分が複数出てしまう。対応する。
		followerIndex := []int{}
		for {
			followers, err := client.Followers(context.Background(), surveyTargetUserInfo.User.ID, &followerOpt)
			if err != nil {
				// レート制限を超えている場合はしばらく待機する
				if strings.Contains(err.Error(), "429") {
					fmt.Println("[Debug] レート制限のため待機し再度実行します...")
					time.Sleep(time.Second * 310) // 5分10秒。10秒はAPIをリクエストしている時間分のバッファを考慮した。
					continue
				}

				// 長時間のHTTP Connectionにより切断された場合は再度実行する
				if strings.Contains(err.Error(), "read: connection reset by peer") {
					fmt.Println("[Debug] HTTP Connectionが切断されたため再度実行します。")
					time.Sleep(time.Second * 60)
					continue
				}

				// レート制限以外のエラーの場合は処理を終了する。
				log.Fatal("[Error] Something Error: " + err.Error())
			}
			for _, v := range followers.Users {
				followerInfo[v.PublicMetrics.FollowersCount] = map[string]string{
					"screenName":  v.UserName,
					"displayName": v.Name,
				}
				followerIndex = append(followerIndex, v.PublicMetrics.FollowersCount)
			}

			fmt.Println("[Debug] " + strconv.Itoa(len(followerIndex)) + "のユーザー情報を取得しました。")
			nextToken = followers.Meta.NextToken
			if nextToken == "" {
				fmt.Println("[Debug] 最終のPaginationTokenまで達したためフォロワー取得処理を終了します。")
				break
			}

			followerOpt.PaginationToken = nextToken
			fmt.Println("[Debug]　 → NextToken: " + nextToken)
		}
		sort.Sort(sort.Reverse(sort.IntSlice(followerIndex)))
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Rank", "DisplayName", "ScreenName", "Followers Count"})
		for i, v := range followerIndex {
			table.Append([]string{strconv.Itoa(i + 1), followerInfo[v]["displayName"], "@" + followerInfo[v]["screenName"], strconv.Itoa(v)})

			// 上位20名のみ表示する。
			const RANKING_LIMIT = 20
			if i >= RANKING_LIMIT-1 {
				break
			}
		}
		fmt.Println("[Info] フォロワー数ランキングの計算が終了しました。")
		table.Render()
	},
}

func init() {
	rootCmd.AddCommand(surveyCmd)

	// Option: screen_name
	surveyCmd.Flags().StringP("screen_name", "s", "", "Twitter user screen name")
}

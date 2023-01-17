package repositories

import (
	"context"
	"fmt"
	"github.com/sivchari/gotwtr"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

type Twitter struct {
	client *gotwtr.Client
}

// NewTwitter
// TODO: doc
func NewTwitter() *Twitter {
	TOKEN := os.Getenv("TWITTER_BEARER_TOKEN")
	if TOKEN == "" {
		log.Fatal("[Error] Twitter API Bearer Tokenの環境変数(TWITTER_BEARER_TOKEN)を設定してください。")
	}
	return &Twitter{
		client: gotwtr.New(TOKEN),
	}
}

// GetFollowers
// TODO: doc
// TODO: refactoring
func (t *Twitter) GetFollowers(userId string) []map[string]string {
	opt := gotwtr.FollowOption{
		UserFields: []gotwtr.UserField{
			gotwtr.UserFieldName,
			gotwtr.UserFieldUserName,
			gotwtr.UserFieldID,
			gotwtr.UserFieldPublicMetrics,
		},
		MaxResults: 1000,
	}

	nextToken := ""
	followerInfo := []map[string]string{}
	for { // フォロワー情報の全てを取得するまで処理を行う
		followers, err := t.client.Followers(context.Background(), userId, &opt)
		if err != nil {
			// レート制限を超えている場合はしばらく待機する
			if strings.Contains(err.Error(), "429") {
				fmt.Println("[Debug] レート制限のため待機し再度実行します...")
				// TODO: config化する
				time.Sleep(time.Second * 120) // レート制限による再度実行のため待機
				continue
			}
			// 長時間のHTTP Connectionにより切断された場合は再度実行する
			if strings.Contains(err.Error(), "read: connection reset by peer") {
				fmt.Println("[Debug] HTTP Connectionが切断されたため再度実行します。")
				// TODO: config化する
				time.Sleep(time.Second * 60) // コネクション切断による再接続のための待機
				continue
			}
			// レート制限以外のエラーの場合は処理を終了する
			log.Fatal("[Error] Something Error: " + err.Error())
		}

		// 取得したフォロワー情報を整形する
		for _, v := range followers.Users {
			followerInfo = append(followerInfo, map[string]string{
				"screenName":    v.UserName,
				"displayName":   v.Name,
				"userId":        v.ID,
				"followerCount": strconv.Itoa(v.PublicMetrics.FollowersCount),
			})
		}

		// TODO: 取得したユーザー数を表示する
		fmt.Println("[Debug] ユーザー情報を取得しました。")
		nextToken = followers.Meta.NextToken
		if nextToken == "" { // フォロワー情報を全て取得した後にフォロワー取得処理を終える。
			fmt.Println("[Debug] 最終のPaginationTokenまで達したためフォロワー取得処理を終了します。")
			break
		}

		opt.PaginationToken = nextToken
		fmt.Println("[Debug]　 → NextToken: " + nextToken)
	}
	return followerInfo

}

//		nextToken := ""
//		followerInfo := map[int]map[string]string{}
//		// TODO: フォロワー数でソートし、フォロワー数をキーにしてランキング結果を表示すると同一のフォロワー数をもつユーザー1人分が複数出てしまう。対応する。
//		followerIndex := []int{}
//		for {
//			followers, err := client.Followers(context.Background(), surveyTargetUserInfo.User.ID, &followerOpt)
//			if err != nil {
//				// レート制限を超えている場合はしばらく待機する
//				if strings.Contains(err.Error(), "429") {
//					fmt.Println("[Debug] レート制限のため待機し再度実行します...")
//					time.Sleep(time.Second * 310) // 5分10秒。10秒はAPIをリクエストしている時間分のバッファを考慮した。
//					continue
//				}
//
//				// 長時間のHTTP Connectionにより切断された場合は再度実行する
//				if strings.Contains(err.Error(), "read: connection reset by peer") {
//					fmt.Println("[Debug] HTTP Connectionが切断されたため再度実行します。")
//					time.Sleep(time.Second * 60)
//					continue
//				}
//
//				// レート制限以外のエラーの場合は処理を終了する。
//				log.Fatal("[Error] Something Error: " + err.Error())
//			}
//			for _, v := range followers.Users {
//				followerInfo[v.PublicMetrics.FollowersCount] = map[string]string{
//					"screenName":  v.UserName,
//					"displayName": v.Name,
//				}
//				followerIndex = append(followerIndex, v.PublicMetrics.FollowersCount)
//			}
//
//			fmt.Println("[Debug] " + strconv.Itoa(len(followerIndex)) + "のユーザー情報を取得しました。")
//			nextToken = followers.Meta.NextToken
//			if nextToken == "" {
//				fmt.Println("[Debug] 最終のPaginationTokenまで達したためフォロワー取得処理を終了します。")
//				break
//			}
//
//			followerOpt.PaginationToken = nextToken
//			fmt.Println("[Debug]　 → NextToken: " + nextToken)
//		}

// GetUserInfo ユーザー情報を取得する。
//
// Args:
//   - screenName string: 情報を取得するユーザーのスクリーンネーム(@は不要)
//
// Returns:
//   - user map[string]string: ユーザー情報
//
// Errors:
//   - レート制限がかかっている場合
//   - ユーザーが存在しないなどの場合
func (t *Twitter) GetUserInfo(screenName string) map[string]string {
	// ユーザー情報を取得する
	opt := gotwtr.RetrieveUserOption{
		UserFields: []gotwtr.UserField{
			gotwtr.UserFieldID,
			gotwtr.UserFieldName,
			gotwtr.UserFieldUserName,
			gotwtr.UserFieldPublicMetrics,
		},
	}
	user, err := t.client.RetrieveSingleUserWithUserName(context.Background(), screenName, &opt)
	if err != nil {
		if strings.Contains(err.Error(), "429") { // TWITTER APIのレートに制限されている場合
			log.Fatal("[Error] レート制限のためしばらくしてから実行してください。")
		}
		// 他の何かしらのエラーの場合
		// TODO: 存在しないユーザーはこの分岐を通らない
		log.Fatal("[Error] 存在しないScreenNameか、もしくは何かしらのエラーが発生しました。")
	}

	return map[string]string{
		"ID":            user.User.ID,
		"screenName":    screenName,
		"displayName":   user.User.Name,
		"followerCount": strconv.Itoa(user.User.PublicMetrics.FollowersCount),
	}
}

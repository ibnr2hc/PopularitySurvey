package repositories

import (
	"context"
	"fmt"
	"github.com/ibnr2hc/PopularitySurvey/cmd/config"
	"github.com/sivchari/gotwtr"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

// Twitter Twitter APIを内包するstruct
type Twitter struct {
	client *gotwtr.Client
}

// ITwitter interface
type ITwitter interface {
	getFollower(userId string, nextToken string) (*gotwtr.FollowersResponse, bool)
	GetAllFollowers(userId string) []map[string]string
	GetUserInfo(screenName string) map[string]string
}

var _ ITwitter = (*Twitter)(nil)

// NewTwitter Twitter APIに関連する機能を持つTwitter structを返す
//
// Returns:
//   - twitter Twitter: Twitter Repositoryを持つstruct
func NewTwitter() *Twitter {
	TOKEN := os.Getenv("TWITTER_BEARER_TOKEN")
	if TOKEN == "" {
		log.Fatal("[Error] Twitter API Bearer Tokenの環境変数(TWITTER_BEARER_TOKEN)を設定してください。")
	}
	return &Twitter{
		client: gotwtr.New(TOKEN),
	}
}

// フォロワーを取得する
//
// Args:
//   - userId string   : フォロワーを取得する対象のユーザーID
//   - nextToken string: TWITTER APIのページネーショントークン
//
// Returns:
//   - followers gotwtr.FollowersResponse: フォロワー情報
//   - isRetry bool                      : リトライが必要か
func (t *Twitter) getFollower(userId string, nextToken string) (*gotwtr.FollowersResponse, bool) {
	opt := gotwtr.FollowOption{
		UserFields: []gotwtr.UserField{
			gotwtr.UserFieldName,
			gotwtr.UserFieldUserName,
			gotwtr.UserFieldID,
			gotwtr.UserFieldPublicMetrics,
		},
		MaxResults: 1000,
	}
	if nextToken != "" { // ページネーションのトークンがある場合
		opt.PaginationToken = nextToken
	}

	followers, err := t.client.Followers(context.Background(), userId, &opt)
	if err != nil {
		// レート制限を超えている場合はしばらく待機する
		if strings.Contains(err.Error(), "429") {
			fmt.Println("[Debug] レート制限のため待機し再度実行します...")
			time.Sleep(config.RATE_LIMIT_WAITING_SECOND) // レート制限による再度実行のため待機
			return &gotwtr.FollowersResponse{}, true
		}
		// 長時間のHTTP Connectionにより切断された場合は再度実行する
		if strings.Contains(err.Error(), "read: connection reset by peer") {
			fmt.Println("[Debug] HTTP Connectionが切断されたため再度実行します。")
			time.Sleep(config.RECONNECT_HTTP_WAITING_SECOND) // コネクション切断による再接続のための待機
			return &gotwtr.FollowersResponse{}, true
		}
		// レート制限以外のエラーの場合は処理を終了する
		log.Fatal("[Error] Something Error: " + err.Error())
	}
	return followers, false
}

// GetAllFollowers 指定したuserIdのフォロワー情報を返却する
//
// Args:
//   - userId string: 調査するユーザーのID
//
// Returns:
//   - followers []map[string]string: フォロワーの情報
func (t *Twitter) GetAllFollowers(userId string) []map[string]string {
	nextToken := ""
	followerInfo := []map[string]string{}
	for { // フォロワー情報の全てを取得するまで処理を行う
		followers, needToRetry := t.getFollower(userId, nextToken)
		if needToRetry { // レート制限などによりリトライが必要な場合
			continue
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

		fmt.Println("[Debug]　 → NextToken: " + nextToken)
	}
	return followerInfo

}

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
		log.Fatal("[Error] 何かしらのエラーが発生しました。")
	}
	if user.User == nil { // 存在しないユーザーの場合
		log.Fatal("[Error] 存在しないScreenNameです。")
	}

	return map[string]string{
		"ID":            user.User.ID,
		"screenName":    screenName,
		"displayName":   user.User.Name,
		"followerCount": strconv.Itoa(user.User.PublicMetrics.FollowersCount),
	}
}

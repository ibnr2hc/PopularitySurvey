package models

import (
	"fmt"
	"github.com/ibnr2hc/PopularitySurvey/cmd/config"
	"github.com/ibnr2hc/PopularitySurvey/cmd/data/repositories"
	"github.com/olekukonko/tablewriter"
	"log"
	"os"
	"sort"
	"strconv"
)

// User Twitterの1ユーザーを表現する
type User struct {
	ID            string // UserID(e.g. 123456789)
	DisplayName   string // ディスプレイネーム(e.g. あいうえお)
	ScreenName    string // スクリーンネーム(e.g. aiueo123)
	FollowerCount int    // フォロワー数
	Followers     []User // フォロワー
}

// IUser interface
type IUser interface {
	GetFollowerForTopRanking() []User
	fetchFollower()
}

var _ IUser = (*User)(nil)

// NewUserForSurveyTarget 調査対象のユーザーのUser structを返す。
//
// Args:
//   - ScreenName string: スクリーンネーム(@は不要)
//
// Returns:
//   - user User: Twitterユーザー
//
// Errors:
//   - 何らかの理由によりフォロワー数が取得できない場合。
func NewUserForSurveyTarget(screenName string) *User {
	user := new(User)
	user.ScreenName = screenName

	// ユーザー情報を取得する
	twitter := repositories.NewTwitter()
	userInfo := twitter.GetUserInfo(screenName)
	followerCount, err := strconv.Atoi(userInfo["followerCount"])
	if err != nil { // 何らかの理由によりフォロワー数が取得できない場合
		log.Fatal(err)
	}
	user.FollowerCount = followerCount
	user.ID = userInfo["ID"]
	user.DisplayName = userInfo["displayName"]
	// フォロワー情報を取得する
	user.fetchFollower()

	fmt.Println("[Info] 調査対象：" + user.DisplayName + "(@" + user.ScreenName + ")")
	return user
}

// ShowFollowers フォロワーをランキング形式で表示する。
//
// Args:
//   - Followers []User: ランキングで表示したいフォロワーの配列
func ShowFollowers(followers []User) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{
		"Rank",            // ランキング
		"Display Name",    // ディスプレイネーム
		"Screen Name",     // スクリーンネーム
		"Followers Count", // フォロワー数
	})

	for i, follower := range followers {
		table.Append([]string{
			strconv.Itoa(i + 1),
			follower.DisplayName,
			"@" + follower.ScreenName,
			strconv.Itoa(follower.FollowerCount),
		})

		// 上位n名のみ表示する。
		if i >= config.SHOW_RANKING_LIMIT-1 {
			break
		}
	}
	table.Render()
}

// GetFollowerForTopRanking 上位からランキングしたフォロワー一覧を取得する
//
// Returns:
//   - users []User: ランキングで順位を付けたフォロワー一覧
func (u *User) GetFollowerForTopRanking() []User {
	// ランキング計算のためにフォロワーのフォロワー数を一覧として取得する。
	// TODO: 処理を最適化したい
	followersCount := []int{}   // フォロワーのフォロワー数一覧(e.g. [1, 50, 100])
	followers := map[int]User{} // フォロワー数をキーとしたフォロワー(e.g. {1: User{}, 50: User{}})
	for _, follower := range u.Followers {
		followersCount = append(followersCount, follower.FollowerCount)
		followers[follower.FollowerCount] = follower
	}

	// フォロワーのフォロワー数を基準に上位からソートする。
	sort.Sort(sort.Reverse(sort.IntSlice(followersCount)))
	var rankingFollowers []User
	for _, followerCount := range followersCount {
		rankingFollowers = append(rankingFollowers, followers[followerCount])
		// TODO: 表示するランキング数をここで制御してもいいかもしれない(処理負荷の軽減)
	}
	fmt.Println("[Info] フォロワー数ランキングの計算が終了しました。")
	return rankingFollowers
}

// フォロワーを取得し、User.followersに保持する。
func (u *User) fetchFollower() {
	twitter := repositories.NewTwitter()
	followers := twitter.GetAllFollowers(u.ID)

	for _, follower := range followers {
		followerCount, _ := strconv.Atoi(follower["followerCount"])
		u.Followers = append(u.Followers, User{
			DisplayName:   follower["displayName"],
			ScreenName:    follower["screenName"],
			FollowerCount: followerCount,
		})

	}
}

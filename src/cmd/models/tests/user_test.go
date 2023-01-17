package tests

import (
	"github.com/ibnr2hc/PopularitySurvey/cmd/models"
	"github.com/stretchr/testify/assert"
	"testing"
)

// 正常系: ランキングでソートされていること
func TestUser_GetFollowerForTopRanking(t *testing.T) {
	// Preparation
	followers := []models.User{
		{
			ID:            "1",
			DisplayName:   "A",
			ScreenName:    "A",
			FollowerCount: 100,
			Followers:     []models.User{{ID: "F1"}},
		},
		{
			ID:            "2",
			DisplayName:   "B",
			ScreenName:    "B",
			FollowerCount: 300,
			Followers:     []models.User{{ID: "F2"}},
		},
		{
			ID:            "3",
			DisplayName:   "C",
			ScreenName:    "C",
			FollowerCount: 200,
			Followers:     []models.User{{ID: "F3"}},
		},
	}
	user := models.User{Followers: followers}
	rankingFollower := user.GetFollowerForTopRanking()

	// Test
	assert.Equal(t, rankingFollower[0].ID, followers[1].ID)
	assert.Equal(t, rankingFollower[0].DisplayName, followers[1].DisplayName)
	assert.Equal(t, rankingFollower[0].ScreenName, followers[1].ScreenName)
	assert.Equal(t, rankingFollower[0].FollowerCount, followers[1].FollowerCount)
	assert.Equal(t, rankingFollower[0].Followers[0].ID, followers[1].Followers[0].ID)
}

// 正常系: 処理がエラーが出ずに終了すること
func TestUser_ShowFollowers(t *testing.T) {
	defer func() {
		err := recover()
		if err != nil {
			t.Error("Unexpected processing.")
		}
	}()

	followers := []models.User{}
	models.ShowFollowers(followers)
}

// TODO: fetchFollower()
// TODO: NewUserForSurveyTarget()

package tests

import (
	"github.com/ibnr2hc/PopularitySurvey/cmd/models"
	"strconv"
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
	if rankingFollower[0].ID != followers[1].ID {
		t.Error("Not expected ID=" + rankingFollower[0].ID + ", expect " + followers[1].ID)
	}
	if rankingFollower[0].DisplayName != followers[1].DisplayName {
		t.Error("Not expected DisplayName=" + rankingFollower[0].DisplayName + ", expect " + followers[1].DisplayName)
	}
	if rankingFollower[0].ScreenName != followers[1].ScreenName {
		t.Error("Not expected ScreenName=" + rankingFollower[0].ScreenName + ", expect " + followers[1].ScreenName)
	}
	if rankingFollower[0].Followers[0].ID != followers[1].Followers[0].ID {
		t.Error("Not expected Followers=" + rankingFollower[0].Followers[0].ID + ", expect " + followers[1].Followers[0].ID)
	}
	if rankingFollower[0].FollowerCount != followers[1].FollowerCount {
		t.Error("Not expected FollowCount=" + strconv.Itoa(rankingFollower[0].FollowerCount) + ", expect " + strconv.Itoa(followers[1].FollowerCount))
	}
}

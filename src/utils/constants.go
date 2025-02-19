package utils

import "sync"

type TxtTasksData struct {
	Tweets   []string
	Comments []string
	Pictures []string
}

type TasksTargetInfo struct {
	FollowTargets            []string
	UnfollowTargets          []string
	RetweetLinks             []string
	LikeLinks                []string
	QuoteTweetLinks          []string
	CommentTweetLinks        []string
	MutualSubsFollowersCount int
	Votes                    struct {
		Links  []string
		Choice string
	}
}

// Account represents a single row from the Excel file
type Account struct {
	AuthToken string
	Proxy     string
	Username  string
	Status    string
}

// ExcelMutex ensures thread-safe access to Excel files
var ExcelMutex sync.Mutex

// MainMenuOptions is the main menu options
var MainMenuOptions = []string{
	"Follow",
	"Tweet",
	"Tweet with Picture",
	"Quote Tweet",
	"Quote Tweet with Picture",
	"Like",
	"Comment",
	"Comment with Picture",
	"Retweet",
	"Unfollow",
	"Check Valid",
	"Vote on Poll",
	"Mutual Subscription",
}

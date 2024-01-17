package additional_twitter_methods

import (
	"strings"
	"twitter/extra"
)

func GetTweetID(tweetLink string) string {
	tweetLink = strings.TrimSpace(tweetLink)

	var tweetID string
	if strings.Contains(tweetLink, "tweet_id=") {
		parts := strings.Split(tweetLink, "tweet_id=")
		tweetID = strings.Split(parts[1], "&")[0]
	} else if strings.Contains(tweetLink, "?") {
		parts := strings.Split(tweetLink, "status/")
		tweetID = strings.Split(parts[1], "?")[0]
	} else if strings.Contains(tweetLink, "status/") {
		parts := strings.Split(tweetLink, "status/")
		tweetID = parts[1]
	} else {
		extra.Logger{}.Error("Failed to get tweet ID from your link: %s", tweetLink)
		return ""
	}

	return tweetID
}

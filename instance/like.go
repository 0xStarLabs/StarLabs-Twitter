package instance

import (
	"compress/gzip"
	"encoding/json"
	"fmt"
	http "github.com/Danny-Dasilva/fhttp"
	"io"
	"strings"
	"twitter/extra"
	"twitter/instance/additional_twitter_methods"
)

// Like press like button
func (twitter *Twitter) Like(tweetLink string) bool {
	likeURL := fmt.Sprintf("https://twitter.com/i/api/graphql/%s/FavoriteTweet", twitter.queryID.LikeID)
	tweetID := additional_twitter_methods.GetTweetID(tweetLink)
	if tweetID == "" {
		return false
	}

	for i := 0; i < twitter.config.Info.MaxTasksRetries; i++ {
		var stringData = fmt.Sprintf(`{"variables":{"tweet_id":"%s"},"queryId":"%s"}`, tweetID, twitter.queryID.LikeID)
		data := strings.NewReader(stringData)

		// Create new request
		req, err := http.NewRequest("POST", likeURL, data)
		if err != nil {
			extra.Logger{}.Error("Failed to build like request: %s", err.Error())
			continue
		}
		req.Header = http.Header{
			"accept":                {"*/*"},
			"accept-encoding":       {"gzip, deflate, br"},
			"authorization":         {twitter.queryID.BearerToken},
			"content-type":          {"application/json"},
			"cookie":                {twitter.cookies.CookiesToHeader()},
			"origin":                {"https://twitter.com"},
			"sec-ch-ua-mobile":      {"?0"},
			"sec-ch-ua-platform":    {`"Windows"`},
			"sec-fetch-dest":        {"empty"},
			"sec-fetch-mode":        {"cors"},
			"sec-fetch-site":        {"same-origin"},
			"x-csrf-token":          {twitter.ct0},
			"x-twitter-active-user": {"yes"},
			"x-twitter-auth-type":   {"OAuth2Session"},
			// x-client-transaction-id
			http.HeaderOrderKey: {
				"accept",
				"accept-encoding",
				"authorization",
				"content-type",
				"cookie",
				"origin",
				"sec-ch-ua-mobile",
				"sec-ch-ua-platform",
				"sec-fetch-dest",
				"sec-fetch-mode",
				"sec-fetch-site",
				"user-agent",
				"x-csrf-token",
				"x-twitter-active-user",
				"x-twitter-auth-type",
			},
			http.PHeaderOrderKey: {":authority", ":method", ":path", ":scheme"},
		}

		resp, err := twitter.client.Do(req)
		if err != nil {
			extra.Logger{}.Error("Failed to do like request: %s", err.Error())
			continue
		}
		defer resp.Body.Close()

		var reader io.ReadCloser
		switch resp.Header.Get("Content-Encoding") {
		case "gzip":
			reader, err = gzip.NewReader(resp.Body)
			if err != nil {
				extra.Logger{}.Error("Failed to create gzip reader while like: %s", err.Error())
				continue
			}
			defer reader.Close()
		default:
			reader = resp.Body
		}

		bodyBytes, err := io.ReadAll(reader)
		if err != nil {
			extra.Logger{}.Error("Failed to read like response body: %s", err.Error())
			continue
		}

		bodyString := string(bodyBytes)

		if resp.StatusCode >= 200 && resp.StatusCode <= 299 {
			if strings.Contains(bodyString, "already") {
				var responseAlreadyLike alreadyLikedResponse
				err = json.Unmarshal(bodyBytes, &responseAlreadyLike)
				if err != nil {
					extra.Logger{}.Error("%d | Failed to do unmarshal in already liked response: %s", twitter.index, err.Error())
					continue
				}
				extra.Logger{}.Success("%d | %s already liked tweet %s", twitter.index, twitter.Username, tweetID)
				return true
			} else if strings.Contains(bodyString, "favorite_tweet") {
				var responseLike likeResponse
				err = json.Unmarshal(bodyBytes, &responseLike)
				if err != nil {
					extra.Logger{}.Error("%d | Failed to do unmarshal in like response: %s", twitter.index, err.Error())
					continue
				}
				if responseLike.Data.FavoriteTweet == "Done" {
					extra.Logger{}.Success("%d | %s liked tweet %s", twitter.index, twitter.Username, tweetID)
					return true
				}
			}

		} else if strings.Contains(bodyString, "this account is temporarily locked") {
			extra.Logger{}.Error("%d | Account is temporarily locked", twitter.index)
			return false

		} else if strings.Contains(bodyString, "Could not authenticate you") {
			extra.Logger{}.Error("%d | Could not authenticate you", twitter.index)
			return false
		} else {
			extra.Logger{}.Error("%d | Unknown response while like: %s", twitter.index, bodyString)
			continue
		}
	}

	extra.Logger{}.Error("%d | Unable to do like", twitter.index)
	return false
}

type likeResponse struct {
	Data struct {
		FavoriteTweet string `json:"favorite_tweet"`
	} `json:"data"`
}

type alreadyLikedResponse struct {
	Errors []struct {
		Message   string `json:"message"`
		Locations []struct {
			Line   int `json:"line"`
			Column int `json:"column"`
		} `json:"locations"`
		Path       []string `json:"path"`
		Extensions struct {
			Name    string `json:"name"`
			Source  string `json:"source"`
			Code    int    `json:"code"`
			Kind    string `json:"kind"`
			Tracing struct {
				TraceID string `json:"trace_id"`
			} `json:"tracing"`
		} `json:"extensions"`
		Code    int    `json:"code"`
		Kind    string `json:"kind"`
		Name    string `json:"name"`
		Source  string `json:"source"`
		Tracing struct {
			TraceID string `json:"trace_id"`
		} `json:"tracing"`
	} `json:"errors"`
	Data struct {
	} `json:"data"`
}

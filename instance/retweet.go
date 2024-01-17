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

// Retweet perform retweet action
func (twitter *Twitter) Retweet(tweetLink string) bool {
	retweetURL := fmt.Sprintf("https://twitter.com/i/api/graphql/%s/CreateRetweet", twitter.queryID.RetweetID)
	tweetID := additional_twitter_methods.GetTweetID(tweetLink)
	if tweetID == "" {
		return false
	}

	for i := 0; i < twitter.config.Info.MaxTasksRetries; i++ {
		var stringData = fmt.Sprintf(`{"variables":{"tweet_id":"%s","dark_request":false},"queryId":"%s"}`, tweetID, twitter.queryID.RetweetID)
		data := strings.NewReader(stringData)

		// Create new request
		req, err := http.NewRequest("POST", retweetURL, data)
		if err != nil {
			extra.Logger{}.Error("Failed to build retweet request: %s", err.Error())
			continue
		}
		req.Header = http.Header{
			"accept":                {"*/*"},
			"accept-encoding":       {"gzip, deflate, br"},
			"authorization":         {twitter.queryID.BearerToken},
			"content-type":          {"application/json"},
			"cookie":                {twitter.cookies.CookiesToHeader()},
			"origin":                {"https://twitter.com"},
			"referer":               {tweetLink},
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
				"referer",
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
			extra.Logger{}.Error("Failed to do retweet request: %s", err.Error())
			continue
		}
		defer resp.Body.Close()

		var reader io.ReadCloser
		switch resp.Header.Get("Content-Encoding") {
		case "gzip":
			reader, err = gzip.NewReader(resp.Body)
			if err != nil {
				extra.Logger{}.Error("Failed to create gzip reader while retweet: %s", err.Error())
				continue
			}
			defer reader.Close()
		default:
			reader = resp.Body
		}

		bodyBytes, err := io.ReadAll(reader)
		if err != nil {
			extra.Logger{}.Error("Failed to read retweet response body: %s", err.Error())
			continue
		}

		bodyString := string(bodyBytes)

		if resp.StatusCode >= 200 && resp.StatusCode <= 299 {
			if strings.Contains(bodyString, "already") {
				var responseAlreadyLike alreadyLikedResponse
				err = json.Unmarshal(bodyBytes, &responseAlreadyLike)
				if err != nil {
					extra.Logger{}.Error("%d | Failed to do unmarshal in already retweeted response: %s", twitter.index, err.Error())
					continue
				}
				extra.Logger{}.Success("%d | %s already retweeted tweet %s", twitter.index, twitter.Username, tweetID)
				return true
			} else if strings.Contains(bodyString, "create_retweet") {
				var responseRetweet retweetResponse
				err = json.Unmarshal(bodyBytes, &responseRetweet)
				if err != nil {
					extra.Logger{}.Error("%d | Failed to do unmarshal in retweeted response: %s", twitter.index, err.Error())
					continue
				}
				extra.Logger{}.Success("%d | %s retweeted tweet %s", twitter.index, twitter.Username, tweetID)
				return true
			}

		} else if strings.Contains(bodyString, "this account is temporarily locked") {
			extra.Logger{}.Error("%d | Account is temporarily locked", twitter.index)
			return false

		} else if strings.Contains(bodyString, "Could not authenticate you") {
			extra.Logger{}.Error("%d | Could not authenticate you", twitter.index)
			return false
		} else {
			extra.Logger{}.Error("%d | Unknown response while retweet: %s", twitter.index, bodyString)
			continue
		}
	}

	extra.Logger{}.Error("%d | Unable to do retweet", twitter.index)
	return false
}

type retweetResponse struct {
	Data struct {
		CreateRetweet struct {
			RetweetResults struct {
				Result struct {
					RestID string `json:"rest_id"`
					Legacy struct {
						FullText string `json:"full_text"`
					} `json:"legacy"`
				} `json:"result"`
			} `json:"retweet_results"`
		} `json:"create_retweet"`
	} `json:"data"`
}

type alreadyRetweetedResponse struct {
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

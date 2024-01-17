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

// Comment perform comment action
func (twitter *Twitter) Comment(tweetContent string, tweetLink string) bool {
	tweetURL := fmt.Sprintf("https://twitter.com/i/api/graphql/%s/CreateTweet", twitter.queryID.TweetID)
	tweetID := additional_twitter_methods.GetTweetID(tweetLink)
	if tweetID == "" {
		return false
	}

	for i := 0; i < twitter.config.Info.MaxTasksRetries; i++ {
		var stringData = fmt.Sprintf(`{"variables":{"tweet_text":"%s","reply":{"in_reply_to_tweet_id":"%s","exclude_reply_user_ids":[]},"dark_request":false,"media":{"media_entities":[],"possibly_sensitive":false},"semantic_annotation_ids":[]},"features":{"c9s_tweet_anatomy_moderator_badge_enabled":true,"tweetypie_unmention_optimization_enabled":true,"responsive_web_edit_tweet_api_enabled":true,"graphql_is_translatable_rweb_tweet_is_translatable_enabled":true,"view_counts_everywhere_api_enabled":true,"longform_notetweets_consumption_enabled":true,"responsive_web_twitter_article_tweet_consumption_enabled":false,"tweet_awards_web_tipping_enabled":false,"longform_notetweets_rich_text_read_enabled":true,"longform_notetweets_inline_media_enabled":true,"rweb_video_timestamps_enabled":true,"responsive_web_graphql_exclude_directive_enabled":true,"verified_phone_label_enabled":false,"freedom_of_speech_not_reach_fetch_enabled":true,"standardized_nudges_misinfo":true,"tweet_with_visibility_results_prefer_gql_limited_actions_policy_enabled":true,"responsive_web_media_download_video_enabled":false,"responsive_web_graphql_skip_user_profile_image_extensions_enabled":false,"responsive_web_graphql_timeline_navigation_enabled":true,"responsive_web_enhance_cards_enabled":false},"queryId":"%s"}`, tweetContent, tweetID, twitter.queryID.TweetID)
		data := strings.NewReader(stringData)

		// Create new request
		req, err := http.NewRequest("POST", tweetURL, data)
		if err != nil {
			extra.Logger{}.Error("Failed to build comment request: %s", err.Error())
			continue
		}
		req.Header = http.Header{
			"accept":                {"*/*"},
			"accept-encoding":       {"gzip, deflate, br"},
			"authorization":         {twitter.queryID.BearerToken},
			"content-type":          {"application/json"},
			"cookie":                {twitter.cookies.CookiesToHeader()},
			"origin":                {"https://twitter.com"},
			"referer":               {"https://twitter.com/compose/tweet"},
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
			extra.Logger{}.Error("Failed to do comment request: %s", err.Error())
			continue
		}
		defer resp.Body.Close()

		var reader io.ReadCloser
		switch resp.Header.Get("Content-Encoding") {
		case "gzip":
			reader, err = gzip.NewReader(resp.Body)
			if err != nil {
				extra.Logger{}.Error("Failed to create gzip reader while comment: %s", err.Error())
				continue
			}
			defer reader.Close()
		default:
			reader = resp.Body
		}

		bodyBytes, err := io.ReadAll(reader)
		if err != nil {
			extra.Logger{}.Error("Failed to read comment response body: %s", err.Error())
			continue
		}

		bodyString := string(bodyBytes)

		if resp.StatusCode >= 200 && resp.StatusCode <= 299 {
			if strings.Contains(bodyString, "duplicate") {
				var responseTweetDuplicate tweetDuplicateResponse
				err = json.Unmarshal(bodyBytes, &responseTweetDuplicate)
				if err != nil {
					extra.Logger{}.Error("%d | Failed to do unmarshal in duplicate comment response: %s", twitter.index, err.Error())
					continue
				}
				extra.Logger{}.Success("%d | %s already commented this tweet. Duplicate", twitter.index, twitter.Username)
				return true

			} else if strings.Contains(bodyString, "rest_id") {
				var responseTweet tweetResponse
				err = json.Unmarshal(bodyBytes, &responseTweet)
				if err != nil {
					extra.Logger{}.Error("%d | Failed to do unmarshal in comment response: %s", twitter.index, err.Error())
					continue
				}
				extra.Logger{}.Success("%d | %s successfully commented", twitter.index, twitter.Username)
				return true
			}

		} else if strings.Contains(bodyString, "this account is temporarily locked") {
			extra.Logger{}.Error("%d | Account is temporarily locked", twitter.index)
			return false

		} else if strings.Contains(bodyString, "Could not authenticate you") {
			extra.Logger{}.Error("%d | Could not authenticate you", twitter.index)
			return false
		} else {
			extra.Logger{}.Error("%d | Unknown response while comment: %s", twitter.index, bodyString)
			continue
		}
	}

	extra.Logger{}.Error("%d | Unable to do comment", twitter.index)
	return false
}

func (twitter *Twitter) CommentWithPicture(tweetContent string, tweetLink string, pictureBase64Encoded string) bool {
	tweetURL := fmt.Sprintf("https://twitter.com/i/api/graphql/%s/CreateTweet", twitter.queryID.TweetID)
	tweetID := additional_twitter_methods.GetTweetID(tweetLink)
	if tweetID == "" {
		return false
	}

	for i := 0; i < twitter.config.Info.MaxTasksRetries; i++ {
		mediaID := twitter.UploadImageToTwitter(pictureBase64Encoded)
		if mediaID == "" {
			return false
		}

		var stringData = fmt.Sprintf(`{"variables":{"tweet_text":"%s","reply":{"in_reply_to_tweet_id":"%s","exclude_reply_user_ids":[]},"dark_request":false,"media":{"media_entities":[{"media_id":"%s","tagged_users":[]}],"possibly_sensitive":false},"semantic_annotation_ids":[]},"features":{"c9s_tweet_anatomy_moderator_badge_enabled":true,"tweetypie_unmention_optimization_enabled":true,"responsive_web_edit_tweet_api_enabled":true,"graphql_is_translatable_rweb_tweet_is_translatable_enabled":true,"view_counts_everywhere_api_enabled":true,"longform_notetweets_consumption_enabled":true,"responsive_web_twitter_article_tweet_consumption_enabled":false,"tweet_awards_web_tipping_enabled":false,"longform_notetweets_rich_text_read_enabled":true,"longform_notetweets_inline_media_enabled":true,"rweb_video_timestamps_enabled":true,"responsive_web_graphql_exclude_directive_enabled":true,"verified_phone_label_enabled":false,"freedom_of_speech_not_reach_fetch_enabled":true,"standardized_nudges_misinfo":true,"tweet_with_visibility_results_prefer_gql_limited_actions_policy_enabled":true,"responsive_web_media_download_video_enabled":false,"responsive_web_graphql_skip_user_profile_image_extensions_enabled":false,"responsive_web_graphql_timeline_navigation_enabled":true,"responsive_web_enhance_cards_enabled":false},"queryId":"%s"}`, tweetContent, tweetID, mediaID, twitter.queryID.TweetID)
		data := strings.NewReader(stringData)

		// Create new request
		req, err := http.NewRequest("POST", tweetURL, data)
		if err != nil {
			extra.Logger{}.Error("Failed to build comment with picture request: %s", err.Error())
			continue
		}
		req.Header = http.Header{
			"accept":                {"*/*"},
			"accept-encoding":       {"gzip, deflate, br"},
			"authorization":         {twitter.queryID.BearerToken},
			"content-type":          {"application/json"},
			"cookie":                {twitter.cookies.CookiesToHeader()},
			"origin":                {"https://twitter.com"},
			"referer":               {"https://twitter.com/compose/tweet"},
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
			extra.Logger{}.Error("Failed to do comment with picture request: %s", err.Error())
			continue
		}
		defer resp.Body.Close()

		var reader io.ReadCloser
		switch resp.Header.Get("Content-Encoding") {
		case "gzip":
			reader, err = gzip.NewReader(resp.Body)
			if err != nil {
				extra.Logger{}.Error("Failed to create gzip reader while comment with picture: %s", err.Error())
				continue
			}
			defer reader.Close()
		default:
			reader = resp.Body
		}

		bodyBytes, err := io.ReadAll(reader)
		if err != nil {
			extra.Logger{}.Error("Failed to read comment with picture response body: %s", err.Error())
			continue
		}

		bodyString := string(bodyBytes)

		if resp.StatusCode >= 200 && resp.StatusCode <= 299 {
			if strings.Contains(bodyString, "rest_id") {
				var responseTweet tweetResponse
				err = json.Unmarshal(bodyBytes, &responseTweet)
				if err != nil {
					extra.Logger{}.Error("%d | Failed to do unmarshal in comment with picture response: %s", twitter.index, err.Error())
					continue
				}
				extra.Logger{}.Success("%d | %s successfully commented with picture", twitter.index, twitter.Username)
				return true
			}

		} else if strings.Contains(bodyString, "this account is temporarily locked") {
			extra.Logger{}.Error("%d | Account is temporarily locked", twitter.index)
			return false

		} else if strings.Contains(bodyString, "Could not authenticate you") {
			extra.Logger{}.Error("%d | Could not authenticate you", twitter.index)
			return false
		} else {
			extra.Logger{}.Error("%d | Unknown response while comment with picture: %s", twitter.index, bodyString)
			continue
		}
	}

	extra.Logger{}.Error("%d | Unable to do comment with picture", twitter.index)
	return false
}

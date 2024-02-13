package instance

import (
	"errors"
	"fmt"
	http "github.com/bogdanfinn/fhttp"
	"io"
	"strings"
	"twitter/extra"
	"twitter/instance/additional_twitter_methods"
)

func (twitter *Twitter) VotePoll(link, answer string) bool {
	for i := 0; i < twitter.config.Info.MaxTasksRetries; i++ {
		tweetID := additional_twitter_methods.GetTweetID(link)
		err, tweetDetail := twitter.getTweetDetail(link)

		pollName := "poll" + strings.Split(strings.Split(tweetDetail, `"name":"poll`)[1], `"`)[0]
		cardID := strings.Split(strings.Split(tweetDetail, `card://`)[1], `"`)[0]

		var stringData = fmt.Sprintf(`twitter%%3Astring%%3Acard_uri=card%%3A%%2F%%2F%s&twitter%%3Along%%3Aoriginal_tweet_id=%s&twitter%%3Astring%%3Aresponse_card_name=%s&twitter%%3Astring%%3Acards_platform=Web-12&twitter%%3Astring%%3Aselected_choice=%s`, cardID, tweetID, pollName, answer)
		data := strings.NewReader(stringData)

		// Create new request
		req, err := http.NewRequest("POST", "https://caps.twitter.com/v2/capi/passthrough/1", data)
		if err != nil {
			extra.Logger{}.Error("Failed to build vote poll request: %s", err.Error())
			continue
		}

		req.Header = http.Header{
			"accept":                {"*/*"},
			"accept-encoding":       {"gzip, deflate, br"},
			"authorization":         {twitter.queryID.BearerToken},
			"content-type":          {"application/x-www-form-urlencoded"},
			"cookie":                {twitter.cookies.CookiesToHeader()},
			"origin":                {"https://twitter.com"},
			"referer":               {"https://twitter.com/"},
			"sec-ch-ua-mobile":      {"?0"},
			"sec-ch-ua-platform":    {`"Windows"`},
			"sec-fetch-dest":        {"empty"},
			"sec-fetch-mode":        {"cors"},
			"sec-fetch-site":        {"same-site"},
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
			extra.Logger{}.Error("Failed to do vote poll request: %s", err.Error())
			continue
		}
		defer resp.Body.Close()

		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			extra.Logger{}.Error("Failed to read vote poll response body: %s", err.Error())
			continue
		}

		bodyString := string(bodyBytes)

		if resp.StatusCode >= 200 && resp.StatusCode <= 299 {
			extra.Logger{}.Success("%d | %s voted in poll %s", twitter.index, twitter.Username, tweetID)
			return true

		} else if strings.Contains(bodyString, "this account is temporarily locked") {
			extra.Logger{}.Error("%d | Account is temporarily locked", twitter.index)
			return false

		} else if strings.Contains(bodyString, "Could not authenticate you") {
			extra.Logger{}.Error("%d | Could not authenticate you", twitter.index)
			return false
		} else {
			extra.Logger{}.Error("%d | Unknown response while vote poll: %s", twitter.index, bodyString)
			continue
		}
	}

	extra.Logger{}.Error("%d | Unable to do vote poll", twitter.index)
	return false
}

func (twitter *Twitter) getTweetDetail(link string) (error, string) {
	for i := 0; i < twitter.config.Info.MaxTasksRetries; i++ {
		tweetID := additional_twitter_methods.GetTweetID(link)
		url := fmt.Sprintf(`https://twitter.com/i/api/graphql/B9_KmbkLhXt6jRwGjJrweg/TweetDetail?variables=%%7B%%22focalTweetId%%22%%3A%%22%s%%22%%2C%%22with_rux_injections%%22%%3Afalse%%2C%%22includePromotedContent%%22%%3Atrue%%2C%%22withCommunity%%22%%3Atrue%%2C%%22withQuickPromoteEligibilityTweetFields%%22%%3Atrue%%2C%%22withBirdwatchNotes%%22%%3Atrue%%2C%%22withVoice%%22%%3Atrue%%2C%%22withV2Timeline%%22%%3Atrue%%7D&features=%%7B%%22responsive_web_graphql_exclude_directive_enabled%%22%%3Atrue%%2C%%22verified_phone_label_enabled%%22%%3Afalse%%2C%%22creator_subscriptions_tweet_preview_api_enabled%%22%%3Atrue%%2C%%22responsive_web_graphql_timeline_navigation_enabled%%22%%3Atrue%%2C%%22responsive_web_graphql_skip_user_profile_image_extensions_enabled%%22%%3Afalse%%2C%%22c9s_tweet_anatomy_moderator_badge_enabled%%22%%3Atrue%%2C%%22tweetypie_unmention_optimization_enabled%%22%%3Atrue%%2C%%22responsive_web_edit_tweet_api_enabled%%22%%3Atrue%%2C%%22graphql_is_translatable_rweb_tweet_is_translatable_enabled%%22%%3Atrue%%2C%%22view_counts_everywhere_api_enabled%%22%%3Atrue%%2C%%22longform_notetweets_consumption_enabled%%22%%3Atrue%%2C%%22responsive_web_twitter_article_tweet_consumption_enabled%%22%%3Atrue%%2C%%22tweet_awards_web_tipping_enabled%%22%%3Afalse%%2C%%22freedom_of_speech_not_reach_fetch_enabled%%22%%3Atrue%%2C%%22standardized_nudges_misinfo%%22%%3Atrue%%2C%%22tweet_with_visibility_results_prefer_gql_limited_actions_policy_enabled%%22%%3Atrue%%2C%%22rweb_video_timestamps_enabled%%22%%3Atrue%%2C%%22longform_notetweets_rich_text_read_enabled%%22%%3Atrue%%2C%%22longform_notetweets_inline_media_enabled%%22%%3Atrue%%2C%%22responsive_web_media_download_video_enabled%%22%%3Atrue%%2C%%22responsive_web_enhance_cards_enabled%%22%%3Afalse%%7D&fieldToggles=%%7B%%22withArticleRichContentState%%22%%3Atrue%%7D`, tweetID)

		// Create new request
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			extra.Logger{}.Error("Failed to build tweet detail request: %s", err.Error())
			continue
		}

		req.Header.Set("authority", "twitter.com")
		req.Header.Set("accept", "*/*")
		req.Header.Set("authorization", twitter.queryID.BearerToken)
		req.Header.Set("content-type", "application/json")
		req.Header.Set("cookie", twitter.cookies.CookiesToHeader())
		req.Header.Set("referer", link)
		req.Header.Set("sec-ch-ua-mobile", "?0")
		req.Header.Set("sec-ch-ua-platform", `"Windows"`)
		req.Header.Set("sec-fetch-dest", "empty")
		req.Header.Set("sec-fetch-mode", "cors")
		req.Header.Set("sec-fetch-site", "same-origin")
		req.Header.Set("x-client-uuid", "0f946cb6-e4db-4b2c-9f1b-a5bf1b46b13c")
		req.Header.Set("x-csrf-token", twitter.ct0)
		req.Header.Set("x-twitter-active-user", "yes")
		req.Header.Set("x-twitter-auth-type", "OAuth2Session")
		req.Header.Set("x-twitter-client-language", "en")

		resp, err := twitter.client.Do(req)
		if err != nil {
			extra.Logger{}.Error("Failed to do tweet detail request: %s", err.Error())
			continue
		}
		defer resp.Body.Close()

		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			extra.Logger{}.Error("Failed to read tweet detail response body: %s", err.Error())
			continue
		}

		return nil, string(bodyBytes)

	}

	twitter.logger.Error("%d | Failed to get tweet detail.", twitter.index)
	return errors.New(""), ""
}

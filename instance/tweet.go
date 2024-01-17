package instance

import (
	"compress/gzip"
	"encoding/json"
	"fmt"
	http "github.com/Danny-Dasilva/fhttp"
	"io"
	"net/url"
	"strings"
	"twitter/extra"
)

// Tweet perform tweet action
func (twitter *Twitter) Tweet(tweetContent string) bool {
	tweetURL := fmt.Sprintf("https://twitter.com/i/api/graphql/%s/CreateTweet", twitter.queryID.TweetID)

	for i := 0; i < twitter.config.Info.MaxTasksRetries; i++ {
		var stringData = fmt.Sprintf(`{"variables":{"tweet_text":"%s","dark_request":false,"media":{"media_entities":[],"possibly_sensitive":false},"semantic_annotation_ids":[]},"features":{"c9s_tweet_anatomy_moderator_badge_enabled":true,"tweetypie_unmention_optimization_enabled":true,"responsive_web_edit_tweet_api_enabled":true,"graphql_is_translatable_rweb_tweet_is_translatable_enabled":true,"view_counts_everywhere_api_enabled":true,"longform_notetweets_consumption_enabled":true,"responsive_web_twitter_article_tweet_consumption_enabled":false,"tweet_awards_web_tipping_enabled":false,"longform_notetweets_rich_text_read_enabled":true,"longform_notetweets_inline_media_enabled":true,"rweb_video_timestamps_enabled":true,"responsive_web_graphql_exclude_directive_enabled":true,"verified_phone_label_enabled":false,"freedom_of_speech_not_reach_fetch_enabled":true,"standardized_nudges_misinfo":true,"tweet_with_visibility_results_prefer_gql_limited_actions_policy_enabled":true,"responsive_web_media_download_video_enabled":false,"responsive_web_graphql_skip_user_profile_image_extensions_enabled":false,"responsive_web_graphql_timeline_navigation_enabled":true,"responsive_web_enhance_cards_enabled":false},"queryId":"%s"}`, tweetContent, twitter.queryID.TweetID)
		data := strings.NewReader(stringData)

		// Create new request
		req, err := http.NewRequest("POST", tweetURL, data)
		if err != nil {
			extra.Logger{}.Error("Failed to build tweet request: %s", err.Error())
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
			extra.Logger{}.Error("Failed to do tweet request: %s", err.Error())
			continue
		}
		defer resp.Body.Close()

		var reader io.ReadCloser
		switch resp.Header.Get("Content-Encoding") {
		case "gzip":
			reader, err = gzip.NewReader(resp.Body)
			if err != nil {
				extra.Logger{}.Error("Failed to create gzip reader while tweet: %s", err.Error())
				continue
			}
			defer reader.Close()
		default:
			reader = resp.Body
		}

		bodyBytes, err := io.ReadAll(reader)
		if err != nil {
			extra.Logger{}.Error("Failed to read tweet response body: %s", err.Error())
			continue
		}

		bodyString := string(bodyBytes)

		if resp.StatusCode >= 200 && resp.StatusCode <= 299 {
			if strings.Contains(bodyString, "duplicate") {
				var responseTweetDuplicate tweetDuplicateResponse
				err = json.Unmarshal(bodyBytes, &responseTweetDuplicate)
				if err != nil {
					extra.Logger{}.Error("%d | Failed to do unmarshal in duplicate tweet response: %s", twitter.index, err.Error())
					continue
				}
				extra.Logger{}.Success("%d | %s already tweeted this tweet. Duplicate", twitter.index, twitter.Username)
				return true

			} else if strings.Contains(bodyString, "rest_id") {
				var responseTweet tweetResponse
				err = json.Unmarshal(bodyBytes, &responseTweet)
				if err != nil {
					extra.Logger{}.Error("%d | Failed to do unmarshal in tweet response: %s", twitter.index, err.Error())
					continue
				}
				extra.Logger{}.Success("%d | %s successfully tweeted", twitter.index, twitter.Username)
				return true
			}

		} else if strings.Contains(bodyString, "this account is temporarily locked") {
			extra.Logger{}.Error("%d | Account is temporarily locked", twitter.index)
			return false

		} else if strings.Contains(bodyString, "Could not authenticate you") {
			extra.Logger{}.Error("%d | Could not authenticate you", twitter.index)
			return false
		} else {
			extra.Logger{}.Error("%d | Unknown response while tweet: %s", twitter.index, bodyString)
			continue
		}
	}

	extra.Logger{}.Error("%d | Unable to do tweet", twitter.index)
	return false
}

func (twitter *Twitter) TweetWithPicture(tweetContent string, pictureBase64Encoded string) bool {
	tweetURL := fmt.Sprintf("https://twitter.com/i/api/graphql/%s/CreateTweet", twitter.queryID.TweetID)

	for i := 0; i < twitter.config.Info.MaxTasksRetries; i++ {
		mediaID := twitter.UploadImageToTwitter(pictureBase64Encoded)
		if mediaID == "" {
			continue
		}

		var stringData = fmt.Sprintf(`{"variables":{"tweet_text":"%s","dark_request":false,"media":{"media_entities":[{"media_id":"%s","tagged_users":[]}],"possibly_sensitive":false},"semantic_annotation_ids":[]},"features":{"c9s_tweet_anatomy_moderator_badge_enabled":true,"tweetypie_unmention_optimization_enabled":true,"responsive_web_edit_tweet_api_enabled":true,"graphql_is_translatable_rweb_tweet_is_translatable_enabled":true,"view_counts_everywhere_api_enabled":true,"longform_notetweets_consumption_enabled":true,"responsive_web_twitter_article_tweet_consumption_enabled":false,"tweet_awards_web_tipping_enabled":false,"longform_notetweets_rich_text_read_enabled":true,"longform_notetweets_inline_media_enabled":true,"rweb_video_timestamps_enabled":true,"responsive_web_graphql_exclude_directive_enabled":true,"verified_phone_label_enabled":false,"freedom_of_speech_not_reach_fetch_enabled":true,"standardized_nudges_misinfo":true,"tweet_with_visibility_results_prefer_gql_limited_actions_policy_enabled":true,"responsive_web_media_download_video_enabled":false,"responsive_web_graphql_skip_user_profile_image_extensions_enabled":false,"responsive_web_graphql_timeline_navigation_enabled":true,"responsive_web_enhance_cards_enabled":false},"queryId":"%s"}`, tweetContent, mediaID, twitter.queryID.TweetID)
		data := strings.NewReader(stringData)

		// Create new request
		req, err := http.NewRequest("POST", tweetURL, data)
		if err != nil {
			extra.Logger{}.Error("Failed to build tweet with picture request: %s", err.Error())
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
			extra.Logger{}.Error("Failed to do with picture request: %s", err.Error())
			continue
		}
		defer resp.Body.Close()

		var reader io.ReadCloser
		switch resp.Header.Get("Content-Encoding") {
		case "gzip":
			reader, err = gzip.NewReader(resp.Body)
			if err != nil {
				extra.Logger{}.Error("Failed to create gzip reader while with picture: %s", err.Error())
				continue
			}
			defer reader.Close()
		default:
			reader = resp.Body
		}

		bodyBytes, err := io.ReadAll(reader)
		if err != nil {
			extra.Logger{}.Error("Failed to read with picture response body: %s", err.Error())
			continue
		}

		bodyString := string(bodyBytes)

		if resp.StatusCode >= 200 && resp.StatusCode <= 299 {
			if strings.Contains(bodyString, "rest_id") {
				var responseTweet tweetResponse
				err = json.Unmarshal(bodyBytes, &responseTweet)
				if err != nil {
					extra.Logger{}.Error("%d | Failed to do unmarshal in with picture response: %s", twitter.index, err.Error())
					continue
				}
				extra.Logger{}.Success("%d | %s successfully tweeted with picture", twitter.index, twitter.Username)
				return true
			}

		} else if strings.Contains(bodyString, "this account is temporarily locked") {
			extra.Logger{}.Error("%d | Account is temporarily locked", twitter.index)
			return false

		} else if strings.Contains(bodyString, "Could not authenticate you") {
			extra.Logger{}.Error("%d | Could not authenticate you", twitter.index)
			return false
		} else {
			extra.Logger{}.Error("%d | Unknown response while tweet with picture: %s", twitter.index, bodyString)
			continue
		}
	}

	extra.Logger{}.Error("%d | Unable to do tweet with picture", twitter.index)
	return false
}

func (twitter *Twitter) TweetQuote(tweetContent string, tweetLink string) bool {
	tweetURL := fmt.Sprintf("https://twitter.com/i/api/graphql/%s/CreateTweet", twitter.queryID.TweetID)

	for i := 0; i < twitter.config.Info.MaxTasksRetries; i++ {
		var stringData = fmt.Sprintf(`{"variables":{"tweet_text":"%s","attachment_url":"%s","dark_request":false,"media":{"media_entities":[],"possibly_sensitive":false},"semantic_annotation_ids":[]},"features":{"c9s_tweet_anatomy_moderator_badge_enabled":true,"tweetypie_unmention_optimization_enabled":true,"responsive_web_edit_tweet_api_enabled":true,"graphql_is_translatable_rweb_tweet_is_translatable_enabled":true,"view_counts_everywhere_api_enabled":true,"longform_notetweets_consumption_enabled":true,"responsive_web_twitter_article_tweet_consumption_enabled":false,"tweet_awards_web_tipping_enabled":false,"longform_notetweets_rich_text_read_enabled":true,"longform_notetweets_inline_media_enabled":true,"rweb_video_timestamps_enabled":true,"responsive_web_graphql_exclude_directive_enabled":true,"verified_phone_label_enabled":false,"freedom_of_speech_not_reach_fetch_enabled":true,"standardized_nudges_misinfo":true,"tweet_with_visibility_results_prefer_gql_limited_actions_policy_enabled":true,"responsive_web_media_download_video_enabled":false,"responsive_web_graphql_skip_user_profile_image_extensions_enabled":false,"responsive_web_graphql_timeline_navigation_enabled":true,"responsive_web_enhance_cards_enabled":false},"queryId":"%s"}`, tweetContent, tweetLink, twitter.queryID.TweetID)
		data := strings.NewReader(stringData)

		// Create new request
		req, err := http.NewRequest("POST", tweetURL, data)
		if err != nil {
			extra.Logger{}.Error("Failed to build quote tweet request: %s", err.Error())
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
			extra.Logger{}.Error("Failed to do quote tweet request: %s", err.Error())
			continue
		}
		defer resp.Body.Close()

		var reader io.ReadCloser
		switch resp.Header.Get("Content-Encoding") {
		case "gzip":
			reader, err = gzip.NewReader(resp.Body)
			if err != nil {
				extra.Logger{}.Error("Failed to create gzip reader while quote tweet: %s", err.Error())
				continue
			}
			defer reader.Close()
		default:
			reader = resp.Body
		}

		bodyBytes, err := io.ReadAll(reader)
		if err != nil {
			extra.Logger{}.Error("Failed to read quote tweet response body: %s", err.Error())
			continue
		}

		bodyString := string(bodyBytes)

		if resp.StatusCode >= 200 && resp.StatusCode <= 299 {
			if strings.Contains(bodyString, "duplicate") {
				var responseTweetDuplicate tweetDuplicateResponse
				err = json.Unmarshal(bodyBytes, &responseTweetDuplicate)
				if err != nil {
					extra.Logger{}.Error("%d | Failed to do unmarshal in duplicate quote tweet response: %s", twitter.index, err.Error())
					continue
				}
				extra.Logger{}.Success("%d | %s already tweeted this tweet. Duplicate", twitter.index, twitter.Username)
				return true

			} else if strings.Contains(bodyString, "rest_id") {
				var responseTweet tweetResponse
				err = json.Unmarshal(bodyBytes, &responseTweet)
				if err != nil {
					extra.Logger{}.Error("%d | Failed to do unmarshal in quote tweet response: %s", twitter.index, err.Error())
					continue
				}
				extra.Logger{}.Success("%d | %s successfully tweeted", twitter.index, twitter.Username)
				return true
			}

		} else if strings.Contains(bodyString, "this account is temporarily locked") {
			extra.Logger{}.Error("%d | Account is temporarily locked", twitter.index)
			return false

		} else if strings.Contains(bodyString, "Could not authenticate you") {
			extra.Logger{}.Error("%d | Could not authenticate you", twitter.index)
			return false
		} else {
			extra.Logger{}.Error("%d | Unknown response while tweet: %s", twitter.index, bodyString)
			continue
		}
	}

	extra.Logger{}.Error("%d | Unable to do quote tweet", twitter.index)
	return false
}

func (twitter *Twitter) UploadImageToTwitter(pictureBase64Encoded string) string {
	data := url.Values{}
	data.Set("media_data", pictureBase64Encoded)

	mediaURL := "https://upload.twitter.com/1.1/media/upload.json"
	req, err := http.NewRequest("POST", mediaURL, strings.NewReader(data.Encode()))
	if err != nil {
		twitter.logger.Error("%d | Failed to build tweet media request: %s", twitter.index, err)
	}
	req.Header = http.Header{
		"authorization": {twitter.queryID.BearerToken},
		"content-type":  {"application/x-www-form-urlencoded"},
		"cookie":        {twitter.cookies.CookiesToHeader()},
		"x-csrf-token":  {twitter.ct0},
	}

	resp, err := twitter.client.Do(req)
	if err != nil {
		extra.Logger{}.Error("%d | Failed to do tweet media request: %s", twitter.index, err.Error())
		return ""
	}
	defer resp.Body.Close()

	var newReader io.Reader

	if resp.Header.Get("Content-Encoding") == "gzip" {
		gzipReader, err := gzip.NewReader(resp.Body)
		if err != nil {
			extra.Logger{}.Error("%d | Failed to gzip upload media response: %s", twitter.index, err.Error())
			return ""
		}
		defer gzipReader.Close()
		newReader = gzipReader
	} else {
		newReader = resp.Body
	}

	body, err := io.ReadAll(newReader)
	if err != nil {
		extra.Logger{}.Error("%d | Failed to read media upload response body: %s", twitter.index, err.Error())
		return ""
	}

	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		extra.Logger{}.Error("%d | Failed to unmarshal media upload response: %s", twitter.index, err.Error())
		return ""
	}
	mediaID, ok := result["media_id_string"].(string)
	if !ok {
		extra.Logger{}.Error("%d | Failed to get media ID in upload response: %s", twitter.index, result)
	}

	return mediaID
}

type tweetResponse struct {
	Data struct {
		CreateTweet struct {
			TweetResults struct {
				Result struct {
					RestID string `json:"rest_id"`
					Core   struct {
						UserResults struct {
							Result struct {
								Typename                   string `json:"__typename"`
								ID                         string `json:"id"`
								RestID                     string `json:"rest_id"`
								AffiliatesHighlightedLabel struct {
								} `json:"affiliates_highlighted_label"`
								HasGraduatedAccess bool   `json:"has_graduated_access"`
								IsBlueVerified     bool   `json:"is_blue_verified"`
								ProfileImageShape  string `json:"profile_image_shape"`
								Legacy             struct {
									CanDm               bool   `json:"can_dm"`
									CanMediaTag         bool   `json:"can_media_tag"`
									CreatedAt           string `json:"created_at"`
									DefaultProfile      bool   `json:"default_profile"`
									DefaultProfileImage bool   `json:"default_profile_image"`
									Description         string `json:"description"`
									Entities            struct {
										Description struct {
											Urls []any `json:"urls"`
										} `json:"description"`
									} `json:"entities"`
									FastFollowersCount      int    `json:"fast_followers_count"`
									FavouritesCount         int    `json:"favourites_count"`
									FollowersCount          int    `json:"followers_count"`
									FriendsCount            int    `json:"friends_count"`
									HasCustomTimelines      bool   `json:"has_custom_timelines"`
									IsTranslator            bool   `json:"is_translator"`
									ListedCount             int    `json:"listed_count"`
									Location                string `json:"location"`
									MediaCount              int    `json:"media_count"`
									Name                    string `json:"name"`
									NeedsPhoneVerification  bool   `json:"needs_phone_verification"`
									NormalFollowersCount    int    `json:"normal_followers_count"`
									PinnedTweetIdsStr       []any  `json:"pinned_tweet_ids_str"`
									PossiblySensitive       bool   `json:"possibly_sensitive"`
									ProfileImageURLHTTPS    string `json:"profile_image_url_https"`
									ProfileInterstitialType string `json:"profile_interstitial_type"`
									ScreenName              string `json:"screen_name"`
									StatusesCount           int    `json:"statuses_count"`
									TranslatorType          string `json:"translator_type"`
									Verified                bool   `json:"verified"`
									WantRetweets            bool   `json:"want_retweets"`
									WithheldInCountries     []any  `json:"withheld_in_countries"`
								} `json:"legacy"`
							} `json:"result"`
						} `json:"user_results"`
					} `json:"core"`
					UnmentionData struct {
					} `json:"unmention_data"`
					EditControl struct {
						EditTweetIds       []string `json:"edit_tweet_ids"`
						EditableUntilMsecs string   `json:"editable_until_msecs"`
						IsEditEligible     bool     `json:"is_edit_eligible"`
						EditsRemaining     string   `json:"edits_remaining"`
					} `json:"edit_control"`
					IsTranslatable bool `json:"is_translatable"`
					Views          struct {
						State string `json:"state"`
					} `json:"views"`
					Source string `json:"source"`
					Legacy struct {
						BookmarkCount     int    `json:"bookmark_count"`
						Bookmarked        bool   `json:"bookmarked"`
						CreatedAt         string `json:"created_at"`
						ConversationIDStr string `json:"conversation_id_str"`
						DisplayTextRange  []int  `json:"display_text_range"`
						Entities          struct {
							Hashtags     []any `json:"hashtags"`
							Symbols      []any `json:"symbols"`
							Timestamps   []any `json:"timestamps"`
							Urls         []any `json:"urls"`
							UserMentions []any `json:"user_mentions"`
						} `json:"entities"`
						FavoriteCount int    `json:"favorite_count"`
						Favorited     bool   `json:"favorited"`
						FullText      string `json:"full_text"`
						IsQuoteStatus bool   `json:"is_quote_status"`
						Lang          string `json:"lang"`
						QuoteCount    int    `json:"quote_count"`
						ReplyCount    int    `json:"reply_count"`
						RetweetCount  int    `json:"retweet_count"`
						Retweeted     bool   `json:"retweeted"`
						UserIDStr     string `json:"user_id_str"`
						IDStr         string `json:"id_str"`
					} `json:"legacy"`
					UnmentionInfo struct {
					} `json:"unmention_info"`
				} `json:"result"`
			} `json:"tweet_results"`
		} `json:"create_tweet"`
	} `json:"data"`
}

type tweetDuplicateResponse struct {
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

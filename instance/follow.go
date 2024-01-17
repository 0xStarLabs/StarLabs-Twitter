package instance

import (
	"compress/gzip"
	"encoding/json"
	"fmt"
	http "github.com/Danny-Dasilva/fhttp"
	"io"
	"strings"
	"twitter/extra"
)

// Follow make subscribe to given username
func (twitter *Twitter) Follow(usernameToFollow string) bool {
	for i := 0; i < twitter.config.Info.MaxTasksRetries; i++ {
		var stringData = fmt.Sprintf(`include_profile_interstitial_type=1&include_blocking=1&include_blocked_by=1&include_followed_by=1&include_want_retweets=1&include_mute_edge=1&include_can_dm=1&include_can_media_tag=1&include_ext_has_nft_avatar=1&include_ext_is_blue_verified=1&include_ext_verified_type=1&include_ext_profile_image_shape=1&skip_status=1&screen_name=%s`, usernameToFollow)
		data := strings.NewReader(stringData)

		// Create new request
		req, err := http.NewRequest("POST", "https://twitter.com/i/api/1.1/friendships/create.json", data)
		if err != nil {
			extra.Logger{}.Error("Failed to build follow request: %s", err.Error())
			continue
		}

		req.Header = http.Header{
			"accept":                {"*/*"},
			"accept-encoding":       {"gzip, deflate, br"},
			"authorization":         {twitter.queryID.BearerToken},
			"content-type":          {"application/x-www-form-urlencoded"},
			"cookie":                {twitter.cookies.CookiesToHeader()},
			"origin":                {"https://twitter.com"},
			"referer":               {fmt.Sprintf("https://twitter.com/%s", usernameToFollow)},
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
			extra.Logger{}.Error("Failed to do follow request: %s", err.Error())
			continue
		}
		defer resp.Body.Close()

		var reader io.ReadCloser
		switch resp.Header.Get("Content-Encoding") {
		case "gzip":
			reader, err = gzip.NewReader(resp.Body)
			if err != nil {
				extra.Logger{}.Error("Failed to create gzip reader while follow: %s", err.Error())
				continue
			}
			defer reader.Close()
		default:
			reader = resp.Body
		}

		bodyBytes, err := io.ReadAll(reader)
		if err != nil {
			extra.Logger{}.Error("Failed to read follow response body: %s", err.Error())
			continue
		}

		bodyString := string(bodyBytes)

		if strings.Contains(bodyString, "screen_name") && resp.StatusCode >= 200 && resp.StatusCode <= 299 {
			var responseDataUsername followResponse
			err = json.Unmarshal(bodyBytes, &responseDataUsername)
			if err != nil {
				extra.Logger{}.Error("%d | Failed to do unmarshal in follow response: %s", twitter.index, err.Error())
				continue
			}
			extra.Logger{}.Success("%d | %s subscribed to %s", twitter.index, twitter.Username, usernameToFollow)
			return true

		} else if strings.Contains(bodyString, "this account is temporarily locked") {
			extra.Logger{}.Error("%d | Account is temporarily locked", twitter.index)
			return false

		} else if strings.Contains(bodyString, "Could not authenticate you") {
			extra.Logger{}.Error("%d | Could not authenticate you", twitter.index)
			return false
		} else {
			extra.Logger{}.Error("%d | Unknown response while follow: %s", twitter.index, bodyString)
			continue
		}
	}

	extra.Logger{}.Error("%d | Unable to do follow", twitter.index)
	return false
}

type followResponse struct {
	ID          int    `json:"id"`
	IDStr       string `json:"id_str"`
	Name        string `json:"name"`
	ScreenName  string `json:"screen_name"`
	Location    string `json:"location"`
	Description string `json:"description"`
	URL         any    `json:"url"`
	Entities    struct {
		Description struct {
			Urls []any `json:"urls"`
		} `json:"description"`
	} `json:"entities"`
	Protected                      bool     `json:"protected"`
	FollowersCount                 int      `json:"followers_count"`
	FastFollowersCount             int      `json:"fast_followers_count"`
	NormalFollowersCount           int      `json:"normal_followers_count"`
	FriendsCount                   int      `json:"friends_count"`
	ListedCount                    int      `json:"listed_count"`
	CreatedAt                      string   `json:"created_at"`
	FavouritesCount                int      `json:"favourites_count"`
	UtcOffset                      any      `json:"utc_offset"`
	TimeZone                       any      `json:"time_zone"`
	GeoEnabled                     bool     `json:"geo_enabled"`
	Verified                       bool     `json:"verified"`
	StatusesCount                  int      `json:"statuses_count"`
	MediaCount                     int      `json:"media_count"`
	Lang                           any      `json:"lang"`
	ContributorsEnabled            bool     `json:"contributors_enabled"`
	IsTranslator                   bool     `json:"is_translator"`
	IsTranslationEnabled           bool     `json:"is_translation_enabled"`
	ProfileBackgroundColor         string   `json:"profile_background_color"`
	ProfileBackgroundImageURL      string   `json:"profile_background_image_url"`
	ProfileBackgroundImageURLHTTPS string   `json:"profile_background_image_url_https"`
	ProfileBackgroundTile          bool     `json:"profile_background_tile"`
	ProfileImageURL                string   `json:"profile_image_url"`
	ProfileImageURLHTTPS           string   `json:"profile_image_url_https"`
	ProfileBannerURL               string   `json:"profile_banner_url"`
	ProfileLinkColor               string   `json:"profile_link_color"`
	ProfileSidebarBorderColor      string   `json:"profile_sidebar_border_color"`
	ProfileSidebarFillColor        string   `json:"profile_sidebar_fill_color"`
	ProfileTextColor               string   `json:"profile_text_color"`
	ProfileUseBackgroundImage      bool     `json:"profile_use_background_image"`
	HasExtendedProfile             bool     `json:"has_extended_profile"`
	DefaultProfile                 bool     `json:"default_profile"`
	DefaultProfileImage            bool     `json:"default_profile_image"`
	PinnedTweetIds                 []int64  `json:"pinned_tweet_ids"`
	PinnedTweetIdsStr              []string `json:"pinned_tweet_ids_str"`
	HasCustomTimelines             bool     `json:"has_custom_timelines"`
	CanDm                          any      `json:"can_dm"`
	CanMediaTag                    bool     `json:"can_media_tag"`
	Following                      bool     `json:"following"`
	FollowRequestSent              bool     `json:"follow_request_sent"`
	Notifications                  bool     `json:"notifications"`
	Muting                         bool     `json:"muting"`
	Blocking                       bool     `json:"blocking"`
	BlockedBy                      bool     `json:"blocked_by"`
	WantRetweets                   bool     `json:"want_retweets"`
	AdvertiserAccountType          string   `json:"advertiser_account_type"`
	AdvertiserAccountServiceLevels []any    `json:"advertiser_account_service_levels"`
	ProfileInterstitialType        string   `json:"profile_interstitial_type"`
	BusinessProfileState           string   `json:"business_profile_state"`
	TranslatorType                 string   `json:"translator_type"`
	WithheldInCountries            []any    `json:"withheld_in_countries"`
	FollowedBy                     bool     `json:"followed_by"`
	ExtIsBlueVerified              bool     `json:"ext_is_blue_verified"`
	ExtHasNftAvatar                bool     `json:"ext_has_nft_avatar"`
	ExtProfileImageShape           string   `json:"ext_profile_image_shape"`
	RequireSomeConsent             bool     `json:"require_some_consent"`
}

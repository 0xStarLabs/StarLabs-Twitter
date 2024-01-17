package additional_twitter_methods

import (
	"compress/gzip"
	"encoding/json"
	"errors"
	"fmt"
	http "github.com/Danny-Dasilva/fhttp"
	"io"
	"strings"
	"twitter/extra"
	"twitter/utils"
)

// GetTwitterUsername get and return Twitter account's username
func GetTwitterUsername(index int, httpClient *http.Client, cookieClient *utils.CookieClient, bearerToken string, csrfToken string) (string, string, error) {
	for i := 0; i < 1; i++ {
		baseURL := "https://api.twitter.com/1.1/account/settings.json?include_mention_filter=true&include_nsfw_user_flag=true&include_nsfw_admin_flag=true&include_ranked_timeline=true&include_alt_text_compose=true&ext=ssoConnections&include_country_code=true&include_ext_dm_nsfw_media_filter=true&include_ext_sharing_audiospaces_listening_data_with_followers=true"

		// Create new request
		req, err := http.NewRequest("GET", baseURL, nil)
		if err != nil {
			extra.Logger{}.Warning("%d | Failed to build activate request: %s", index, err.Error())
			continue
		}
		req.Header = http.Header{
			"accept":                {"*/*"},
			"accept-encoding":       {"gzip, deflate, br"},
			"authorization":         {bearerToken},
			"cookie":                {cookieClient.CookiesToHeader()},
			"origin":                {"https://twitter.com"},
			"referer":               {"https://twitter.com/"},
			"sec-ch-ua-mobile":      {"?0"},
			"sec-ch-ua-platform":    {`"Windows"`},
			"sec-fetch-dest":        {"empty"},
			"sec-fetch-mode":        {"cors"},
			"sec-fetch-site":        {"same-site"},
			"x-csrf-token":          {csrfToken},
			"x-twitter-active-user": {"no"},
			// x-client-transaction-id
			http.HeaderOrderKey: {
				"accept",
				"accept-encoding",
				"authorization",
				"cookie",
				"origin",
				"referer",
				"sec-ch-ua",
				"sec-ch-ua-mobile",
				"sec-ch-ua-platform",
				"sec-fetch-dest",
				"sec-fetch-mode",
				"sec-fetch-site",
				"user-agent",
				"x-csrf-token",
				"x-twitter-active-user",
			},
			http.PHeaderOrderKey: {":authority", ":method", ":path", ":scheme"},
		}

		resp, err := httpClient.Do(req)
		if err != nil {
			extra.Logger{}.Warning("%d | Failed to do get username request: %s", index, err.Error())
			continue
		}
		defer resp.Body.Close()

		cookieClient.SetCookieFromResponse(resp)

		csrfToken, ok := cookieClient.GetCookieValue("ct0")
		if ok != true {
			extra.Logger{}.Error("%d | Failed to get new csrf token", index)
			continue
		}

		var reader io.ReadCloser
		switch resp.Header.Get("Content-Encoding") {
		case "gzip":
			reader, err = gzip.NewReader(resp.Body)
			if err != nil {
				extra.Logger{}.Error("%d | Failed to create gzip reader while getting username: %s", index, err.Error())
				continue
			}
			defer reader.Close()
		default:
			reader = resp.Body
		}

		bodyBytes, err := io.ReadAll(reader)
		if err != nil {
			extra.Logger{}.Error("%d | Failed to read response body: %s", index, err.Error())
			continue
		}

		bodyString := string(bodyBytes)

		if strings.Contains(bodyString, "screen_name") {
			var responseDataUsername getUsernameJSON
			err = json.Unmarshal(bodyBytes, &responseDataUsername)
			if err != nil {
				extra.Logger{}.Error("%d | Failed to do unmarshal in get username response: %s", index, err.Error())
				continue
			}
			return responseDataUsername.ScreenName, csrfToken, nil

		} else if strings.Contains(bodyString, "this account is temporarily locked") {
			extra.Logger{}.Error("%d | Account is temporarily locked", index)
			return "locked", csrfToken, errors.New("")

		} else if strings.Contains(bodyString, "Could not authenticate you") {
			extra.Logger{}.Error("%d | Could not authenticate you", index)
			return "failed_auth", csrfToken, errors.New("")
		} else {
			fmt.Println(bodyString)
		}
	}

	extra.Logger{}.Error("%d | Unable to get twitter username", index)
	return "", "", errors.New("unable to get twitter username")
}

type getUsernameJSON struct {
	Protected                bool   `json:"protected"`
	ScreenName               string `json:"screen_name"`
	AlwaysUseHTTPS           bool   `json:"always_use_https"`
	UseCookiePersonalization bool   `json:"use_cookie_personalization"`
	SleepTime                struct {
		Enabled   bool `json:"enabled"`
		EndTime   any  `json:"end_time"`
		StartTime any  `json:"start_time"`
	} `json:"sleep_time"`
	GeoEnabled                                      bool   `json:"geo_enabled"`
	Language                                        string `json:"language"`
	DiscoverableByEmail                             bool   `json:"discoverable_by_email"`
	DiscoverableByMobilePhone                       bool   `json:"discoverable_by_mobile_phone"`
	DisplaySensitiveMedia                           bool   `json:"display_sensitive_media"`
	PersonalizedTrends                              bool   `json:"personalized_trends"`
	AllowMediaTagging                               string `json:"allow_media_tagging"`
	AllowContributorRequest                         string `json:"allow_contributor_request"`
	AllowAdsPersonalization                         bool   `json:"allow_ads_personalization"`
	AllowLoggedOutDevicePersonalization             bool   `json:"allow_logged_out_device_personalization"`
	AllowLocationHistoryPersonalization             bool   `json:"allow_location_history_personalization"`
	AllowSharingDataForThirdPartyPersonalization    bool   `json:"allow_sharing_data_for_third_party_personalization"`
	AllowDmsFrom                                    string `json:"allow_dms_from"`
	AlwaysAllowDmsFromSubscribers                   any    `json:"always_allow_dms_from_subscribers"`
	AllowDmGroupsFrom                               string `json:"allow_dm_groups_from"`
	TranslatorType                                  string `json:"translator_type"`
	CountryCode                                     string `json:"country_code"`
	NsfwUser                                        bool   `json:"nsfw_user"`
	NsfwAdmin                                       bool   `json:"nsfw_admin"`
	RankedTimelineSetting                           any    `json:"ranked_timeline_setting"`
	RankedTimelineEligible                          any    `json:"ranked_timeline_eligible"`
	AddressBookLiveSyncEnabled                      bool   `json:"address_book_live_sync_enabled"`
	UniversalQualityFilteringEnabled                string `json:"universal_quality_filtering_enabled"`
	DmReceiptSetting                                string `json:"dm_receipt_setting"`
	AltTextComposeEnabled                           any    `json:"alt_text_compose_enabled"`
	MentionFilter                                   string `json:"mention_filter"`
	AllowAuthenticatedPeriscopeRequests             bool   `json:"allow_authenticated_periscope_requests"`
	ProtectPasswordReset                            bool   `json:"protect_password_reset"`
	RequirePasswordLogin                            bool   `json:"require_password_login"`
	RequiresLoginVerification                       bool   `json:"requires_login_verification"`
	ExtSharingAudiospacesListeningDataWithFollowers bool   `json:"ext_sharing_audiospaces_listening_data_with_followers"`
	Ext                                             struct {
		SsoConnections struct {
			R struct {
				Ok []any `json:"ok"`
			} `json:"r"`
			TTL int `json:"ttl"`
		} `json:"ssoConnections"`
	} `json:"ext"`
	DmQualityFilter  string `json:"dm_quality_filter"`
	AutoplayDisabled bool   `json:"autoplay_disabled"`
	SettingsMetadata struct {
	} `json:"settings_metadata"`
}

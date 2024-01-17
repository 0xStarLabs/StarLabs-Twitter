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

// ChangeUsername change account username
func (twitter *Twitter) ChangeUsername(newUsername string) bool {
	changeUsernameURL := fmt.Sprintf("https://twitter.com/i/api/1.1/account/settings.json")

	for i := 0; i < twitter.config.Info.MaxTasksRetries; i++ {
		var stringData = fmt.Sprintf(`screen_name=%s`, newUsername)
		data := strings.NewReader(stringData)

		// Create new request
		req, err := http.NewRequest("POST", changeUsernameURL, data)
		if err != nil {
			extra.Logger{}.Error("Failed to build change username request: %s", err.Error())
			continue
		}
		req.Header = http.Header{
			"accept":                {"*/*"},
			"accept-encoding":       {"gzip, deflate, br"},
			"authorization":         {twitter.queryID.BearerToken},
			"content-type":          {"application/x-www-form-urlencoded"},
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
			extra.Logger{}.Error("Failed to do change username request: %s", err.Error())
			continue
		}
		defer resp.Body.Close()

		var reader io.ReadCloser
		switch resp.Header.Get("Content-Encoding") {
		case "gzip":
			reader, err = gzip.NewReader(resp.Body)
			if err != nil {
				extra.Logger{}.Error("Failed to create gzip reader while change username: %s", err.Error())
				continue
			}
			defer reader.Close()
		default:
			reader = resp.Body
		}

		bodyBytes, err := io.ReadAll(reader)
		if err != nil {
			extra.Logger{}.Error("Failed to read change username response body: %s", err.Error())
			continue
		}

		bodyString := string(bodyBytes)

		if resp.StatusCode >= 200 && resp.StatusCode <= 299 {
			if strings.Contains(bodyString, "screen_name") {
				var changeUsernameResponse ChangeUsernameResponse
				err = json.Unmarshal(bodyBytes, &changeUsernameResponse)
				if err != nil {
					extra.Logger{}.Error("%d | Failed to do unmarshal in change username response: %s", twitter.index, err.Error())
					continue
				}
				if changeUsernameResponse.ScreenName == newUsername {
					twitter.logger.Success("%d | %s | Successfully changed username", twitter.index, twitter.Username)
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
			extra.Logger{}.Error("%d | Unknown response while change username: %s", twitter.index, bodyString)
			continue
		}
	}

	extra.Logger{}.Error("%d | %s | Unable to change username", twitter.index, twitter.Username)
	return false
}

type ChangeUsernameResponse struct {
	Protected                bool   `json:"protected"`
	ScreenName               string `json:"screen_name"`
	AlwaysUseHTTPS           bool   `json:"always_use_https"`
	UseCookiePersonalization bool   `json:"use_cookie_personalization"`
	SleepTime                struct {
		Enabled   bool `json:"enabled"`
		EndTime   any  `json:"end_time"`
		StartTime any  `json:"start_time"`
	} `json:"sleep_time"`
	GeoEnabled                                   bool   `json:"geo_enabled"`
	Language                                     string `json:"language"`
	DiscoverableByEmail                          bool   `json:"discoverable_by_email"`
	DiscoverableByMobilePhone                    bool   `json:"discoverable_by_mobile_phone"`
	DisplaySensitiveMedia                        bool   `json:"display_sensitive_media"`
	PersonalizedTrends                           bool   `json:"personalized_trends"`
	AllowMediaTagging                            string `json:"allow_media_tagging"`
	AllowContributorRequest                      string `json:"allow_contributor_request"`
	AllowAdsPersonalization                      bool   `json:"allow_ads_personalization"`
	AllowLoggedOutDevicePersonalization          bool   `json:"allow_logged_out_device_personalization"`
	AllowLocationHistoryPersonalization          bool   `json:"allow_location_history_personalization"`
	AllowSharingDataForThirdPartyPersonalization bool   `json:"allow_sharing_data_for_third_party_personalization"`
	AllowDmsFrom                                 string `json:"allow_dms_from"`
	AlwaysAllowDmsFromSubscribers                any    `json:"always_allow_dms_from_subscribers"`
	AllowDmGroupsFrom                            string `json:"allow_dm_groups_from"`
	TranslatorType                               string `json:"translator_type"`
	CountryCode                                  string `json:"country_code"`
	AddressBookLiveSyncEnabled                   bool   `json:"address_book_live_sync_enabled"`
	UniversalQualityFilteringEnabled             string `json:"universal_quality_filtering_enabled"`
	DmReceiptSetting                             string `json:"dm_receipt_setting"`
	AllowAuthenticatedPeriscopeRequests          bool   `json:"allow_authenticated_periscope_requests"`
	ProtectPasswordReset                         bool   `json:"protect_password_reset"`
	RequirePasswordLogin                         bool   `json:"require_password_login"`
	RequiresLoginVerification                    bool   `json:"requires_login_verification"`
	DmQualityFilter                              string `json:"dm_quality_filter"`
	AutoplayDisabled                             bool   `json:"autoplay_disabled"`
	SettingsMetadata                             struct {
		IsEu string `json:"is_eu"`
	} `json:"settings_metadata"`
}

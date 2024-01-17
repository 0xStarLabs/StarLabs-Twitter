package instance

import (
	"compress/gzip"
	"fmt"
	http "github.com/Danny-Dasilva/fhttp"
	"io"
	"strings"
	"twitter/extra"
)

// ChangeBackground change account background
func (twitter *Twitter) ChangeBackground(pictureBase64Encoded string) bool {
	pictureURL := fmt.Sprintf("https://api.twitter.com/1.1/account/update_profile_banner.json")

	for i := 0; i < twitter.config.Info.MaxTasksRetries; i++ {
		mediaID := twitter.UploadImageToTwitter(pictureBase64Encoded)
		if mediaID == "" {
			continue
		}

		var stringData = fmt.Sprintf(`include_profile_interstitial_type=1&include_blocking=1&include_blocked_by=1&include_followed_by=1&include_want_retweets=1&include_mute_edge=1&include_can_dm=1&include_can_media_tag=1&include_ext_has_nft_avatar=1&include_ext_is_blue_verified=1&include_ext_verified_type=1&include_ext_profile_image_shape=1&skip_status=1&return_user=true&media_id=%s`, mediaID)
		data := strings.NewReader(stringData)

		// Create new request
		req, err := http.NewRequest("POST", pictureURL, data)
		if err != nil {
			twitter.logger.Error("%d | Failed to build change background request: %s", twitter.index, err.Error())
			continue
		}

		req.Header = http.Header{
			"accept":                {"*/*"},
			"accept-encoding":       {"gzip, deflate, br"},
			"authorization":         {twitter.queryID.BearerToken},
			"content-type":          {"application/x-www-form-urlencoded"},
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
			twitter.logger.Error("%d | Failed to do change background request: %s", twitter.index, err.Error())
			continue
		}
		defer resp.Body.Close()

		var reader io.ReadCloser
		switch resp.Header.Get("Content-Encoding") {
		case "gzip":
			reader, err = gzip.NewReader(resp.Body)
			if err != nil {
				twitter.logger.Error("%d | Failed to create gzip reader while change background: %s", twitter.index, err.Error())
				continue
			}
			defer reader.Close()
		default:
			reader = resp.Body
		}

		bodyBytes, err := io.ReadAll(reader)
		if err != nil {
			twitter.logger.Error("%d | Failed to read change background response body: %s", twitter.index, err.Error())
			continue
		}

		bodyString := string(bodyBytes)

		if resp.StatusCode >= 200 && resp.StatusCode <= 299 {
			if strings.Contains(bodyString, "screen_name") {
				twitter.logger.Success("%d | %s successfully changed background", twitter.index, twitter.Username)
				return true
			}

		} else if strings.Contains(bodyString, "this account is temporarily locked") {
			twitter.logger.Error("%d | Account is temporarily locked", twitter.index)
			return false

		} else if strings.Contains(bodyString, "Could not authenticate you") {
			twitter.logger.Error("%d | Could not authenticate you", twitter.index)
			return false
		} else {
			twitter.logger.Error("%d | Unknown response while change background: %s", twitter.index, bodyString)
			continue
		}
	}

	extra.Logger{}.Error("%d | Unable to change background", twitter.index)
	return false
}

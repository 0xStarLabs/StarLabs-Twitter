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

// ChangeBirthdate change account birthdate
func (twitter *Twitter) ChangeBirthdate(newBirthdate string) bool {
	changeBirthdateURL := fmt.Sprintf("https://twitter.com/i/api/1.1/account/update_profile.json")

	dateParts := strings.Split(newBirthdate, ":")

	for i := 0; i < twitter.config.Info.MaxTasksRetries; i++ {
		var stringData = fmt.Sprintf(`birthdate_day=%s&birthdate_month=%s&birthdate_year=%s&birthdate_visibility=mutualfollow&birthdate_year_visibility=self`, dateParts[0], dateParts[1], dateParts[2])
		data := strings.NewReader(stringData)

		// Create new request
		req, err := http.NewRequest("POST", changeBirthdateURL, data)
		if err != nil {
			extra.Logger{}.Error("%d | Failed to build change birthdate request: %s", twitter.index, err.Error())
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
			extra.Logger{}.Error("%d | Failed to do change birthdate request: %s", twitter.index, err.Error())
			continue
		}
		defer resp.Body.Close()

		var reader io.ReadCloser
		switch resp.Header.Get("Content-Encoding") {
		case "gzip":
			reader, err = gzip.NewReader(resp.Body)
			if err != nil {
				extra.Logger{}.Error("%d | Failed to create gzip reader while change birthdate: %s", twitter.index, err.Error())
				continue
			}
			defer reader.Close()
		default:
			reader = resp.Body
		}

		bodyBytes, err := io.ReadAll(reader)
		if err != nil {
			extra.Logger{}.Error("%d | Failed to read change birthdate response body: %s", twitter.index, err.Error())
			continue
		}

		bodyString := string(bodyBytes)

		if resp.StatusCode >= 200 && resp.StatusCode <= 299 {
			if strings.Contains(bodyString, "screen_name") {
				var changeUsernameResponse ChangeUsernameResponse
				err = json.Unmarshal(bodyBytes, &changeUsernameResponse)
				if err != nil {
					extra.Logger{}.Error("%d | Failed to do unmarshal in change birthdate response: %s", twitter.index, err.Error())
					continue
				}
				twitter.logger.Success("%d | %s | Successfully changed birthdate", twitter.index, twitter.Username)
				return true
			}

		} else if strings.Contains(bodyString, "this account is temporarily locked") {
			extra.Logger{}.Error("%d | Account is temporarily locked", twitter.index)
			return false

		} else if strings.Contains(bodyString, "Could not authenticate you") {
			extra.Logger{}.Error("%d | Could not authenticate you", twitter.index)
			return false
		} else {
			extra.Logger{}.Error("%d | Unknown response while change birthdate: %s", twitter.index, bodyString)
			continue
		}
	}

	extra.Logger{}.Error("%d | %s | Unable to change username", twitter.index, twitter.Username)
	return false
}

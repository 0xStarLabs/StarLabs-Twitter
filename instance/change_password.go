package instance

import (
	"compress/gzip"
	"encoding/json"
	"fmt"
	"github.com/Danny-Dasilva/CycleTLS/cycletls"
	http "github.com/Danny-Dasilva/fhttp"
	"io"
	"strings"
	"twitter/extra"
)

// ChangePassword change account password
// this function returns new auth token if action was successful
// returns empty string "" if failed to change a password
func (twitter *Twitter) ChangePassword(oldPassword string, newPassword string) string {
	var newAuthToken string

	changeLocationURL := fmt.Sprintf("https://twitter.com/i/api/i/account/change_password.json")

	for i := 0; i < twitter.config.Info.MaxTasksRetries; i++ {
		var stringData = fmt.Sprintf(`current_password=%s&password=%s&password_confirmation=%s`, oldPassword, newPassword, newPassword)
		data := strings.NewReader(stringData)

		// Create new request
		req, err := http.NewRequest("POST", changeLocationURL, data)
		if err != nil {
			extra.Logger{}.Error("%d | Failed to build change password request: %s", err.Error())
			continue
		}
		req.Header = http.Header{
			"accept":                {"*/*"},
			"accept-encoding":       {"gzip, deflate, br"},
			"authorization":         {twitter.queryID.BearerToken},
			"content-type":          {"application/x-www-form-urlencoded"},
			"cookie":                {twitter.cookies.CookiesToHeader()},
			"origin":                {"https://twitter.com"},
			"referer":               {"https://twitter.com/settings/password"},
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
			extra.Logger{}.Error("%d | Failed to do change password request: %s", err.Error())
			continue
		}
		defer resp.Body.Close()

		var reader io.ReadCloser
		switch resp.Header.Get("Content-Encoding") {
		case "gzip":
			reader, err = gzip.NewReader(resp.Body)
			if err != nil {
				extra.Logger{}.Error("%d | Failed to create gzip reader while change password: %s", err.Error())
				continue
			}
			defer reader.Close()
		default:
			reader = resp.Body
		}

		bodyBytes, err := io.ReadAll(reader)
		if err != nil {
			extra.Logger{}.Error("%d | Failed to read change password response body: %s", err.Error())
			continue
		}

		bodyString := string(bodyBytes)

		for _, httpCookie := range resp.Cookies() {
			newCookie := cycletls.Cookie{
				Name:  httpCookie.Name,
				Value: httpCookie.Value,
			}
			if newCookie.Name == "auth_token" {
				newAuthToken = newCookie.Value
			}
		}

		if resp.StatusCode >= 200 && resp.StatusCode <= 299 {
			if strings.Contains(bodyString, "status") {
				var changePasswordResponse ChangePasswordResponse
				err = json.Unmarshal(bodyBytes, &changePasswordResponse)
				if err != nil {
					extra.Logger{}.Error("%d | Failed to do unmarshal in change password response: %s", twitter.index, err.Error())
					continue
				}
				if changePasswordResponse.Status == "ok" {
					twitter.logger.Success("%d | %s | Successfully changed password", twitter.index, twitter.Username)
					return newAuthToken
				}
			}

		} else if strings.Contains(bodyString, "this account is temporarily locked") {
			extra.Logger{}.Error("%d | Account is temporarily locked", twitter.index)
			return newAuthToken

		} else if strings.Contains(bodyString, "Could not authenticate you") {
			extra.Logger{}.Error("%d | Could not authenticate you", twitter.index)
			return newAuthToken
		} else {
			extra.Logger{}.Error("%d | Unknown response while change password: %s", twitter.index, bodyString)
			continue
		}
	}

	extra.Logger{}.Error("%d | %s | Unable to change password", twitter.index, twitter.Username)
	return newAuthToken
}

type ChangePasswordResponse struct {
	Status string `json:"status"`
}

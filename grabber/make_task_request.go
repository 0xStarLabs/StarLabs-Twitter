package grabber

import (
	"compress/gzip"
	"encoding/json"
	"fmt"
	http "github.com/Danny-Dasilva/fhttp"
	"io"
	"strings"
	"twitter/extra"
)

func (grabber *Grabber) MakeTaskRequest(params string, jsonData string) bool {
	var taskURL string
	if params != "" {
		taskURL = fmt.Sprintf("https://api.twitter.com/1.1/onboarding/task.json?%s", params)
	} else {
		taskURL = fmt.Sprintf("https://api.twitter.com/1.1/onboarding/task.json")
	}

	data := strings.NewReader(jsonData)

	for i := 0; i < 1; i++ {
		// Create new request
		req, err := http.NewRequest("POST", taskURL, data)
		if err != nil {
			grabber.logger.Error("Failed to build task request: %s", err.Error())
			continue
		}
		req.Header = http.Header{
			"accept":                    {"text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7"},
			"accept-encoding":           {"gzip, deflate, br"},
			"authorization":             {grabber.queryID.BearerToken},
			"content-type":              {"application/json"},
			"cookie":                    {grabber.cookies.CookiesToHeader()},
			"origin":                    {"https://twitter.com"},
			"referer":                   {"https://twitter.com/"},
			"sec-ch-ua-mobile":          {"?0"},
			"sec-ch-ua-platform":        {`"Windows"`},
			"upgrade-insecure-requests": {"1"},
			"x-guest-token":             {grabber.guestToken},
			http.HeaderOrderKey: {
				"accept",
				"accept-encoding",
				"content-type",
				"origin",
				"referer",
				"sec-ch-ua-mobile",
				"sec-ch-ua-platform",
				"upgrade-insecure-requests",
				"user-agent",
			},
			//http.PHeaderOrderKey: {":authority", ":method", ":path", ":scheme"},
		}

		resp, err := grabber.client.Do(req)
		if err != nil {
			grabber.logger.Error("Failed to do task request: %s", err.Error())
			continue
		}
		grabber.cookies.SetCookieFromResponse(resp)

		defer resp.Body.Close()

		var reader io.ReadCloser
		switch resp.Header.Get("Content-Encoding") {
		case "gzip":
			reader, err = gzip.NewReader(resp.Body)
			if err != nil {
				extra.Logger{}.Error("Failed to create gzip reader while task request: %s", err.Error())
				continue
			}
			defer reader.Close()
		default:
			reader = resp.Body
		}

		bodyBytes, err := io.ReadAll(reader)
		if err != nil {
			extra.Logger{}.Error("Failed to read task response body: %s", err.Error())
			continue
		}

		bodyString := string(bodyBytes)

		if strings.Contains(bodyString, "flow_token") && resp.StatusCode >= 200 && resp.StatusCode <= 299 {
			var task MakeTaskResponse
			err = json.Unmarshal(bodyBytes, &task)
			if err != nil {
				extra.Logger{}.Error("%d | Failed to do unmarshal in make task request: %s", grabber.index, err.Error())
				continue
			}
			if strings.Contains(bodyString, "LoginTwoFactorAuthChallenge") || strings.Contains(bodyString, "LoginTwoFactorAuthChooseMethod") {
				grabber.logger.Warning("%d | Account requires two-factor authentication.")
				return false
			}
			if strings.Contains(bodyString, "AccountDuplicationCheck") || strings.Contains(bodyString, "AccountDuplicationCheck_true") {
				grabber.logger.Warning("%d | Account requires duplication check.", grabber.index)
				return false
			}
			grabber.flowToken = task.FlowToken
			return true

		} else if strings.Contains(bodyString, "we've sent a confirmation code to") {
			grabber.logger.Warning("%d | Twitter sent a confirmation code to the mail.")
			return false
		}
	}
	grabber.logger.Error("%d | Failed to task request", grabber.index)
	return false
}

type MakeTaskResponse struct {
	FlowToken string `json:"flow_token"`
	Status    string `json:"status"`
	Subtasks  []struct {
		SubtaskID         string `json:"subtask_id"`
		JsInstrumentation struct {
			URL       string `json:"url"`
			TimeoutMs int    `json:"timeout_ms"`
			NextLink  struct {
				LinkType string `json:"link_type"`
				LinkID   string `json:"link_id"`
			} `json:"next_link"`
		} `json:"js_instrumentation"`
	} `json:"subtasks"`
}

package grabber

import (
	"compress/gzip"
	"encoding/json"
	"github.com/Danny-Dasilva/CycleTLS/cycletls"
	http "github.com/Danny-Dasilva/fhttp"
	"io"
)

func (grabber *Grabber) Flow() bool {
	ok := grabber.initRequest()
	if ok == false {
		return false
	}
	ok = grabber.activateGuestToken()
	if ok == false {
		return false
	}
	ok = grabber.MakeTaskRequest("flow_name=login", GenerateFirstTaskString())
	if ok == false {
		return false
	}
	ok = grabber.MakeTaskRequest("", GenerateSecondTaskString(grabber.flowToken))
	if ok == false {
		return false
	}
	ok = grabber.MakeTaskRequest("", GenerateThirdTaskString(grabber.flowToken, grabber.login))
	if ok == false {
		return false
	}
	ok = grabber.MakeTaskRequest("", GenerateFourthTaskString(grabber.flowToken, grabber.password))
	if ok == false {
		return false
	}
	ok = grabber.MakeTaskRequest("", GenerateFifthTaskString(grabber.flowToken))
	if ok == false {
		return false
	}
	return true
}

func (grabber *Grabber) initRequest() bool {
	loginURL := "https://twitter.com/i/flow/login"

	for i := 0; i < 1; i++ {
		// Create new request
		req, err := http.NewRequest("GET", loginURL, nil)
		if err != nil {
			grabber.logger.Error("Failed to build init grabber request: %s", err.Error())
			continue
		}
		req.Header = http.Header{
			"accept":             {"text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7"},
			"accept-encoding":    {"gzip, deflate, br"},
			"sec-ch-ua-mobile":   {"?0"},
			"sec-ch-ua-platform": {`"Windows"`},
			http.HeaderOrderKey: {
				"accept",
				"accept-encoding",
				"sec-ch-ua-mobile",
				"sec-ch-ua-platform",
			},
			http.PHeaderOrderKey: {":authority", ":method", ":path", ":scheme"},
		}

		resp, err := grabber.client.Do(req)
		if err != nil {
			grabber.logger.Error("%d | Failed to do init grabber request: %s", grabber.index, err.Error())
			continue
		}
		grabber.cookies.SetCookieFromResponse(resp)
		return true
	}
	grabber.logger.Error("%d | Failed to init grabber request", grabber.index)
	return false
}

func (grabber *Grabber) activateGuestToken() bool {
	activateURL := "https://api.twitter.com/1.1/guest/activate.json"

	for i := 0; i < 1; i++ {
		// Create new request
		req, err := http.NewRequest("POST", activateURL, nil)
		if err != nil {
			grabber.logger.Error("%d | Failed to build guest token grabber request: %s", grabber.index, err.Error())
			continue
		}
		req.Header = http.Header{
			"accept":                  {"text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7"},
			"accept-encoding":         {"gzip, deflate, br"},
			"authorization":           {grabber.queryID.BearerToken},
			"content-type":            {"application/x-www-form-urlencoded"},
			"origin":                  {"https://twitter.com"},
			"referer":                 {"https://twitter.com/"},
			"sec-ch-ua-mobile":        {"?0"},
			"sec-ch-ua-platform":      {`"Windows"`},
			"x-twitter-activate-user": {"yes"},
			http.HeaderOrderKey:       {},
			http.PHeaderOrderKey:      {":authority", ":method", ":path", ":scheme"},
		}
		resp, err := grabber.client.Do(req)
		if err != nil {
			grabber.logger.Error("%d | Failed to do guest token grabber request: %s", grabber.index, err.Error())
			continue
		}
		grabber.cookies.SetCookieFromResponse(resp)

		defer resp.Body.Close()

		var reader io.ReadCloser
		switch resp.Header.Get("Content-Encoding") {
		case "gzip":
			reader, err = gzip.NewReader(resp.Body)
			if err != nil {
				grabber.logger.Error("%d | Failed to create gzip reader while guest token grabber: %s", grabber.index, err.Error())
				continue
			}
			defer reader.Close()
		default:
			reader = resp.Body
		}

		bodyBytes, err := io.ReadAll(reader)
		if err != nil {
			grabber.logger.Error("%d | Failed to read check guest token response body: %s", grabber.index, err.Error())
			continue
		}

		var guestToken GuestTokenResponse
		err = json.Unmarshal(bodyBytes, &guestToken)
		if err != nil {
			grabber.logger.Error("%d | Failed to do unmarshal in guest token response: %s", grabber.index, err.Error())
			continue
		}
		grabber.guestToken = guestToken.GuestToken
		gtCookie := cycletls.Cookie{
			Name:  "gt",
			Value: grabber.guestToken,
		}
		grabber.cookies.Cookies = append(grabber.cookies.Cookies, gtCookie)
		return true
	}

	grabber.logger.Error("%d | Failed to get guest token", grabber.index)
	return false
}

type GuestTokenResponse struct {
	GuestToken string `json:"guest_token"`
}

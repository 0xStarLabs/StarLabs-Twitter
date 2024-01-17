package instance

import (
	"compress/gzip"
	"fmt"
	api2captcha "github.com/2captcha/2captcha-go"
	http "github.com/Danny-Dasilva/fhttp"
	"io"
	"strings"
	"time"
	"twitter/extra"
)

func (twitter *Twitter) Unfreeze() bool {
	for i := 0; i < twitter.config.Info.MaxTasksRetries; i++ {
		req, err := http.NewRequest("GET", "https://twitter.com/account/access", nil)
		if err != nil {
			twitter.logger.Error("%d | Failed to build unfreeze get request: %s", twitter.index, err.Error())
			continue
		}
		req.Header = http.Header{
			"accept":                    {"text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7"},
			"accept-encoding":           {"gzip, deflate, br"},
			"cookie":                    {twitter.cookies.CookiesToHeader()},
			"sec-ch-ua-mobile":          {"?0"},
			"sec-ch-ua-platform":        {`"Windows"`},
			"sec-fetch-dest":            {"document"},
			"sec-fetch-mode":            {"navigate"},
			"sec-fetch-site":            {"same-origin"},
			"upgrade-insecure-requests": {"1"},
			http.HeaderOrderKey: {
				"accept",
				"accept-encoding",
				"cookie",
				"sec-ch-ua-mobile",
				"sec-ch-ua-platform",
				"sec-fetch-dest",
				"sec-fetch-mode",
				"sec-fetch-site",
				"user-agent",
				"upgrade-insecure-requests",
			},
			http.PHeaderOrderKey: {":authority", ":method", ":path", ":scheme"},
		}

		resp, err := twitter.client.Do(req)
		if err != nil {
			extra.Logger{}.Error("%d | Failed to do unfreeze get request: %s", twitter.index, err.Error())
			continue
		}
		defer resp.Body.Close()
		twitter.cookies.SetCookieFromResponse(resp)
		var reader io.ReadCloser
		switch resp.Header.Get("Content-Encoding") {
		case "gzip":
			reader, err = gzip.NewReader(resp.Body)
			if err != nil {
				extra.Logger{}.Error("%d | Failed to create gzip reader while unfreeze get: %s", twitter.index, err.Error())
				continue
			}
			defer reader.Close()
		default:
			reader = resp.Body
		}

		bodyBytes, err := io.ReadAll(reader)
		if err != nil {
			extra.Logger{}.Error("%d | Failed to read unfreeze GET response body: %s", twitter.index, err.Error())
			continue
		}

		bodyString := string(bodyBytes)

		if strings.Contains(bodyString, "Verify email") || strings.Contains(bodyString, "Verify your email address") {
			twitter.logger.Warning("%d | You need to verify an EMAIL to unlock your account.", twitter.index)
			return false
		} else if strings.Contains(bodyString, "change your password") {
			twitter.logger.Warning("%d | You need to change your password to unlock your account.", twitter.index)
			return false
		} else if strings.Contains(bodyString, "Pass an Arkose challenge") || strings.Contains(bodyString, "arkose_iframe") {
			twitter.logger.Info("%d | Starting to pass Arkose captcha", twitter.index)

		} else {
			twitter.logger.Warning("%d | Unknown server response: %s", twitter.index, bodyString)
			return false
		}

		var authenticityToken string
		var assignmentToken string

		if strings.Contains(bodyString, "authenticity_token") && strings.Contains(bodyString, "assignment_token") {
			parts := strings.Split(bodyString, `name="authenticity_token" value="`)
			if len(parts) > 1 {
				authenticityToken = strings.Split(parts[1], `"`)[0]
			} else {
				twitter.logger.Error("%d | Failed to get authenticity token while unfreeze", twitter.index)
				return false
			}
			parts = strings.Split(bodyString, `assignment_token" value="`)
			if len(parts) > 1 {
				assignmentToken = strings.Split(parts[1], `"`)[0]
				assignmentToken = strings.Replace(assignmentToken, "-", "", 1)
			} else {
				twitter.logger.Error("%d | Failed to get assignment token while unfreeze", twitter.index)
				return false
			}
		}
		fmt.Println(authenticityToken)
		fmt.Println(assignmentToken)
		solution := twitter.SolveCaptcha()

		if solution == "" {
			return false
		}
		tokenID := strings.Split(solution, "|")[0]
		tokenID = strings.Replace(tokenID, "token:", "", 1)
		fmt.Println("solution token is this: ", tokenID)

		//
		//twitter.BypassFuncaptcha(solution)
		//fmt.Println(1)
		//return true
		dataStr := fmt.Sprintf(`authenticity_token=%s&assignment_token=%s&lang=en&flow=`, authenticityToken, assignmentToken)
		var data = strings.NewReader(dataStr)

		req, err = http.NewRequest("POST", "https://twitter.com/account/access", data)
		if err != nil {
			twitter.logger.Error("%d | Failed to build captcha result request: %s", twitter.index, err)
			return false
		}
		req.Header = http.Header{
			"accept":                    {"text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7"},
			"content-type":              {"application/x-www-form-urlencoded"},
			"cookie":                    {twitter.cookies.CookiesToHeader()},
			"origin":                    {"https://twitter.com"},
			"pragma":                    {"no-cache"},
			"referer":                   {"https://twitter.com/account/access"},
			"sec-ch-ua-mobile":          {"?0"},
			"sec-ch-ua-platform":        {`"Windows"`},
			"sec-fetch-dest":            {"document"},
			"sec-fetch-mode":            {"navigate"},
			"sec-fetch-site":            {"same-origin"},
			"upgrade-insecure-requests": {"1"},
			http.HeaderOrderKey: {
				"accept",
				"content-type",
				"cookie",
				"origin",
				"pragma",
				"referer",
				"sec-ch-ua-mobile",
				"sec-ch-ua-platform",
				"sec-fetch-dest",
				"sec-fetch-mode",
				"sec-fetch-site",
				"user-agent",
				"upgrade-insecure-requests",
			},
			http.PHeaderOrderKey: {":authority", ":method", ":path", ":scheme"},
		}
		resp, err = twitter.client.Do(req)
		if err != nil {
			twitter.logger.Error("%d | Failed to send captcha result request: %s", twitter.index, err)
			return false
		}
		defer resp.Body.Close()

		twitter.cookies.SetCookieFromResponse(resp)
		fmt.Println(resp.StatusCode)

		URL := fmt.Sprintf("https://client-api.arkoselabs.com/fc/gc/?token=%s", tokenID)
		req, err = http.NewRequest("GET", URL, nil)
		if err != nil {
			twitter.logger.Error("%d | Failed to build captcha result request: %s", twitter.index, err)
			return false
		}
		req.Header = http.Header{
			"accept":                    {"application/json, text/plain, */*"},
			"cookie":                    {twitter.cookies.CookiesToHeader()},
			"pragma":                    {"no-cache"},
			"referer":                   {"https://twitter.com/account/access"},
			"sec-ch-ua-mobile":          {"?0"},
			"sec-ch-ua-platform":        {`"Windows"`},
			"sec-fetch-dest":            {"empty"},
			"sec-fetch-mode":            {"cors"},
			"sec-fetch-site":            {"same-origin"},
			"upgrade-insecure-requests": {"1"},
			http.HeaderOrderKey: {
				"accept",
				"cookie",
				"pragma",
				"referer",
				"sec-ch-ua-mobile",
				"sec-ch-ua-platform",
				"sec-fetch-dest",
				"sec-fetch-mode",
				"sec-fetch-site",
				"user-agent",
				"upgrade-insecure-requests",
			},
			http.PHeaderOrderKey: {":authority", ":method", ":path", ":scheme"},
		}
		resp, err = twitter.client.Do(req)
		if err != nil {
			twitter.logger.Error("%d | Failed to send captcha result request: %s", twitter.index, err)
			return false
		}
		defer resp.Body.Close()
		twitter.cookies.SetCookieFromResponse(resp)
		bodyBytes, err = io.ReadAll(reader)
		if err != nil {
			extra.Logger{}.Error("%d | Failed to read unfreeze GET response body: %s", twitter.index, err.Error())
			continue
		}

		fmt.Println("Ниже статус код логина")
		fmt.Println(resp.StatusCode)

		dataStr = fmt.Sprintf(`token=%s&sid=eu-west-1&render_type=canvas&lang=en&isAudioGame=false&analytics_tier=40&apiBreakerVersion=green`, tokenID)
		data = strings.NewReader(dataStr)

		req, err = http.NewRequest("POST", "https://client-api.arkoselabs.com/fc/gfct/", data)
		if err != nil {
			twitter.logger.Error("%d | Failed to build captcha result request: %s", twitter.index, err)
			return false
		}
		req.Header = http.Header{
			"accept":               {"*/*"},
			"content-type":         {"application/x-www-form-urlencoded; charset=UTF-8"},
			"cookie":               {twitter.cookies.CookiesToHeader()},
			"origin":               {"https://client-api.arkoselabs.com"},
			"pragma":               {"no-cache"},
			"referer":              {fmt.Sprintf(`https://client-api.arkoselabs.com/fc/assets/ec-game-core/game-core/1.17.1/standard/index.html?session=%s&r=eu-west-1&meta=3&meta_width=558&meta_height=523&metabgclr=transparent&metaiconclr=%23555555&guitextcolor=%23000000&lang=en&pk=0152B4EB-D2DC-460A-89A1-629838B529C9&at=40&ag=101&cdn_url=https%3A%2F%2Fclient-api.arkoselabs.com%2Fcdn%2Ffc&lurl=https%3A%2F%2Faudio-eu-west-1.arkoselabs.com&surl=https%3A%2F%2Fclient-api.arkoselabs.com&smurl=https%3A%2F%2Fclient-api.arkoselabs.com%2Fcdn%2Ffc%2Fassets%2Fstyle-manager&theme=default`), tokenID},
			"sec-ch-ua-mobile":     {"?0"},
			"sec-ch-ua-platform":   {`"Windows"`},
			"sec-fetch-dest":       {"empty"},
			"sec-fetch-mode":       {"cors"},
			"sec-fetch-site":       {"same-origin"},
			"x-newrelic-timestamp": {time.Now().Format(time.RFC3339Nano)},
			"x-requested-with":     {"XMLHttpRequest"},
			http.HeaderOrderKey: {
				"accept",
				"content-type",
				"cookie",
				"origin",
				"pragma",
				"referer",
				"sec-ch-ua-mobile",
				"sec-ch-ua-platform",
				"sec-fetch-dest",
				"sec-fetch-mode",
				"sec-fetch-site",
				"user-agent",
				"x-newrelic-timestamp",
				"x-requested-with",
			},
			http.PHeaderOrderKey: {":authority", ":method", ":path", ":scheme"},
		}
		resp, err = twitter.client.Do(req)
		if err != nil {
			twitter.logger.Error("%d | Failed to send captcha result request: %s", twitter.index, err)
			return false
		}
		defer resp.Body.Close()
		bodyBytes, err = io.ReadAll(reader)
		if err != nil {
			extra.Logger{}.Error("%d | Failed to read unfreeze GET response body: %s", twitter.index, err.Error())
			continue
		}

		twitter.cookies.SetCookieFromResponse(resp)
		fmt.Println("ниже код fc/a/ без гейм токена")
		fmt.Println(resp.StatusCode)
		fmt.Println(string(bodyBytes))
	}

	return false
}

func (twitter *Twitter) BypassFuncaptcha(solution string) bool {
	tokenID := strings.Split(solution, "|")[0]

	fmt.Println("solution token is this: ", tokenID)

	URL := fmt.Sprintf("https://client-api.arkoselabs.com/fc/assets/ec-game-core/game-core/1.17.1/standard/index.html?session=%s&r=eu-west-1&meta=3&meta_width=558&meta_height=523&metabgclr=transparent&metaiconclr=%23555555&guitextcolor=%23000000&lang=en&pk=0152B4EB-D2DC-460A-89A1-629838B529C9&at=40&ag=101&cdn_url=https%3A%2F%2Fclient-api.arkoselabs.com%2Fcdn%2Ffc&lurl=https%3A%2F%2Faudio-eu-west-1.arkoselabs.com&surl=https%3A%2F%2Fclient-api.arkoselabs.com&smurl=https%3A%2F%2Fclient-api.arkoselabs.com%2Fcdn%2Ffc%2Fassets%2Fstyle-manager&theme=default", tokenID)

	req, err := http.NewRequest("GET", URL, nil)
	if err != nil {
		twitter.logger.Error("%d | Failed to build funcaptcha init request: %s", twitter.index, err)
		return false
	}
	req.Header = http.Header{
		"accept":             {"text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7"},
		"cookie":             {twitter.cookies.CookiesToHeader()},
		"sec-ch-ua-mobile":   {"?0"},
		"sec-ch-ua-platform": {`"Windows"`},
		"sec-fetch-dest":     {"iframe"},
		"sec-fetch-mode":     {"navigate"},
		"sec-fetch-site":     {"same-origin"},
		http.HeaderOrderKey: {
			"accept",
			"cookie",
			"sec-ch-ua-mobile",
			"sec-ch-ua-platform",
			"sec-fetch-dest",
			"sec-fetch-mode",
			"sec-fetch-site",
			"user-agent",
		},
		http.PHeaderOrderKey: {":authority", ":method", ":path", ":scheme"},
	}
	resp, err := twitter.client.Do(req)
	if err != nil {
		twitter.logger.Error("%d | Failed to send funcaptcha init request: %s", twitter.index, err)
		return false
	}
	defer resp.Body.Close()
	bodyText, err := io.ReadAll(resp.Body)
	if err != nil {
		twitter.logger.Error("%d | Failed to read response of funcaptcha init request: %s", twitter.index, err)
		return false
	}
	fmt.Println(resp.StatusCode)
	fmt.Println(string(bodyText))

	return true
}

func (twitter *Twitter) SolveCaptcha() string {
	captcha := api2captcha.NewClient(twitter.config.Info.TwoCaptchaKey)
	captcha.SoftId = 13289704
	captcha.DefaultTimeout = 120
	captcha.PollingInterval = 5

	captchaData := api2captcha.FunCaptcha{
		SiteKey:   "0152B4EB-D2DC-460A-89A1-629838B529C9",
		Url:       "https://twitter.com/account/access",
		Surl:      "https://client-api.arkoselabs.com",
		UserAgent: twitter.userAgent,
	}
	req := captchaData.ToRequest()

	req.SetProxy("HTTP", twitter.proxy)

	solution, err := captcha.Solve(req)
	if err != nil {
		twitter.logger.Error("%d | Failed to solve funcaptcha: %s", err)
		return ""
	}

	twitter.logger.Success("%d | Captcha solved!", twitter.index)
	return solution
}

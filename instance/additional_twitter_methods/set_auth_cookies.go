package additional_twitter_methods

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Danny-Dasilva/CycleTLS/cycletls"
	"math/rand"
	"strings"
	"twitter/utils"
)

// SetAuthCookies sets authentication cookies and updates headers for a client. Returns auth_token, csrf_token and err
func SetAuthCookies(accountIndex int, cookieClient *utils.CookieClient, twitterAuth string) (string, string, error) {
	csrfToken := ""
	authToken := ""

	// json cookies
	if strings.Contains(twitterAuth, "[") && strings.Contains(twitterAuth, "]") {
		jsonPart := strings.Split(strings.Split(twitterAuth, "[")[1], "]")[0]
		var cookiesJson []map[string]string
		if err := json.Unmarshal([]byte(jsonPart), &cookiesJson); err != nil {
			return "", "", fmt.Errorf("%d | Failed to decode account json cookies: %v", accountIndex, err)
		}

		for _, cookie := range cookiesJson {
			if name, ok := cookie["name"]; ok {
				value := cookie["value"]
				cookieClient.AddCookies([]cycletls.Cookie{{Name: name, Value: value}})
				if name == "ct0" {
					csrfToken = value
				}
				if name == "auth_token" {
					authToken = value
				}
			}
		}

		// auth token
	} else if len(twitterAuth) < 60 {
		csrfToken = generateMD5Token()
		cookieClient.AddCookies([]cycletls.Cookie{
			{Name: "auth_token", Value: twitterAuth},
			{Name: "ct0", Value: csrfToken},
			{Name: "des_opt_in", Value: "Y"},
		})
		authToken = twitterAuth
	}

	if csrfToken == "" {
		return "", "", errors.New("failed to get csrf token")
	}

	return authToken, csrfToken, nil
}

func generateMD5Token() string {
	randBytes := make([]byte, 32)
	rand.Read(randBytes)
	hasher := md5.New()
	hasher.Write(randBytes)
	return hex.EncodeToString(hasher.Sum(nil))
}

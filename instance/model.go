package instance

import (
	http "github.com/Danny-Dasilva/fhttp"
	"twitter/extra"
	"twitter/instance/additional_twitter_methods"
	"twitter/utils"
)

type Twitter struct {
	index     int
	authToken string
	proxy     string
	config    extra.Config
	queryID   extra.QueryIDs
	userAgent string

	ct0      string
	Username string

	client  *http.Client
	cookies *utils.CookieClient
	logger  extra.Logger
}

func (twitter *Twitter) InitTwitter(index int, authToken string, proxy string, config extra.Config, queryID extra.QueryIDs) (bool, string) {
	twitter.index = index
	twitter.authToken = authToken
	twitter.proxy = proxy
	twitter.config = config
	twitter.queryID = queryID
	return twitter.prepareClient()
}

func (twitter *Twitter) prepareClient() (bool, string) {
	var err error

	for i := 0; i < twitter.config.Info.MaxTasksRetries; i++ {
		tlsBuildInstance := utils.GetRandomTLSConfig()

		twitter.userAgent = tlsBuildInstance.UserAgent
		twitter.client = utils.CreateHttpClient(twitter.proxy, tlsBuildInstance)
		twitter.cookies = utils.NewCookieClient()
		twitter.authToken, twitter.ct0, err = additional_twitter_methods.SetAuthCookies(twitter.index, twitter.cookies, twitter.authToken)
		if err != nil {
			continue
		}

		twitter.Username, twitter.ct0, err = additional_twitter_methods.GetTwitterUsername(twitter.index, twitter.client, twitter.cookies, twitter.queryID.BearerToken, twitter.ct0)

		if err != nil {
			if twitter.Username == "locked" {
				return false, "locked"

			} else if twitter.Username == "failed_auth" {
				return false, "failed_auth"
			}
		} else {
			return true, "ok"
		}
	}

	twitter.logger.Error("Failed to prepare client.")
	return false, "unknown"
}

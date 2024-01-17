package grabber

import (
	http "github.com/Danny-Dasilva/fhttp"
	"twitter/extra"
	"twitter/utils"
)

type Grabber struct {
	index     int
	login     string
	password  string
	proxy     string
	config    extra.Config
	queryID   extra.QueryIDs
	userAgent string

	authToken  string
	ct0        string
	username   string
	guestToken string
	flowToken  string

	client  *http.Client
	cookies *utils.CookieClient
	logger  extra.Logger
}

func (grabber *Grabber) InitGrabber(index int, login string, password string, proxy string, config extra.Config, queryIDs extra.QueryIDs) bool {
	grabber.index = index + 1
	grabber.login = login
	grabber.password = password
	grabber.proxy = proxy
	grabber.config = config
	grabber.queryID = queryIDs
	return grabber.prepareClient()
}

func (grabber *Grabber) prepareClient() bool {
	//var err error

	for i := 0; i < 5; i++ {
		tlsBuildInstance := utils.GetRandomTLSConfig()

		grabber.userAgent = tlsBuildInstance.UserAgent
		grabber.client = utils.CreateHttpClient(grabber.proxy, tlsBuildInstance)
		grabber.cookies = utils.NewCookieClient()

		return true
	}

	grabber.logger.Error("Failed to prepare grabber client.")
	return false
}

func (grabber *Grabber) GrabTheCookie() string {
	ok := grabber.Flow()
	if ok {
		return grabber.authToken
	}
	return ""
}

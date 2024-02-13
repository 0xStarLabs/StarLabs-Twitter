package grabber

import (
	tlsClient "github.com/bogdanfinn/tls-client"
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

	client  tlsClient.HttpClient
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
	var err error

	for i := 0; i < 5; i++ {
		grabber.client, err = utils.CreateHttpClient(grabber.proxy)
		if err != nil {
			return false
		}
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

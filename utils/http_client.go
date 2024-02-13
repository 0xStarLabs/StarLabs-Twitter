package utils

import (
	"fmt"
	http "github.com/bogdanfinn/fhttp"
	tlsClient "github.com/bogdanfinn/tls-client"
	"github.com/bogdanfinn/tls-client/profiles"
	"strings"
	"twitter/extra"
)

func CreateHttpClient(proxies string) (tlsClient.HttpClient, error) {

	options := []tlsClient.HttpClientOption{
		tlsClient.WithClientProfile(profiles.Chrome_120),
		tlsClient.WithRandomTLSExtensionOrder(),
		tlsClient.WithInsecureSkipVerify(),
		tlsClient.WithTimeoutSeconds(30),
	}
	if proxies != "" {
		options = append(options, tlsClient.WithProxyUrl(fmt.Sprintf("http://%s", proxies)))
	}

	client, err := tlsClient.NewHttpClient(tlsClient.NewNoopLogger(), options...)
	if err != nil {

		extra.Logger{}.Error("Failed to create Http Client: %s", err)
		return nil, err
	}

	return client, nil
}

func CookiesToHeader(allCookies map[string][]*http.Cookie) string {
	var cookieStrs []string
	for _, cookies := range allCookies {
		for _, cookie := range cookies {
			cookieStrs = append(cookieStrs, cookie.Name+"="+cookie.Value)
		}
	}
	return strings.Join(cookieStrs, "; ")
}

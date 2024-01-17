package utils

import (
	"context"
	"fmt"
	"github.com/Danny-Dasilva/CycleTLS/cycletls"
	http "github.com/Danny-Dasilva/fhttp"
	"golang.org/x/net/proxy"
	"math/rand"
	"net"
	"net/url"
	"time"
)

func CreateHttpClient(proxies string, tlsBuildInstance TLSBuild) *http.Client {
	var client *http.Client

	if proxies != "" {
		proxyFormat := fmt.Sprintf("http://%s", proxies)

		proxyDialer, err := createProxyDialer(proxyFormat)
		if err != nil {
			fmt.Println("Error creating proxy dialer:", err)
		}

		client = &http.Client{
			Transport: cycletls.NewTransportWithProxy(tlsBuildInstance.JA3, tlsBuildInstance.UserAgent, proxyDialer),
			Timeout:   60 * time.Second,
		}
	} else {
		client = &http.Client{
			Transport: cycletls.NewTransport(tlsBuildInstance.JA3, tlsBuildInstance.UserAgent),
			Timeout:   60 * time.Second}
	}

	return client
}

func createProxyDialer(proxyStr string) (proxy.ContextDialer, error) {
	if proxyStr == "" {
		return nil, fmt.Errorf("proxy string is empty")
	}

	proxyURL, err := url.Parse(proxyStr)
	if err != nil {
		return nil, fmt.Errorf("error parsing proxy URL: %v", err)
	}

	dialer, err := proxy.FromURL(proxyURL, proxy.Direct)
	if err != nil {
		return nil, fmt.Errorf("error creating proxy dialer from URL: %v", err)
	}

	return &contextDialerWrapper{dialer: dialer}, nil
}

type contextDialerWrapper struct {
	dialer proxy.Dialer
}

func (d *contextDialerWrapper) DialContext(ctx context.Context, network, address string) (net.Conn, error) {
	return d.dialer.Dial(network, address)
}

func GetRandomTLSConfig() TLSBuild {
	rand.Seed(time.Now().UnixNano())

	builds := make([]string, 0, len(TLSBuilds))
	for k := range TLSBuilds {
		builds = append(builds, k)
	}
	randomBuild := builds[rand.Intn(len(builds))]

	randomBuildInstance := TLSBuilds[randomBuild]

	return TLSBuild{
		JA3:       randomBuildInstance["ja3"],
		UserAgent: randomBuildInstance["user-agent"],
		Version:   randomBuildInstance["version"],
	}
}

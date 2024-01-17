package utils

import (
	"github.com/Danny-Dasilva/CycleTLS/cycletls"
	http "github.com/Danny-Dasilva/fhttp"
	"strings"
)

type CookieClient struct {
	Cookies []cycletls.Cookie
}

func NewCookieClient() *CookieClient {
	return &CookieClient{Cookies: []cycletls.Cookie{}}
}

func (jar *CookieClient) GetCookieValue(name string) (string, bool) {
	for _, cookie := range jar.Cookies {
		if cookie.Name == name {
			return cookie.Value, true
		}
	}
	return "", false
}

func (jar *CookieClient) AddCookies(cookies []cycletls.Cookie) {
	for _, cookie := range cookies {
		jar.Cookies = append(jar.Cookies, cookie)
	}
}

func (jar *CookieClient) SetCookieFromResponse(resp *http.Response) {
	for _, httpCookie := range resp.Cookies() {
		updated := false
		newCookie := cycletls.Cookie{
			Name:  httpCookie.Name,
			Value: httpCookie.Value,
		}

		if httpCookie.Path != "" {
			newCookie.Path = httpCookie.Path
		}

		if httpCookie.Domain != "" {
			newCookie.Domain = httpCookie.Domain
		}

		if httpCookie.MaxAge > 0 {
			newCookie.MaxAge = httpCookie.MaxAge
		}

		if httpCookie.Secure {
			newCookie.Secure = httpCookie.Secure
		}

		if httpCookie.HttpOnly {
			newCookie.HTTPOnly = httpCookie.HttpOnly
		}

		if !httpCookie.Expires.IsZero() {
			newCookie.Expires = httpCookie.Expires
		}

		for i, existingCookie := range jar.Cookies {
			if existingCookie.Name == httpCookie.Name {
				jar.Cookies[i] = newCookie
				updated = true
				break
			}
		}

		if !updated {
			jar.Cookies = append(jar.Cookies, newCookie)
		}
	}
}

// CookiesToHeader this function converts cookies from the []cycletls.Cookie{} format to a header format
// that simply contains a string of all cookies in the name=value format;
func (jar *CookieClient) CookiesToHeader() string {
	var cookieStrs []string
	for _, cookie := range jar.Cookies {
		cookieStrs = append(cookieStrs, cookie.Name+"="+cookie.Value)
	}
	return strings.Join(cookieStrs, "; ")
}

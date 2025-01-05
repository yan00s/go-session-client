package gosessionclient

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"golang.org/x/net/http2"
	"h12.io/socks"
)

// Scheme: protocol://username:passw@ip:port or protocol://ip:port
// Example: http://kyabx3d:hjehxapxxa@178.0.9.209:4791
// And also need enter timeout in seconds
func (session *Session) SetProxy(proxy string, timeout int) error {

	proxy = strings.Replace(proxy, "\r", "", -1)

	parsedURL, err := url.Parse(proxy)
	if err != nil {
		return customErr("failed to parse proxy URL", err)
	}

	tr, err := getTransport(parsedURL)

	if err != nil {
		return err
	}

	if err := http2.ConfigureTransport(tr); err != nil {
		return err
	}

	session.Client.Transport = tr
	session.Client.Timeout = time.Duration(timeout) * time.Second
	return nil
}

func getTransport(parsedURL *url.URL) (*http.Transport, error) {
	var tr *http.Transport

	switch parsedURL.Scheme {
	case "http":
		tr = &http.Transport{Proxy: http.ProxyURL(parsedURL), TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}
	case "socks4", "socks4a", "socks5": // "socks5://user:password@127.0.0.1:1080?timeout=5s"
		dialSocksProxy := socks.Dial(parsedURL.String())
		tr = &http.Transport{Dial: dialSocksProxy}
	default:
		return nil, fmt.Errorf("unsupported proxy scheme: %s", parsedURL.Scheme)
	}
	return tr, nil
}

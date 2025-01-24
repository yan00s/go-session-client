package gosessionclient

import (
	"crypto/tls"
	"strings"
)

func checkLocalHost(proxy string) *tls.Config {
	if strings.Contains(proxy, "127.0.0.1") || strings.Contains(proxy, "localhost") {
		return &tls.Config{InsecureSkipVerify: true}
	} else {
		return nil
	}
}

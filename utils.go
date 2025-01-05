package gosessionclient

import (
	"crypto/tls"
	"fmt"
	"strings"
)

func checkLocalHost(proxy string) *tls.Config {
	if strings.Contains(proxy, "127.0.0.1") || strings.Contains(proxy, "localhost") {
		return &tls.Config{InsecureSkipVerify: true}
	} else {
		return nil
	}
}

func customErr(text string, err error) error {
	return fmt.Errorf("%v: %w", text, err)
}

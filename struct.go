package gosessionclient

import "net/http"

type Session struct {
	Client  *http.Client
	Headers http.Header
}

type Response struct {
	Body    []byte
	Headers http.Header
	Cookies []*http.Cookie
	Status  int
	Err     error
}

package gosessionclient

import "net/http"

type Session struct {
	Client  *http.Client
	Headers http.Header
}

type Response struct {
	Body   []byte
	Status int
	Err    error
}

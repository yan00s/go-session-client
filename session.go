package gosessionclient

import (
	"bytes"
	"context"
	"crypto/tls"
	"io"
	"net/http"
	"net/http/cookiejar"
	"time"

	"golang.org/x/net/http2"
	"golang.org/x/net/publicsuffix"
)

func CreateSession(pcAgent bool) Session {
	session := Session{}
	jar, _ := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})

	tr := &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}

	client := &http.Client{Transport: tr, Jar: jar}
	http2.ConfigureTransport(client.Transport.(*http.Transport))

	session.Client = client
	session.Headers = GenerateHeaders(pcAgent)
	return session
}

func (session *Session) sendReq(ctx context.Context, urlstr string, method string, reader io.Reader) *Response {
	result := &Response{}

	req, err := http.NewRequestWithContext(ctx, method, urlstr, reader)
	if err != nil {
		result.Err = customErr("Error on create request", err)
		return result
	}

	req.Header = session.Headers

	resp, err := session.Client.Do(req)

	if err != nil {
		result.Err = customErr("Error on send requests", err)
		return result
	}

	body, err := readBody(resp)

	if err != nil {
		result.Err = customErr("Error on read result response", err)
		return result
	}

	result.Body = body
	result.Status = resp.StatusCode
	result.Headers = resp.Header
	result.Cookies = resp.Cookies()

	return result
}

func (session *Session) sendReqWithRetry(ctx context.Context, urlstr string, method string, reader io.Reader, retries int, delay time.Duration) *Response {
	var result *Response
	for i := 0; i < retries; i++ {
		result = session.sendReq(ctx, urlstr, method, reader)
		if result.Err == nil {
			break
		}
		time.Sleep(delay)
	}
	return result
}

// SendReq sends an HTTP request with the specified method, retry count, timeout, and retry delay, along with optional data (if dataStr is provided, it will be used as the request body).
// Parameters:
// - url: The URL to which the request is sent.
// - method: The HTTP method (GET, POST, PUT, DELETE, PATCH, HEAD, OPTIONS, TRACE).
// - retryCount: The number of retry attempts in case of an error.
// - timeout: The timeout duration for each request (e.g., 30 * time.Second).
// - retryDelay: The delay between retry attempts (e.g., 2 * time.Second).
// - dataStr: Optional data to be sent in the request body.
// Returns: A pointer to a Response struct containing the result of the request.
func (session *Session) SendReq(url, method string, retryCount int, timeout time.Duration, retryDelay time.Duration, dataStr ...string) *Response {
	var reader io.Reader

	if len(dataStr) > 0 {
		data := []byte(dataStr[0])
		reader = bytes.NewReader(data)
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	return session.sendReqWithRetry(ctx, url, method, reader, retryCount, retryDelay) // Example retries
}

func (res *Response) String() string {
	return string(res.Body)
}

func readBody(resp *http.Response) ([]byte, error) {
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

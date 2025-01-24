package gosessionclient

import (
	"bytes"
	"context"
	"fmt"
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

	tr := &http.Transport{}

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
		result.Err = fmt.Errorf("Error on create request: %w", err)
		return result
	}

	req.Header = session.Headers

	resp, err := session.Client.Do(req)

	if err != nil {
		result.Err = fmt.Errorf("Error on send requests: %w", err)
		return result
	}

	body, err := readBody(resp)
	defer resp.Body.Close()

	if err != nil {
		result.Err = fmt.Errorf("Error on read result response: %w", err)
		return result
	}

	result.Body = body
	result.Status = resp.StatusCode
	result.Headers = resp.Header
	result.Cookies = resp.Cookies()

	return result
}

// SendReqWithRetry sends an HTTP request with retry logic.
// This function allows specifying the HTTP method, timeout duration, number of retries,
// delay between retries, and optional request body data.
//
// Parameters:
// - url: The target URL for the request.
// - method: The HTTP method to use (e.g., GET, POST, PUT, DELETE, PATCH, HEAD, OPTIONS, TRACE).
// - timeout: The maximum duration for the one request context (e.g., 30 * time.Second).
// - retryCount: The number of retry attempts if the request fails.
// - retryDelay: The duration to wait between retries (e.g., 2 * time.Second).
// - dataStr: (Optional) A string to be sent as the request body, if provided.
//
// Returns:
// A pointer to a Response struct containing the result of the request, including any errors.
func (session *Session) SendReqWithRetry(url, method string, timeout time.Duration, retryCount int, retryDelay time.Duration, dataStr ...string) *Response {
	var result *Response
	var reader io.Reader

	if len(dataStr) > 0 {
		// Use only first element as data for request
		data := []byte(dataStr[0])
		reader = bytes.NewReader(data)
	}

	for i := 0; i < retryCount; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()
		result = session.sendReq(ctx, url, method, reader)

		if result.Err == nil {
			break
		}

		if ctx.Err() != nil {
			result.Err = fmt.Errorf("Request timed out: %w", ctx.Err())
			break
		}

		time.Sleep(retryDelay)
	}
	return result
}

// SendReq sends an HTTP request without retry logic.
// This function allows specifying the HTTP method, timeout duration, and optional request body data.
//
// Parameters:
// - url: The target URL for the request.
// - method: The HTTP method to use (e.g., GET, POST, PUT, DELETE, PATCH, HEAD, OPTIONS, TRACE).
// - timeout: The maximum duration for the request context (e.g., 30 * time.Second).
// - dataStr: (Optional) A string to be sent as the request body, if provided.
//
// Returns:
// A pointer to a Response struct containing the result of the request, including any errors.
func (session *Session) SendReq(url, method string, timeout time.Duration, dataStr ...string) *Response {
	var reader io.Reader

	if len(dataStr) > 0 {
		// Use only first element as data for request
		data := []byte(dataStr[0])
		reader = bytes.NewReader(data)
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	return session.sendReq(ctx, url, method, reader)
}

func (res *Response) String() string {
	return string(res.Body)
}

func readBody(resp *http.Response) ([]byte, error) {
	// limit 100MB
	limitedReader := io.LimitReader(resp.Body, 100*1024*1024)
	body, err := io.ReadAll(limitedReader)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}
	return body, nil
}

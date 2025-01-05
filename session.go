package gosessionclient

import (
	"bytes"
	"crypto/tls"
	"io"
	"net/http"
	"net/http/cookiejar"

	"golang.org/x/net/http2"
	"golang.org/x/net/publicsuffix"
)

func CreateSession() Session {
	session := Session{}
	jar, _ := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})

	tr := &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}

	client := &http.Client{Transport: tr, Jar: jar}
	http2.ConfigureTransport(client.Transport.(*http.Transport))

	session.Client = client
	session.Headers = GenerateHeaders()
	return session
}

func (session *Session) sendReq(urlstr string, method string, reader io.Reader) *Response {
	result := &Response{}

	req, err := http.NewRequest(method, urlstr, reader)
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

	return result
}

// SendReq sends an HTTP request with the specified method and optional data (if dataStr is provided, it will be used as request body).
// Supported HTTP methods:
// 1. GET     - Retrieve data from the server
// 2. POST    - Send data to the server
// 3. PUT     - Update data on the server (full replacement)
// 4. DELETE  - Remove data from the server
// 5. PATCH   - Update data on the server (partial update)
// 6. HEAD    - Request headers only (no body)
// 7. OPTIONS - Retrieve supported HTTP methods for a resource
// 8. TRACE   - Perform a trace of the route to the server
func (session *Session) SendReq(url, method string, dataStr ...string) *Response {
	var reader io.Reader

	if len(dataStr) > 0 {
		data := []byte(dataStr[0])
		reader = bytes.NewReader(data)
	}

	return session.sendReq(url, method, reader)
}

func (res *Response) String() string {
	return string(res.Body)
}

func readBody(resp *http.Response) ([]byte, error) {
	body, err := io.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}
	return body, nil
}

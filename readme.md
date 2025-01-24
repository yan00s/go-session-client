# Go-Session-Client

Go-Session-Client is a simple and lightweight HTTP client for Go, with support for sessions, cookies and proxies.

## Features
- Persistent sessions with cookies
- Support for SOCKS4(A), SOCKS5, and HTTP proxies
- Simplified request methods (GET, POST, etc.)
- A lot of user agent (mobile and PC)

## Installation

go get github.com/yan00s/go-session-client

Example Usage:

```Go
package main

import (
	"fmt"
	"time"

	session "github.com/yan00s/go-session-client"
)

func main() {
	ses := session.CreateSession(true)

	/// Trying to set Proxy
	proxyStr := "http://username:passw@100.100.100.100:2000"
	timeout := 3 // seconds
	if err := ses.SetProxy(proxyStr, timeout); err != nil {
		fmt.Println(fmt.Errorf("Error in setting up proxy: %w", err))
	}
	///

	fmt.Println()

	/// Trying to make a get request
	fmt.Println("Trying to make request with 1 try on icanhazip.com with a 10 second timeout per request")
	resp := ses.SendReq("https://icanhazip.com", "GET", 5*time.Second)

	if resp.Err != nil {
		fmt.Println("Error in making request:", resp.Err.Error())
		return
	}

	fmt.Println("Response:", resp.String())
	fmt.Println("Status:", resp.Status)
	///
}
```

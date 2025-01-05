# Go-Session-Client

Go-Session-Client is a simple and lightweight HTTP client for Go, with support for sessions, cookies and proxies.

## Features
- Persistent sessions with cookies
- Support for SOCKS4(A), SOCKS5, and HTTP proxies
- Simplified request methods (GET, POST, etc.)
- A lot of user agent (mobile and PC)

## Installation

<<<<<<< HEAD

=======
>>>>>>> 09802ea (Moved files from src to root and delete example.go)
go get github.com/yan00s/go-session-client

Example Usage:

```Go
package main

import (
	"fmt"
	"os"

	session "github.com/yan00s/go-session-client"
)

func main() {
	ses := session.CreateSession()
	timeout := 15 // in seconds

	if err := ses.SetProxy("http://username:passw@100.100.100.100:2000", timeout); err != nil {
		fmt.Println(fmt.Errorf("err in set up proxy: %w", err))
		os.Exit(1)
	}

	resp := ses.SendReq("https://icanhazip.com", "GET")

	fmt.Println("response:", resp.String())
	fmt.Println("status:", resp.Status)
	fmt.Println("errors:", resp.Err)
}
```

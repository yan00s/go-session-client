package main

import (
	"fmt"
	"os"
	"time"

	session "github.com/yan00s/go-session-client"
)

func main() {
	ses := session.CreateSession(true)
	timeout := 15 // in seconds

	if err := ses.SetProxy("http://username:passw@100.100.100.100:2000", timeout); err != nil {
		fmt.Println(fmt.Errorf("err in set up proxy: %w", err))
		os.Exit(1)
	}

	resp := ses.SendReq("https://icanhazip.com", "GET", 1, 10*time.Second, 1*time.Second)

	fmt.Println("response:", resp.String())
	fmt.Println("status:", resp.Status)
	fmt.Println("errors:", resp.Err)
}

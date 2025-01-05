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

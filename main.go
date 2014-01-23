package main

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
)

func main() {
	fmt.Fprintf(os.Stderr, "IRCCloud Streamer")
	h, e := http.NewRequest("GET", "https://www.irccloud.com/chat/stream", nil)
	if e != nil {
		panic(e)
	}
	client := &http.Client{}

	resp, e := client.Do(h)
	if e != nil {
		panic(e)
	}
	h.Header.Add("Cookie", "session=asdfasdfasd")

	reader := bufio.NewReader(resp.Body)
	for {
		s, err := reader.ReadString('\n')
		if err != nil {
			panic(err)
		}
		fmt.Printf("%s", s)
	}
}

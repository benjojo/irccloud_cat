package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

func main() {
	fmt.Fprintf(os.Stderr, "IRCCloud Streamer")

	sb, e := ioutil.ReadFile("./.session")
	if e != nil {
		log.Fatal("cannot read session cookie. Please make a .session file and put your cookie in there.")
	}
	lines := strings.Split(string(sb), "\n")
	if len(lines) != 1 {
		log.Fatal("There is more than one line in the session file. wat.")
	}
	h, e := http.NewRequest("GET", "https://www.irccloud.com/chat/stream", nil)
	if e != nil {
		panic(e)
	}
	client := &http.Client{}
	h.Header.Add("Cookie", fmt.Sprintf("session=%s", lines[0]))

	resp, e := client.Do(h)
	if e != nil {
		panic(e)
	}

	reader := bufio.NewReader(resp.Body)
	for {
		s, err := reader.ReadString('\n')
		if err != nil {
			panic(err)
		}
		fmt.Printf("%s", s)
	}
}

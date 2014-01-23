package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

type MessageDump struct {
	Bid  int    `json:"bid"`
	Type string `json:"type"`
	URL  string `json:"url"`
}

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
		wat := MessageDump{}
		json.Unmarshal([]byte(s), &wat)
		if wat.URL != "" && wat.Type == "oob_include" {
			h, e = http.NewRequest("GET", "https://www.irccloud.com"+wat.URL, nil)
			h.Header.Add("Cookie", fmt.Sprintf("session=%s", lines[0]))
			tresp, _ := client.Do(h)
			ioutil.ReadAll(tresp.Body)
			fmt.Fprintf(os.Stderr, "Did the OOB include, the stream should work ~forever~ now.")
		}
		if err != nil {
			panic(err)
		}
		fmt.Printf("%s", s)
	}
}

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

type Msg struct {
	Bid         float64 `json:"bid"`
	Chan        string  `json:"chan"`
	Cid         float64 `json:"cid"`
	Eid         float64 `json:"eid"`
	From        string  `json:"from"`
	FromHost    string  `json:"from_host"`
	FromMode    string  `json:"from_mode"`
	FromName    string  `json:"from_name"`
	Hostmask    string  `json:"hostmask"`
	IdentPrefix string  `json:"ident_prefix"`
	Msg         string  `json:"msg"`
	Type        string  `json:"type"`
}

var FileLogging bool

func main() {
	FileLogging = true
	fmt.Fprintf(os.Stderr, "IRCCloud Streamer")
	streamtoken := ""
	if len(os.Args) == 2 {
		streamtoken = os.Args[1]
	} else {
		sb, e := ioutil.ReadFile("./.session")
		if e != nil {
			log.Fatal("cannot read session cookie. Please make a .session file and put your cookie in there.")
		}
		lines := strings.Split(string(sb), "\n")
		if len(lines) != 1 {
			log.Fatal("There is more than one line in the session file. wat.")
		}
		streamtoken = lines[0]
	}
	h, e := http.NewRequest("GET", "https://www.irccloud.com/chat/stream", nil)
	if e != nil {
		panic(e)
	}
	client := &http.Client{}
	h.Header.Add("Cookie", fmt.Sprintf("session=%s", streamtoken))

	resp, e := client.Do(h)
	if e != nil {
		panic(e)
	}
	if FileLogging {
		f, err := os.OpenFile("./log", os.O_APPEND|os.O_WRONLY, 0600)
		if err != nil {
			panic("cannot open the log file.")
		}
	}
	reader := bufio.NewReader(resp.Body)
	for {
		s, err := reader.ReadString('\n')
		wat := MessageDump{}
		json.Unmarshal([]byte(s), &wat)
		if wat.URL != "" && wat.Type == "oob_include" {
			h, e = http.NewRequest("GET", "https://www.irccloud.com"+wat.URL, nil)
			h.Header.Add("Cookie", fmt.Sprintf("session=%s", streamtoken))
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

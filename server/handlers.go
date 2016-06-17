// Copyright 2016 Olivier Wulveryck
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package server

import (
	"log"
	"net/http"
	"regexp"
	"time"

	"github.com/gorilla/websocket"
)

type result struct {
	Topic string    `json:"topic"`
	Date  time.Time `json:"date"`
	Total int64     `json:"total"`
	Score float64   `json:"score"`
}

type message struct {
	Topic  string    `json:"topic"`
	Sender string    `json:"sender"`
	msg    string    `json:"message"`
	Like   bool      `json:"like"`
	Date   time.Time `json:"-"`
}

type phoneMessage struct {
	Name    string `json:"name"`
	Device  string `json:"device"`
	Message string `json:"message"`
}

type communication struct {
	msg  phoneMessage
	Chan chan phoneMessage
}

var topics map[string][]int64

func init() {
	topics = make(map[string][]int64, 0)
}

var upgrader = websocket.Upgrader{} // use default options

var communicationChannel = make(chan communication)

func slideJoin(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return

	}
	defer c.Close()
	type tempo struct {
		Channel chan phoneMessage
		State   string
	}
	//	var attendee = make(map[string]tempo)
	nodeid := 0
	for {
		message := <-communicationChannel
		log.Println("message received ", message)
		type node struct {
			Name string `json:"name"`
			ID   int    `json:"id"`
			Icon string `json:"icon"`
		}
		type link struct {
			Source int `json:"source"`
			Target int `json:"target"`
		}
		type reply struct {
			Nodes []node `json:"nodes"`
			Links []link `json:"links"`
		}
		var ios = regexp.MustCompile(`(?i).*ios|iphone.*`)
		var android = regexp.MustCompile(`(?i).*android.*`)
		var icon string
		switch {
		case ios.MatchString(message.msg.Device):
			icon = "/img/iphone-phone-color.png"
		case android.MatchString(message.msg.Device):
			icon = "/img/android-phone-color.png"
		default:
			icon = "/img/smartphone.png"
		}
		nodeid = nodeid + 1
		myreply := reply{
			[]node{
				node{message.msg.Name, nodeid, icon},
			},
			[]link{
				link{0, 0},
			},
		}
		log.Println("Sending", myreply)
		err = websocket.WriteJSON(c, &myreply)
		if err != nil {
			log.Println(err)
		}
		log.Println("sent", myreply)
		/*
			attendee[message.msg.Name] = tempo{Channel: message.Chan, State: message.msg.State}
			for att, temp := range attendee {
				go func(att string, temp tempo) {
					channel := temp.Channel
					state := temp.State
					log.Println(state)
					if state == "start" {
						log.Printf("Sending to %v on channel %v", att, channel)
						channel <- msg{"A", "running"}
					}
					if state == "stop" {
						log.Printf("Sending to %v on channel %v", att, channel)
						channel <- msg{"A", "stopped"}
					}
				}(att, temp)
			}
		*/
	}
}

func phone(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return

	}
	defer c.Close()

	// Read incoming message and pass it to the hub
	var channel = make(chan phoneMessage)
	// launch a goroutine and wait
	go func(channel chan phoneMessage) {
		for {
			response := <-channel
			log.Println("about to send ", response)

			err = websocket.WriteJSON(c, response)
			if err != nil {
				log.Println("Unable to send message", err)
			}
		}
	}(channel)
	for {

		var message phoneMessage
		err := websocket.ReadJSON(c, &message)
		if err != nil {
			log.Println("Unable to read message", err)
		} else {
			log.Printf("=> %v is talking (%v)", message.Name, message)
			communicationChannel <- communication{msg: message, Chan: channel}
			log.Printf("=> Advertized ")

		}
	}
}

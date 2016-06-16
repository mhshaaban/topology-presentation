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

type msg struct {
	Name  string `json:"name"`
	State string `json:"state"`
}

type communication struct {
	msg  msg
	Chan chan msg
}

var topics map[string][]int64

func init() {
	topics = make(map[string][]int64, 0)
}

var upgrader = websocket.Upgrader{} // use default options

var communicationChannel = make(chan communication)

func orchestrator(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return

	}
	defer c.Close()
	type tempo struct {
		Channel chan msg
		State   string
	}
	var attendee = make(map[string]tempo)
	for {
		message := <-communicationChannel
		log.Println("message received ", message)
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
	var channel = make(chan msg)
	// launch a goroutine and wait
	go func(channel chan msg) {
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

		var message msg
		err := websocket.ReadJSON(c, &message)
		if err != nil {
			log.Println("Unable to read message", err)
		} else {
			log.Printf("=> %v is talking", message.Name)
			communicationChannel <- communication{msg: message, Chan: channel}
			log.Printf("=> Advertized ")

		}
	}
}

func progress(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return

	}
	defer c.Close()

	for {

		var message message
		err := websocket.ReadJSON(c, &message)
		if err != nil {
			log.Println("Unable to read message", err)
		} else {
			log.Printf("=> message: \n==> Topic:%v\n==>Sender:%v\n==>Date:%v ", message.Topic, message.Sender, message.Date)
		}

		var increment int64
		if message.Like {
			increment = 1
		} else {
			increment = -1
		}
		if _, ok := topics[message.Topic]; !ok {
			topics[message.Topic] = []int64{
				1,
				1,
			}
		}
		topics[message.Topic] = []int64{
			topics[message.Topic][0] + 1,
			topics[message.Topic][1] + increment,
		}
		var response result
		response.Topic = message.Topic
		response.Date = time.Now()
		response.Total = topics[message.Topic][0]
		response.Score = float64(topics[message.Topic][1] * 100 / topics[message.Topic][0])
		log.Println("about to send ", response)

		err = websocket.WriteJSON(c, response)
		if err != nil {
			log.Println("Unable to send message", err)
		}
	}
}

func getJSON(w http.ResponseWriter, r *http.Request) {
}

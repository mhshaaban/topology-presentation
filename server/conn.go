// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package server

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// Conn is an middleman between the websocket connection and the hub.
type Conn struct {
	// The websocket connection.
	ws *websocket.Conn

	// Buffered channel of outbound messages.
	send chan Message
}

// readPump pumps messages from the websocket connection to the hub.
func (c *Conn) readPump(ID string) {
	defer func() {
		hubs.RLock()
		hubs.h[ID].unregister <- c
		hubs.RUnlock()
		c.ws.Close()
	}()
	c.ws.SetReadLimit(maxMessageSize)
	c.ws.SetReadDeadline(time.Now().Add(pongWait))
	c.ws.SetPongHandler(func(string) error { c.ws.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		var message Message
		//err := websocket.ReadJSON(c.ws, &message)
		_, b, err := c.ws.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				log.Printf("==> [%v] error: %v", ID, err)
			}
			break
		}
		err = json.Unmarshal(b, &message)
		if err != nil {
			log.Printf("==> [%v] Received message but cannot unmarshal it, %v", ID, err)
		}
		log.Printf("==> [%v] Received message: %v (%s)", ID, message, string(b))

		hubs.RLock()
		hubs.h[ID].broadcast <- message
		hubs.RUnlock()
	}
}

// write writes a message with the given message type and payload.
func (c *Conn) write(mt int, payload Message) error {
	c.ws.SetWriteDeadline(time.Now().Add(writeWait))
	return websocket.WriteJSON(c.ws, &payload)
}

// writePump pumps messages from the hub to the websocket connection.
func (c *Conn) writePump(ID string) {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.ws.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				// The hub closed the channel.
				c.write(websocket.CloseMessage, Message{})
				return
			}

			c.ws.SetWriteDeadline(time.Now().Add(writeWait))
			w, err := c.ws.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			b, _ := json.Marshal(message)
			w.Write(b)

			//			websocket.WriteJSON(c.ws, &message)

			// Add queued chat messages to the current websocket message.
			n := len(c.send)
			for i := 0; i < n; i++ {
				//websocket.WriteJSON(c.ws, <-c.send)
				b, _ := json.Marshal(<-c.send)
				w.Write([]byte{'\n'})

				w.Write(b)
				//w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			if err := c.write(websocket.PingMessage, Message{}); err != nil {
				return
			}
		}
	}
}

// serveWs handles websocket requests from the peer.
func serveWs(w http.ResponseWriter, r *http.Request) {
	//Let's get the ID
	vars := mux.Vars(r)
	ID := vars["id"]
	log.Printf("=> Connection to %v", ID)
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	conn := &Conn{send: make(chan Message, 256), ws: ws}
	if ID == "" {
		log.Printf("No ID provided, bailing out")
		return
	}
	hubs.Lock()
	if _, ok := hubs.h[ID]; ok {
		log.Printf("==> [%v] A hub alredy exist for this ID", ID)
		hubs.h[ID].register <- conn
	} else {
		log.Printf("==> [%v] Creating a new hub", ID)
		hubs.h[ID] = &Hub{
			broadcast:   make(chan Message),
			register:    make(chan *Conn),
			unregister:  make(chan *Conn),
			connections: make(map[*Conn]bool),
		}
		go hubs.h[ID].run()
		hubs.h[ID].register <- conn
	}

	hubs.Unlock()
	go conn.writePump(ID)
	conn.readPump(ID)
}

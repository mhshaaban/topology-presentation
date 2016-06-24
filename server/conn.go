// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package server

import (
	"encoding/json"
	log "github.com/Sirupsen/logrus"
	"regexp"
	"time"

	"github.com/gorilla/websocket"
)

func init() {
	log.SetLevel(log.DebugLevel)
}

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Minute

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
func (c *Conn) readPump(Tag Tag, h *hub) {
	var contextLogger = log.WithFields(log.Fields{
		"Tag": Tag,
	})
	defer func() {
		h.unregister <- c
		c.ws.Close()
	}()
	c.ws.SetReadLimit(maxMessageSize)
	c.ws.SetReadDeadline(time.Now().Add(pongWait))
	c.ws.SetPongHandler(func(string) error { c.ws.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		var message Node
		//err := websocket.ReadJSON(c.ws, &message)
		t, b, err := c.ws.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				contextLogger.Error(err)
			}
			break
		}
		if t == websocket.PingMessage {
			continue
		}
		err = json.Unmarshal(b, &message)
		if err != nil {
			contextLogger.Error(err)
			return
		}
		contextLogger.Debug(message)

		found := false
		nodeTag := 0
		for i, node := range h.message.Nodes {
			if node.UUID == message.UUID {
				h.message.Nodes[i].Status = message.Status
				switch message.Status {
				case "initial":
					message.Color = "black"
				case "configured":
					message.Color = "cyan"
				case "started":
					message.Color = "green"
				case "stopped":
					message.Color = "orange"
				case "deleted":
					message.Color = "red"
				default:
					message.Color = "black"
				}
				h.message.Nodes[i].Color = message.Color
				found = true
			}
			nodeTag = i + 1
		}
		// New node
		if !found {

			var ios = regexp.MustCompile(`(?i).*ios|iphone.*`)
			var android = regexp.MustCompile(`(?i).*android.*`)
			if message.Icon == "" {
				var icon string
				switch {
				case ios.MatchString(message.Device):
					icon = "/img/iphone-phone-color.png"
				case android.MatchString(message.Device):
					icon = "/img/android-phone-color.png"
				default:
					icon = "/img/smartphone.png"
				}
				message.Icon = icon
			}
			message.Tag = nodeTag
			h.message.Nodes = append(h.message.Nodes, message)
			// Add a link to the previous node
			if nodeTag >= 1 {
				h.message.Links = append(h.message.Links, Link{Source: nodeTag, Target: nodeTag - 1})

			}
		}
		if message.Status == "error" {
			h.message.Message = "error"
		} else {
			h.message.Message = "info"
		}
		if len(h.message.Links) == 0 {
			// Add a dummy link for d3.js
			h.message.Links = append(h.message.Links, Link{Source: 0, Target: 0})
		}
		//log.Debug("==> [%v] Broadcasting message: %v", Tag, h.message)
		//h.broadcast <- *h.message
		h.process <- *h.message
		h.broadcast <- *h.message
	}
}

// write writes a message with the given message type and payload.
func (c *Conn) write(mt int, payload Message) error {
	c.ws.SetWriteDeadline(time.Now().Add(writeWait))
	return websocket.WriteJSON(c.ws, &payload)
}

// writePump pumps messages from the hub to the websocket connection.
func (c *Conn) writePump(Tag Tag) {
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
			if err := c.write(websocket.PingMessage, Message{Message: "ping"}); err != nil {
				return
			}
		}
	}
}

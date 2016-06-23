// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package server

// Hub maintains the set of active connections and broadcasts messages to the
// connections.
type Hub struct {
	// ID of the Hubs
	ID int

	// Registered connections.
	connections map[*Conn]bool

	// Inbound messages from the connections.
	broadcast chan Message

	// Register requests from the connections.
	register chan *Conn

	// Unregister requests from connections.
	unregister chan *Conn
}

// Hubs maintains the set of active Hubs
type Hubs struct {
	// Unregister
	unregister chan int

	// Registered hubs.
	hubs map[int]*Hub

	// Register requests from the connections.
	Request chan *Reply
}

type Reply struct {
	ID  int
	Rep chan *Hub
}

var AllHubs = Hubs{
	Request: make(chan *Reply),
	hubs:    make(map[int]*Hub),
}

// The main routine for registering the hubs
func (h *Hubs) Run() {
	for {
		select {
		case r := <-h.Request:
			if hub, ok := h.hubs[r.ID]; ok {
				r.Rep <- hub
			} else {
				//TODO create a new hub
				// And add it to the hubs map
			}
		case hub := <-h.unregister:
			if _, ok := h.hubs[hub]; ok {
				delete(h.hubs, hub)
			}
		}
	}
}
func (h *Hub) run() {
	for {
		select {
		case conn := <-h.register:
			h.connections[conn] = true
		case conn := <-h.unregister:
			if _, ok := h.connections[conn]; ok {
				delete(h.connections, conn)
				close(conn.send)
			}
			// If the last element has been removed exit)
			if len(h.connections) == 0 {
				AllHubs.unregister <- h.ID
				return
			}
		case message := <-h.broadcast:
			for conn := range h.connections {
				select {
				case conn.send <- message:
				default:
					close(conn.send)
					//delete(hub.connections, conn)
					delete(h.connections, conn)
				}
			}
		}
	}
}

// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package server

// Hub maintains the set of active connections and broadcasts messages to the
// connections.
type Hub struct {
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
	// Registered connections.
	connections map[*Hub]bool

	// Register requests from the connections.
	register chan *Hub

	// Unregister requests from connections.
	unregister chan *Hub
}

var AllHubs = Hubs{
	register:    make(chan *Hub),
	unregister:  make(chan *Hub),
	connections: make(map[*Hub]bool),
}

func (h *Hubs) run() {
	for {
		select {
		case hub := <-h.register:
			h.connections[hub] = true
		case hub := <-h.unregister:
			if _, ok := h.connections[hub]; ok {
				delete(h.connections, hub)
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
				AllHubs.unregister <- h
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

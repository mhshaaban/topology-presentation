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

/*
var hub = Hub{
	broadcast:   make(chan Message),
	register:    make(chan *Conn),
	unregister:  make(chan *Conn),
	connections: make(map[*Conn]bool),
}
*/

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
			// TODO
			// If the last element has been removed, clean the Hubs map
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

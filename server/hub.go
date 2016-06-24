// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package server

import (
	log "github.com/Sirupsen/logrus"
)

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
			if _, ok := h.hubs[r.ID]; !ok {
				//TODO create a new hub
				var hub = &Hub{
					ID:          r.ID,
					broadcast:   make(chan Message),
					register:    make(chan *Conn),
					unregister:  make(chan *Conn),
					connections: make(map[*Conn]bool),
				}
				var contextLogger = log.WithFields(log.Fields{
					"ID":  r.ID,
					"Hub": &hub,
				})
				contextLogger.Debug("New HUB")
				go hub.run()
				h.hubs[r.ID] = hub
				// By the end reply to the sender
			}
			r.Rep <- h.hubs[r.ID]
		case hub := <-h.unregister:
			log.Debug("In the hubs' unregister")
			if _, ok := h.hubs[hub]; ok {
				var contextLogger = log.WithFields(log.Fields{
					"ID":  hub,
					"Hub": h.hubs[hub],
				})
				contextLogger.Debug("Unregistering HUB")
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
			log.WithFields(log.Fields{
				"Connections": len(h.connections),
				"Connection":  conn,
				"Hub":         &h,
			}).Debug("Registerng connection")
		case conn := <-h.unregister:
			if _, ok := h.connections[conn]; ok {
				log.WithFields(log.Fields{
					"Connections": len(h.connections),
					"Connection":  conn,
					"Hub":         &h,
				}).Debug("Unregisterng connection")
				delete(h.connections, conn)
				close(conn.send)
			}
			// If the last element has been removed exit)
			if len(h.connections) == 0 {
				AllHubs.unregister <- h.ID
				return
			}
		case message := <-h.broadcast:
			log.WithFields(log.Fields{
				"Hub": &h,
			}).Debug("Broadcast")
			for conn := range h.connections {
				log.WithFields(log.Fields{
					"Connection": conn,
					"Hub":        &h,
				}).Debug("Sending...")
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

package server

import (
	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

// serveWs handles websocket requests from the peer.
func serveWs(w http.ResponseWriter, r *http.Request) {
	//Let's get the ID
	vars := mux.Vars(r)
	ID, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Warn("No ID provided, bailing out")
		return
	}
	var contextLogger = log.WithFields(log.Fields{
		"ID":   ID,
		"From": r.RemoteAddr,
	})
	contextLogger.Info("New connection")
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	conn := &Conn{send: make(chan Message, 256), ws: ws}
	topologies.t[ID] = &Message{}
	reply := &Reply{
		ID:  ID,
		Rep: make(chan *Hub),
	}
	defer close(reply.Rep)

	AllHubs.Request <- reply
	hub := <-reply.Rep
	hub.register <- conn
	go conn.writePump(ID)
	conn.readPump(ID, hub)
	contextLogger.Info("Connection ended")
}

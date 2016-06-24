package server

import (
	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
	"net/http"
)

// serveWs handles websocket requests from the peer.
func serveWs(w http.ResponseWriter, r *http.Request) {
	//Let's get the Tag
	vars := mux.Vars(r)
	Tag, err := stringToTag(vars["tag"])
	if err != nil {
		log.Warn("No Tag provided, bailing out")
		return
	}
	var contextLogger = log.WithFields(log.Fields{
		"Tag":  Tag,
		"From": r.RemoteAddr,
	})
	contextLogger.Info("New connection")
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	conn := &Conn{send: make(chan OutMessage, 256), ws: ws}
	reply := &reply{
		Message: createMessage(),
		Tag:     Tag,
		Rep:     make(chan *hub),
	}
	defer close(reply.Rep)

	AllHubs.Request <- reply
	hub := <-reply.Rep
	hub.register <- conn
	go conn.writePump(Tag)
	conn.readPump(Tag, hub)
	contextLogger.Info("Connection ended")
}

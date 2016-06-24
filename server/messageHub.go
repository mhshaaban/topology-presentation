package server

import (
	//	log "github.com/Sirupsen/logrus"
	"strconv"
)

// The Tag for channel separation
type Tag int

func stringToTag(s string) (Tag, error) {
	i, err := strconv.Atoi(s)
	return Tag(i), err
}

// MessageHub is a messaging hub
type MessageHub struct {
	// Unregister
	unregister chan int

	// Registered hubs.
	message map[Tag]*Message

	// Register requests from the connections.
	Request chan *Reply
}

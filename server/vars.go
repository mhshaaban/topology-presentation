package server

import (
	"sync"
)

var topologies = struct {
	sync.RWMutex
	t map[Tag]*Message
}{t: make(map[Tag]*Message)}

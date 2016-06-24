package server

import (
	"sync"
)

var topologies = struct {
	sync.RWMutex
	t map[int]*Message
}{t: make(map[int]*Message)}

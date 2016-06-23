package server

import (
	"sync"
)

var topologies = struct {
	sync.RWMutex
	t map[string]*Message
}{t: make(map[string]*Message)}

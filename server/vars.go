package server

import (
	"sync"
)

var topologies = struct {
	sync.RWMutex
	t map[string]*Message
}{t: make(map[string]*Message)}

var hubs = struct {
	sync.RWMutex
	h map[string]*Hub
}{h: make(map[string]*Hub)}

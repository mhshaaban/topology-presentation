package server

// Node is a structure describing a single node
type Node struct {
	Name   string `json:"name"`
	UUID   string `json:"id"`
	Device string `json:"device"`
	Icon   string `json:"icon"`
	Status string `json:"status"`
}

//Link is describing a link between two nodes
type Link struct {
	Source int `json:"source"`
	Target int `json:"target"`
}

//Message is the top envelop for message communication between nodes
type Message struct {
	UUID    string `json:"id"`
	Message string `json:"message"`
	Nodes   []Node `json:"nodes"`
	Links   []Link `json:"links"`
}

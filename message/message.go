package message

import (
	"encoding/json"
	"regexp"
)

// Node is a structure describing a single node
type Node struct {
	Name   string `json:"name"`
	Tag    int    `json:"id"`
	UUID   string `json:"uuid"`
	Device string `json:"device"`
	Icon   string `json:"icon"`
	Status string `json:"status"`
	Color  string `json:"color"`
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

// CreateMessage creates a new message and returns a pointer
func CreateMessage() *Message {
	return &Message{}
}

// Serialize returns a byte array of the message
func (m *Message) Serialize() ([]byte, error) {
	return json.Marshal(m)
}

// Set function updates the content of message m awwording to input n
// And it fills the Msg's interface Contract
func (m *Message) Set(n []byte) error {
	var message Node
	err := json.Unmarshal(n, &message)
	if err != nil {
		return err
	}

	found := false
	nodeTag := 0
	for i, node := range m.Nodes {
		if node.UUID == message.UUID {
			m.Nodes[i].Status = message.Status
			switch message.Status {
			case "initial":
				message.Color = "black"
			case "configured":
				message.Color = "cyan"
			case "started":
				message.Color = "green"
			case "stopped":
				message.Color = "orange"
			case "deleted":
				message.Color = "red"
			default:
				message.Color = "black"
			}
			m.Nodes[i].Color = message.Color
			found = true
		}
		nodeTag = i + 1
	}
	// New node
	if !found {

		var ios = regexp.MustCompile(`(?i).*ios|iphone.*`)
		var android = regexp.MustCompile(`(?i).*android.*`)
		if message.Icon == "" {
			var icon string
			switch {
			case ios.MatchString(message.Device):
				icon = "/img/iphone-phone-color.png"
			case android.MatchString(message.Device):
				icon = "/img/android-phone-color.png"
			default:
				icon = "/img/smartphone.png"
			}
			message.Icon = icon
		}
		message.Tag = nodeTag
		m.Nodes = append(m.Nodes, message)
		// Add a link to the previous node
		if nodeTag >= 1 {
			m.Links = append(m.Links, Link{Source: nodeTag, Target: nodeTag - 1})
		}
	}
	if message.Status == "error" {
		m.Message = "error"
	} else {
		m.Message = "info"
	}
	if len(m.Links) == 0 {
		// Add a dummy link for d3.js
		m.Links = append(m.Links, Link{Source: 0, Target: 0})
	}
	return nil
}

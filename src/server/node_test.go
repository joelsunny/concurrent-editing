package main

import (
	"encoding/json"
	"testing"
	"time"
)

// TestNodeDeltaHandler checks the the node deltaHandle implementation
func TestNodeDeltaHandler(t *testing.T) {
	d := NewDocument()
	c := &MockConn{ReadChan: make(chan Message), WriteChan: make(chan Message)}
	d.AddNode(c)
	// send message to the mock connection
	m := InMsg{Index: 1, Op: Operation{0, "h", 0}, LastIndex: 0}
	b, _ := json.Marshal(m)
	c.WriteMessage(1, b)
	time.Sleep(1 * time.Second)

}

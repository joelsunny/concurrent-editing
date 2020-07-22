package main

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"
)

// TestAddNode checks the the node addition implementation
func TestAddNode(t *testing.T) {
	d := NewDocument()
	c := &MockConn{ReadChan: make(chan Message), WriteChan: make(chan Message)}
	d.AddNode(c)
	fmt.Println(d)
	if (len(d.Nodes) != 1) || (len(d.lastCommits) != 1) {
		t.Errorf("Expected msgCount to be %d but instead got %d!", 1, len(d.Nodes))
	}
}

func TestDeltaHandler(t *testing.T) {
	d := NewDocument()
	c := &MockConn{ReadChan: make(chan Message), WriteChan: make(chan Message)}
	d.AddNode(c)
	// send message to the mock connection
	m := InMsg{Index: 1, Op: Operation{0, "h", 0}, LastIndex: 0}
	b, _ := json.Marshal(m)
	c.WriteMessage(1, b)
	time.Sleep(1 * time.Second)

}

func TestBroadcast(t *testing.T) {
	d := NewDocument()
	c1 := &MockConn{ReadChan: make(chan Message), WriteChan: make(chan Message)}
	d.AddNode(c1)
	c2 := &MockConn{ReadChan: make(chan Message), WriteChan: make(chan Message)}
	d.AddNode(c2)

	// send message to the mock connection
	m := OutChanMsg{log: Log{}, lastCommits: []int{0, 1}}
	d.broadcast(m)
	time.Sleep(2 * time.Second)
}

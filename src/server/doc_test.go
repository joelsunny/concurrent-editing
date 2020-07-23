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
	fmt.Println(d.Nodes)
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
	fmt.Println(d.Doc)
}

func TestBroadcast(t *testing.T) {
	d := NewDocument()
	c1 := &MockConn{ReadChan: make(chan Message), WriteChan: make(chan Message)}
	n1 := d.AddNode(c1)
	c2 := &MockConn{ReadChan: make(chan Message), WriteChan: make(chan Message)}
	n2 := d.AddNode(c2)

	// send message to the mock connection
	m := OutChanMsg{log: Log{}, lastCommits: map[int]int{n1.me: 0, n2.me: 1}, Origin: 3}
	d.broadcast(m)
	time.Sleep(2 * time.Second)
	fmt.Println(d.Doc)
}

package main

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"
)

// TestNodeDeltaHandler checks the the node deltaHandle implementation
func TestTransform(t *testing.T) {
	d := NewDocument()
	c1 := &MockConn{ReadChan: make(chan Message), WriteChan: make(chan Message)}
	d.AddNode(c1)
	c2 := &MockConn{ReadChan: make(chan Message), WriteChan: make(chan Message)}
	d.AddNode(c2)

	m := InMsg{Index: 1, Op: Operation{0, "hello", 0}, LastIndex: 0}
	b, _ := json.Marshal(m)
	c1.WriteMessage(1, b)

	m = InMsg{Index: 1, Op: Operation{0, "world", 0}, LastIndex: 0}
	b, _ = json.Marshal(m)
	c2.WriteMessage(1, b)

	time.Sleep(2 * time.Second)

	fmt.Println(d.Doc)

}

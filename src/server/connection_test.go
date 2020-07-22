package main

import (
	"testing"
	"time"
)

// TestConnection checks the MockConn implementation
func TestConnection(t *testing.T) {
	c := &MockConn{ReadChan: make(chan Message), WriteChan: make(chan Message)}
	n1 := NewNode(c, 1, &Document{})
	n1.Conn.WriteMessage(1, []byte("hello"))
	time.Sleep(1 * time.Second)
	if n1.msgCount != 1 {
		t.Errorf("Expected msgCount to be %d but instead got %d!", 1, n1.msgCount)
	}
}

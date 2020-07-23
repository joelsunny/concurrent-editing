package main

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"
)

// Node structure keeps track of individual connections to a document
type Node struct {
	mu       sync.Mutex
	Conn     Connection
	me       int // id of this node
	OutChan  chan []byte
	StopChan chan bool // Channel to signal closing of document
	Doc      *Document // pointer to the document being modified
	msgCount int
}

// NewNode returns a new Node
func NewNode(Conn Connection, id int, d *Document) *Node {
	n := Node{Conn: Conn, me: id, Doc: d}
	n.OutChan = make(chan []byte, 100)
	n.StopChan = make(chan bool, 10)
	go n.Run()
	return &n
}

// Handle Delta
func (n *Node) deltaHandler(message []byte) {
	n.mu.Lock()
	defer n.mu.Unlock()

	// parse to IncomingDelta
	var m InMsg
	_ = json.Unmarshal(message, &m)
	// m.Print()
	n.msgCount = n.msgCount + 1
	// send to Doc.ReceiveChan
	fmt.Printf("inc insert:%v.\n", m.Op.Insert)
	n.Doc.RecieveChan <- IncomingDelta{index: m.Index, lastIndex: m.LastIndex, op: m.Op, Origin: n.me}
}

// sendDelta back to client
func (n *Node) sendDelta() {
	for {
		select {
		case b := <-n.OutChan:
			// fmt.Printf("%d: received message to be forwarded\n%v\n", n.me, b)
			n.Conn.WriteMessage(1, b)
		case <-n.StopChan:
		}
	}
}

// Run method acts as the main routine for a node
func (n *Node) Run() {

	go n.sendDelta()

	for {
		_, message, err := n.Conn.ReadMessage()
		// fmt.Println("got message ", message)
		if err != nil {
			log.Println("read:", err)
			if err.Error() == "websocket: close 1001 (going away)" {
				fmt.Println("connection terminated")
				n.Doc.RemoveNode(n.me)
				return
			}
		}

		go n.deltaHandler(message)
	}

}

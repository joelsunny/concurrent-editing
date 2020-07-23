package main

import (
	"fmt"
	"sync"
)

// Document structure for representing the document in memory
type Document struct {
	mu          sync.Mutex
	Nodes       map[int]*Node // list of nodes modifying the document
	Doc         []Log         // Document is stored as a list of deltas
	nextNode    int           // unique id for next node
	nextIndex   int           // next log index
	lastCommits map[int]int   // list of last indices received from the nodes, clients maintain there own independent indices
	// this needs to be transmitted back to clients so they can perform transforms on updates received from the server
	// Try: a dumb client that won't do any transforms, the document at the server should still be correct
	RecieveChan chan IncomingDelta // channel to communicate delta
	OutChan     chan OutChanMsg    // Outwards channel of the document
	StopChan    chan bool          // Channel to signal closing of document
}

// NewDocument creates a new document instance
func NewDocument() *Document {
	d := Document{}
	d.nextIndex = 0
	d.nextNode = 0
	d.Nodes = make(map[int]*Node)
	d.lastCommits = make(map[int]int)
	d.RecieveChan = make(chan IncomingDelta, 100)
	d.OutChan = make(chan OutChanMsg, 100)
	d.StopChan = make(chan bool, 10)

	go d.Run()
	return &d
}

// AddNode adds a new user to the document
func (d *Document) AddNode(Conn Connection) *Node {
	d.mu.Lock()
	defer d.mu.Unlock()
	n := NewNode(Conn, d.nextNode, d)
	d.Nodes[d.nextNode] = n
	d.lastCommits[d.nextNode] = 0
	d.nextNode = d.nextNode + 1
	return n
}

// RemoveNode removes the node
func (d *Document) RemoveNode(id int) {
	d.mu.Lock()
	defer d.mu.Unlock()
	fmt.Printf("removing node: %v\n", id)
	delete(d.Nodes, id)
	delete(d.lastCommits, id)
}

// handle delta
func (d *Document) handleDelta(m IncomingDelta) {
	fmt.Printf("received new delta from node %d\n", m.Origin)
	d.mu.Lock()
	defer d.mu.Unlock()
	var unseen []Log
	if m.lastIndex < d.nextIndex-1 {
		unseen = d.Doc[m.lastIndex+1:] // todo?
	}

	// transform and send unseen delta back to node and append incoming delta to Log, send incoming data to remaining nodes
	log := m.Transform(unseen, d.nextIndex)
	fmt.Println("transformed delta, log: ", log)
	d.Doc = append(d.Doc, log)
	d.lastCommits[m.Origin] = m.lastIndex
	d.OutChan <- OutChanMsg{log: log, lastCommits: d.lastCommits, Origin: m.Origin}
	d.nextIndex = d.nextIndex + 1

}

// function to broadcast delta to all nodes
func (d *Document) broadcast(m OutChanMsg) {
	// transform log to outgoingdelta
	fmt.Println("broadcasting")
	for id, n := range d.Nodes {
		if id != m.Origin {
			n.OutChan <- OutgoingDelta{Type: "delta", Log: m.log, LastCommit: m.lastCommits[id]}
		}
	}
}

// Run method monitors the channels of the document
func (d *Document) Run() {
	for {
		select {
		case m := <-d.RecieveChan:
			go d.handleDelta(m)
		case m := <-d.OutChan:
			go d.broadcast(m)
		case <-d.StopChan:
			// signal nodes to stop and exit
			for _, n := range d.Nodes {
				n.StopChan <- true
			}
			return
		}
	}
}

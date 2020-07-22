package main

import (
	"fmt"
	"sync"
)

// Document structure for representing the document in memory
type Document struct {
	mu          sync.Mutex
	Nodes       []*Node // list of nodes modifying the document
	Doc         []Log   // Document is stored as a list of deltas
	nextNode    int     // unique id for next node
	nextIndex   int     // next log index
	lastCommits []int   // list of last indices received from the nodes, clients maintain there own independent indices
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
	d.RecieveChan = make(chan IncomingDelta, 100)
	d.OutChan = make(chan OutChanMsg, 100)
	d.StopChan = make(chan bool, 10)

	go d.Run()
	return &d
}

// AddNode adds a new user to the document
func (d *Document) AddNode(Conn Connection) {
	d.mu.Lock()
	defer d.mu.Unlock()
	n := NewNode(Conn, d.nextNode, d)
	d.nextNode = d.nextNode + 1
	d.Nodes = append(d.Nodes, n)
	d.lastCommits = append(d.lastCommits, 0)
}

// handle delta
func (d *Document) handleDelta(m IncomingDelta) {
	fmt.Printf("received new delta from node %d\n", m.Origin)
	d.mu.Lock()
	defer d.mu.Unlock()
	unseen := d.Doc[m.lastIndex:] // todo?

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
	for i := 0; i < len(d.Nodes); i++ {
		if d.Nodes[i].me == m.Origin {
			continue
		}
		d.Nodes[i].OutChan <- OutgoingDelta{log: m.log, lastCommit: m.lastCommits[i]}
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
			for i := 0; i < len(d.Nodes); i++ {
				d.Nodes[i].StopChan <- true
			}
			return
		}
	}
}

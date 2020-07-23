package main

import "fmt"

// Log structure
type Log struct {
	Index  int       `json:"index"`
	Delta  Operation `json:"delta"`
	origin int
}

// Operation struture
type Operation struct {
	Retain int    `json:"retain"`
	Insert string `json:"insert"`
	Delete int    `json:"delete"`
}

// InMsg describes the structure of message received by nodes
type InMsg struct {
	Type      string    `json:"type"`
	Index     int       `json:"index"` // index maintained by the client
	Op        Operation `json:"op"`
	LastIndex int       `json:"lastIndex"` // last server index that the client knows about
}

// Print InMsg
func (m InMsg) Print() {
	fmt.Println("InMsg: ", m.Index)
	fmt.Println("op: ", m.Op)
}

// OutChanMsg includes the log to be broadcasted and a snapshot of the lastCommits field in the Document
type OutChanMsg struct {
	log         Log
	Origin      int
	lastCommits map[int]int
}

// IncomingDelta structure
type IncomingDelta struct {
	index     int // index maintained by the client
	Origin    int // id of the node
	op        Operation
	lastIndex int // last server index that the client knows about
}

// OutgoingDelta structure
type OutgoingDelta struct {
	Type       string `json:"type"`
	Log        Log    `json:"log"`
	LastCommit int    `json:"lastCommit"`
}

// Transform the delta with respect to unseen deltas
func (m *IncomingDelta) Transform(unseen []Log, index int) Log {

	var log = Log{Index: index, Delta: Operation{}}
	for i := 0; i < len(unseen); i++ {
		uop := unseen[i].Delta
		if (uop.Retain <= m.op.Retain) && (unseen[i].origin != m.Origin) {
			m.op.Retain = m.op.Retain + len(uop.Insert) - uop.Delete
			// // send the unseen update to the client --> no need as all deltas at server will be broadcasted
			// m.Origin.RecieveChan <- uop
		}
	}
	log.Delta = m.op
	log.origin = m.Origin
	return log
}

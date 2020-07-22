package main

import "errors"

// Connection Interface that needs to be satisfied by Node.Conn
type Connection interface {
	ReadMessage() (int, []byte, error)
	WriteMessage(int, []byte) error
	Close() error
}

// Message structure
type Message struct {
	Type int
	Data []byte
}

// MockConn implements Connection interface
type MockConn struct {
	ReadChan  chan Message
	WriteChan chan Message
	CloseChan chan bool
}

// ReadMessage method
func (m *MockConn) ReadMessage() (int, []byte, error) {
	select {
	case d := <-m.ReadChan:
		return d.Type, d.Data, nil
	case <-m.CloseChan:
		return -1, nil, errors.New("closed channel")
	}
}

// WriteMessage method
func (m *MockConn) WriteMessage(mtype int, data []byte) error {
	d := Message{Type: mtype, Data: data}
	m.ReadChan <- d
	return nil
}

// Close the connection
func (m *MockConn) Close() error {
	m.CloseChan <- true
	return nil
}

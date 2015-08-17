package server

import (
	"bufio"
	"encoding/json"
	"net"
)

var currentAgentId int = 0

type agentConn struct {
	socket   net.Conn
	incoming *bufio.Reader
	id       int
	device   *device
}

// create a new connection
func newAgentConnection(conn net.Conn) *agentConn {
	currentAgentId++
	newconn := &agentConn{
		socket:   conn,
		incoming: bufio.NewReader(conn),
		id:       currentAgentId,
	}
	return newconn
}

// read data from a connection
func (conn *agentConn) readData() {
	for {
		line, _, err := conn.incoming.ReadLine()
		if err != nil {
			return
		}
		handleAgentEvent(conn, line)
	}
}

// handle a JSON event
func handleAgentEvent(conn *agentConn, data []byte) bool {
	var i interface{}
	err := json.Unmarshal(data, &i)
	if err != nil {
		return false
	}

	// should be an array.
	c, found := i.([]interface{})
	if !found {
		return false
	}

	// first arg is command, a string
	command, found := c[0].(string)
	if !found {
		return false
	}

	// second arg is the parameters, an object
	params, found := c[1].(map[string]interface{})
	if !found {
		return false
	}

	// if a handler for this command exists, run it
	if agentEventHandlers[command] != nil {
		agentEventHandlers[command](conn, command, params)
	}

	return true
}

// send a JSON event
func (conn *agentConn) send(command string, params map[string]interface{}) bool {
	b, err := json.Marshal(params)
	if err != nil {
		return false
	}
	b = append(b, '\n')
	if _, err = conn.socket.Write(b); err != nil {
		return false
	}
	return true
}

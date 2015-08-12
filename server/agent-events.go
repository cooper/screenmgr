package server

var agentEventHandlers = make(map[string]func(conn *agentConn, name string, params map[string]interface{}))

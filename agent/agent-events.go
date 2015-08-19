package agent

var eventHandlers = make(map[string]func(conn *agentConn, name string, params map[string]interface{}))

package clients

import (
	"lets-go-framework/lets"
)

// Define TCP client, it will used by controller.
var TcpExample = &tcpExample{}

// TCP client definition.
type tcpExample struct {
}

// Send insert user data to TCP server.
func (g *tcpExample) OnConnect() {
	lets.LogD("WE ARE CONNECTED")
	return
}

// Send insert user data to TCP server.
func (g *tcpExample) OnDisconnect() {
	lets.LogD("WE ARE DISCONNECTED")
	return
}
